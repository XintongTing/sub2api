package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/shopspring/decimal"
)

const (
	paypalSandboxAPIBase    = "https://api-m.sandbox.paypal.com"
	paypalLiveAPIBase       = "https://api-m.paypal.com"
	paypalHTTPTimeout      = 15 * time.Second
	paypalMaxResponseSize  = 1 << 20
	paypalMaxErrorSummary  = 512
	paypalAccessTokenSkew  = 2 * time.Minute
	paypalWebhookEventPaid = "CHECKOUT.ORDER.APPROVED"
)

type PayPal struct {
	instanceID string
	config     map[string]string
	httpClient *http.Client
}

type paypalTokenState struct {
	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

var paypalAccessTokens sync.Map

func NewPayPal(instanceID string, config map[string]string) (*PayPal, error) {
	for _, k := range []string{"clientId", "clientSecret", "webhookId"} {
		if strings.TrimSpace(config[k]) == "" {
			return nil, fmt.Errorf("paypal config missing required key: %s", k)
		}
	}
	cfg := cloneStringMap(config)
	currency, err := payment.NormalizePaymentCurrency(cfg["currency"])
	if err != nil {
		return nil, fmt.Errorf("paypal config currency: %w", err)
	}
	cfg["currency"] = currency
	env := strings.ToLower(strings.TrimSpace(cfg["environment"]))
	if env == "" {
		env = "sandbox"
	}
	if env != "sandbox" && env != "live" {
		return nil, fmt.Errorf("paypal environment must be sandbox or live")
	}
	cfg["environment"] = env
	apiBase := strings.TrimSpace(cfg["apiBase"])
	if apiBase == "" {
		if env == "live" {
			apiBase = paypalLiveAPIBase
		} else {
			apiBase = paypalSandboxAPIBase
		}
	}
	normalizedBase, err := normalizePayPalAPIBase(apiBase)
	if err != nil {
		return nil, err
	}
	cfg["apiBase"] = normalizedBase
	return &PayPal{
		instanceID: instanceID,
		config:     cfg,
		httpClient: &http.Client{Timeout: paypalHTTPTimeout},
	}, nil
}

func normalizePayPalAPIBase(raw string) (string, error) {
	base := strings.TrimSpace(raw)
	parsed, err := url.Parse(base)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return "", fmt.Errorf("paypal apiBase must be an HTTPS URL")
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	parsed.RawPath = ""
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	return parsed.String(), nil
}

func (p *PayPal) Name() string        { return "PayPal" }
func (p *PayPal) ProviderKey() string { return payment.TypePayPal }
func (p *PayPal) SupportedTypes() []payment.PaymentType {
	return []payment.PaymentType{payment.TypePayPal}
}

func (p *PayPal) MerchantIdentityMetadata() map[string]string {
	if p == nil {
		return nil
	}
	return map[string]string{
		"currency":    p.currency(),
		"environment": p.config["environment"],
	}
}

func (p *PayPal) currency() string {
	if p == nil {
		return payment.DefaultPaymentCurrency
	}
	currency, err := payment.NormalizePaymentCurrency(p.config["currency"])
	if err != nil {
		return payment.DefaultPaymentCurrency
	}
	return currency
}

func (p *PayPal) CreatePayment(ctx context.Context, req payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	amount, err := decimal.NewFromString(strings.TrimSpace(req.Amount))
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("paypal create order: invalid amount %s", req.Amount)
	}
	token, err := p.accessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("paypal auth: %w", err)
	}
	currency := p.currency()
	payload := paypalCreateOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []paypalPurchaseUnit{{
			ReferenceID: req.OrderID,
			CustomID:    req.OrderID,
			Description: req.Subject,
			Amount: paypalAmount{
				CurrencyCode: currency,
				Value:        payment.FormatAmountForCurrency(amount.InexactFloat64(), currency),
			},
		}},
		ApplicationContext: paypalApplicationContext{
			ReturnURL: req.ReturnURL,
			CancelURL: req.ReturnURL,
			UserAction: "PAY_NOW",
		},
	}
	var order paypalOrder
	if err := p.doJSON(ctx, http.MethodPost, "/v2/checkout/orders", token, payload, &order, "create-order-"+req.OrderID); err != nil {
		return nil, fmt.Errorf("paypal create order: %w", err)
	}
	payURL := paypalOrderLink(order.Links, "approve")
	if strings.TrimSpace(order.ID) == "" {
		return nil, fmt.Errorf("paypal create order: missing order id")
	}
	if payURL == "" {
		return nil, fmt.Errorf("paypal create order: missing approve URL")
	}
	return &payment.CreatePaymentResponse{
		TradeNo:    order.ID,
		PayURL:     payURL,
		Currency:   currency,
		PaymentEnv: p.config["environment"],
		ResultType: payment.CreatePaymentResultOrderCreated,
	}, nil
}

func (p *PayPal) QueryOrder(ctx context.Context, tradeNo string) (*payment.QueryOrderResponse, error) {
	orderID := strings.TrimSpace(tradeNo)
	if orderID == "" {
		return nil, fmt.Errorf("paypal query order: missing order id")
	}
	token, err := p.accessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("paypal auth: %w", err)
	}
	order, err := p.getOrder(ctx, token, orderID)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(order.Status, "APPROVED") {
		if err := p.doJSON(ctx, http.MethodPost, "/v2/checkout/orders/"+url.PathEscape(orderID)+"/capture", token, nil, &order, "capture-order-"+orderID); err != nil {
			return nil, fmt.Errorf("paypal capture order: %w", err)
		}
	}
	return paypalQueryResponse(order, orderID, p.currency()), nil
}

func (p *PayPal) getOrder(ctx context.Context, token, orderID string) (paypalOrder, error) {
	var order paypalOrder
	if err := p.doJSON(ctx, http.MethodGet, "/v2/checkout/orders/"+url.PathEscape(orderID), token, nil, &order, ""); err != nil {
		return paypalOrder{}, fmt.Errorf("paypal query order: %w", err)
	}
	return order, nil
}

func (p *PayPal) VerifyNotification(ctx context.Context, rawBody string, headers map[string]string) (*payment.PaymentNotification, error) {
	token, err := p.accessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("paypal auth: %w", err)
	}
	var event paypalWebhookEvent
	if err := json.Unmarshal([]byte(rawBody), &event); err != nil {
		return nil, fmt.Errorf("paypal parse webhook: %w", err)
	}
	ok, err := p.verifyWebhook(ctx, token, rawBody, headers, event)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("paypal notification invalid signature")
	}
	if event.EventType != paypalWebhookEventPaid {
		return nil, nil
	}
	var order paypalOrder
	if len(event.Resource) > 0 {
		if err := json.Unmarshal(event.Resource, &order); err != nil {
			return nil, fmt.Errorf("paypal parse order resource: %w", err)
		}
	}
	if strings.TrimSpace(order.ID) == "" {
		return nil, fmt.Errorf("paypal webhook missing order id")
	}
	query, err := p.QueryOrder(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	if query.Status != payment.ProviderStatusPaid {
		return nil, nil
	}
	orderID := paypalOrderReference(order)
	if orderID == "" {
		orderID = query.Metadata["order_id"]
	}
	if orderID == "" {
		return nil, fmt.Errorf("paypal webhook missing merchant order id")
	}
	return &payment.PaymentNotification{
		TradeNo: query.TradeNo,
		OrderID: orderID,
		Amount:  query.Amount,
		Status:  payment.NotificationStatusSuccess,
		RawData: rawBody,
		Metadata: map[string]string{
			"currency": query.Metadata["currency"],
			"status":   query.Metadata["status"],
		},
	}, nil
}

func (p *PayPal) Refund(ctx context.Context, req payment.RefundRequest) (*payment.RefundResponse, error) {
	ref := strings.TrimSpace(req.TradeNo)
	if ref == "" {
		return nil, fmt.Errorf("paypal refund: missing order or capture id")
	}
	token, err := p.accessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("paypal auth: %w", err)
	}
	captureID, err := p.refundCaptureID(ctx, token, ref)
	if err != nil {
		return nil, err
	}
	payload := paypalRefundRequest{
		Amount: paypalAmount{
			CurrencyCode: p.currency(),
			Value:        req.Amount,
		},
		NoteToPayer: strings.TrimSpace(req.Reason),
	}
	var refund paypalRefund
	if err := p.doJSON(ctx, http.MethodPost, "/v2/payments/captures/"+url.PathEscape(captureID)+"/refund", token, payload, &refund, "refund-"+req.OrderID+"-"+captureID); err != nil {
		return nil, fmt.Errorf("paypal refund: %w", err)
	}
	return &payment.RefundResponse{
		RefundID: refund.ID,
		Status:   paypalRefundStatus(refund.Status),
	}, nil
}

func (p *PayPal) refundCaptureID(ctx context.Context, token, ref string) (string, error) {
	order, err := p.getOrder(ctx, token, ref)
	if err != nil {
		return ref, nil
	}
	captureID, _, _, _ := paypalCaptureMetadata(order)
	if captureID == "" {
		return "", fmt.Errorf("paypal refund: missing capture id")
	}
	return captureID, nil
}

func (p *PayPal) accessToken(ctx context.Context) (string, error) {
	cacheKey := p.config["apiBase"] + "|" + p.config["clientId"]
	stateAny, _ := paypalAccessTokens.LoadOrStore(cacheKey, &paypalTokenState{})
	state := stateAny.(*paypalTokenState)
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.token != "" && time.Now().Add(paypalAccessTokenSkew).Before(state.expiresAt) {
		return state.token, nil
	}
	endpoint, err := joinPayPalURL(p.config["apiBase"], "/v1/oauth2/token")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(p.config["clientId"], p.config["clientSecret"])
	client := p.client()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, paypalMaxResponseSize))
	if err != nil {
		return "", err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, summarizePayPalResponse(raw))
	}
	var token paypalAccessTokenResponse
	if err := json.Unmarshal(raw, &token); err != nil {
		return "", fmt.Errorf("decode token JSON: %w", err)
	}
	if strings.TrimSpace(token.AccessToken) == "" {
		return "", fmt.Errorf("paypal auth: missing access token")
	}
	expiresIn := time.Duration(token.ExpiresIn) * time.Second
	if expiresIn <= 0 {
		expiresIn = time.Hour
	}
	state.token = token.AccessToken
	state.expiresAt = time.Now().Add(expiresIn)
	return state.token, nil
}

func (p *PayPal) verifyWebhook(ctx context.Context, token, rawBody string, headers map[string]string, event paypalWebhookEvent) (bool, error) {
	payload := paypalVerifyWebhookRequest{
		AuthAlgo:         firstPayPalHeader(headers, "paypal-auth-algo"),
		CertURL:          firstPayPalHeader(headers, "paypal-cert-url"),
		TransmissionID:   firstPayPalHeader(headers, "paypal-transmission-id"),
		TransmissionSig:  firstPayPalHeader(headers, "paypal-transmission-sig"),
		TransmissionTime: firstPayPalHeader(headers, "paypal-transmission-time"),
		WebhookID:        p.config["webhookId"],
		WebhookEvent:     json.RawMessage(rawBody),
	}
	if payload.AuthAlgo == "" || payload.CertURL == "" || payload.TransmissionID == "" || payload.TransmissionSig == "" || payload.TransmissionTime == "" {
		return false, fmt.Errorf("paypal notification missing signature headers")
	}
	if strings.TrimSpace(event.ID) == "" {
		return false, fmt.Errorf("paypal notification missing event id")
	}
	var resp paypalVerifyWebhookResponse
	if err := p.doJSON(ctx, http.MethodPost, "/v1/notifications/verify-webhook-signature", token, payload, &resp, ""); err != nil {
		return false, fmt.Errorf("paypal verify notification: %w", err)
	}
	return strings.EqualFold(resp.VerificationStatus, "SUCCESS"), nil
}

func (p *PayPal) doJSON(ctx context.Context, method, path, token string, payload any, out any, idempotencyKey string) error {
	endpoint, err := joinPayPalURL(p.config["apiBase"], path)
	if err != nil {
		return err
	}
	var body io.Reader
	if payload != nil {
		raw, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return marshalErr
		}
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if strings.TrimSpace(idempotencyKey) != "" {
		req.Header.Set("PayPal-Request-Id", idempotencyKey)
	}
	client := p.client()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, paypalMaxResponseSize))
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, summarizePayPalResponse(raw))
	}
	if out == nil || len(bytes.TrimSpace(raw)) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}
	return nil
}

func (p *PayPal) client() *http.Client {
	if p != nil && p.httpClient != nil {
		return p.httpClient
	}
	return &http.Client{Timeout: paypalHTTPTimeout}
}

func joinPayPalURL(base, path string) (string, error) {
	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("paypal API path is required")
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return parsed.String(), nil
}

func paypalQueryResponse(order paypalOrder, fallbackTradeNo, fallbackCurrency string) *payment.QueryOrderResponse {
	currency := fallbackCurrency
	amount := decimal.Zero
	orderID := paypalOrderReference(order)
	if len(order.PurchaseUnits) > 0 {
		pu := order.PurchaseUnits[0]
		if strings.TrimSpace(pu.Amount.CurrencyCode) != "" {
			currency = pu.Amount.CurrencyCode
		}
		if parsed, err := decimal.NewFromString(pu.Amount.Value); err == nil {
			amount = parsed
		}
	}
	captureID, captureStatus, captureAmount, captureCurrency := paypalCaptureMetadata(order)
	if captureAmount.GreaterThan(decimal.Zero) {
		amount = captureAmount
	}
	if captureCurrency != "" {
		currency = captureCurrency
	}
	metadata := map[string]string{
		"currency": currency,
		"status":   order.Status,
		"order_id": orderID,
	}
	if captureID != "" {
		metadata["capture_id"] = captureID
	}
	return &payment.QueryOrderResponse{
		TradeNo:  firstNonEmptyPayPal(order.ID, fallbackTradeNo),
		Status:   paypalProviderStatus(order.Status, captureStatus),
		Amount:   amount.InexactFloat64(),
		Metadata: metadata,
	}
}

func paypalCaptureMetadata(order paypalOrder) (string, string, decimal.Decimal, string) {
	for _, unit := range order.PurchaseUnits {
		for _, capture := range unit.Payments.Captures {
			if strings.TrimSpace(capture.ID) == "" {
				continue
			}
			amount := decimal.Zero
			if parsed, err := decimal.NewFromString(capture.Amount.Value); err == nil {
				amount = parsed
			}
			return capture.ID, capture.Status, amount, capture.Amount.CurrencyCode
		}
	}
	return "", "", decimal.Zero, ""
}

func paypalOrderReference(order paypalOrder) string {
	for _, unit := range order.PurchaseUnits {
		if ref := firstNonEmptyPayPal(unit.CustomID, unit.ReferenceID, unit.InvoiceID); ref != "" {
			return ref
		}
	}
	return ""
}

func paypalOrderLink(links []paypalLink, rel string) string {
	for _, link := range links {
		if strings.EqualFold(link.Rel, rel) && strings.TrimSpace(link.Href) != "" {
			return strings.TrimSpace(link.Href)
		}
	}
	return ""
}

func paypalProviderStatus(orderStatus, captureStatus string) string {
	if strings.EqualFold(captureStatus, "COMPLETED") {
		return payment.ProviderStatusPaid
	}
	switch strings.ToUpper(strings.TrimSpace(orderStatus)) {
	case "COMPLETED":
		return payment.ProviderStatusPaid
	case "VOIDED", "FAILED", "EXPIRED":
		return payment.ProviderStatusFailed
	default:
		return payment.ProviderStatusPending
	}
}

func paypalRefundStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "COMPLETED":
		return payment.ProviderStatusSuccess
	case "FAILED", "CANCELLED":
		return payment.ProviderStatusFailed
	default:
		return payment.ProviderStatusPending
	}
}

func summarizePayPalResponse(body []byte) string {
	summary := strings.Join(strings.Fields(string(body)), " ")
	if summary == "" {
		return "<empty>"
	}
	if len(summary) > paypalMaxErrorSummary {
		return summary[:paypalMaxErrorSummary] + "..."
	}
	return summary
}

func firstPayPalHeader(headers map[string]string, name string) string {
	if len(headers) == 0 {
		return ""
	}
	for key, value := range headers {
		if strings.EqualFold(key, name) {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstNonEmptyPayPal(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

type paypalAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type paypalCreateOrderRequest struct {
	Intent             string                    `json:"intent"`
	PurchaseUnits      []paypalPurchaseUnit      `json:"purchase_units"`
	ApplicationContext paypalApplicationContext `json:"application_context,omitempty"`
}

type paypalApplicationContext struct {
	ReturnURL string `json:"return_url,omitempty"`
	CancelURL string `json:"cancel_url,omitempty"`
	UserAction string `json:"user_action,omitempty"`
}

type paypalPurchaseUnit struct {
	ReferenceID string         `json:"reference_id,omitempty"`
	CustomID    string         `json:"custom_id,omitempty"`
	InvoiceID   string         `json:"invoice_id,omitempty"`
	Description string         `json:"description,omitempty"`
	Amount      paypalAmount   `json:"amount"`
	Payments    paypalPayments `json:"payments,omitempty"`
}

type paypalAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type paypalPayments struct {
	Captures []paypalCapture `json:"captures,omitempty"`
}

type paypalCapture struct {
	ID     string       `json:"id"`
	Status string       `json:"status"`
	Amount paypalAmount `json:"amount"`
}

type paypalOrder struct {
	ID            string               `json:"id"`
	Status        string               `json:"status"`
	Links         []paypalLink         `json:"links"`
	PurchaseUnits []paypalPurchaseUnit `json:"purchase_units"`
}

type paypalLink struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type paypalWebhookEvent struct {
	ID        string          `json:"id"`
	EventType string          `json:"event_type"`
	Resource  json.RawMessage `json:"resource"`
}

type paypalVerifyWebhookRequest struct {
	AuthAlgo         string          `json:"auth_algo"`
	CertURL          string          `json:"cert_url"`
	TransmissionID   string          `json:"transmission_id"`
	TransmissionSig  string          `json:"transmission_sig"`
	TransmissionTime string          `json:"transmission_time"`
	WebhookID        string          `json:"webhook_id"`
	WebhookEvent     json.RawMessage `json:"webhook_event"`
}

type paypalVerifyWebhookResponse struct {
	VerificationStatus string `json:"verification_status"`
}

type paypalRefundRequest struct {
	Amount      paypalAmount `json:"amount"`
	NoteToPayer string       `json:"note_to_payer,omitempty"`
}

type paypalRefund struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

var (
	_ payment.Provider                 = (*PayPal)(nil)
	_ payment.MerchantIdentityProvider = (*PayPal)(nil)
)
