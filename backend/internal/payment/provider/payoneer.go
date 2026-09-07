package provider

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/shopspring/decimal"
)

const (
	payoneerSandboxAPIBase     = "https://api.sandbox.payoneer.com"
	payoneerLiveAPIBase        = "https://api.payoneer.com"
	payoneerDefaultCreatePath  = "/checkout/payment"
	payoneerDefaultQueryPath   = "/checkout/payment/{trade_no}"
	payoneerHTTPTimeout       = 15 * time.Second
	payoneerMaxResponseSize   = 1 << 20
	payoneerMaxErrorSummary   = 512
	payoneerWebhookTolerance  = 5 * time.Minute
)

type Payoneer struct {
	instanceID string
	config     map[string]string
	httpClient *http.Client
}

func NewPayoneer(instanceID string, config map[string]string) (*Payoneer, error) {
	for _, k := range []string{"clientId", "clientSecret", "webhookSecret"} {
		if strings.TrimSpace(config[k]) == "" {
			return nil, fmt.Errorf("payoneer config missing required key: %s", k)
		}
	}
	cfg := cloneStringMap(config)
	currency, err := payment.NormalizePaymentCurrency(cfg["currency"])
	if err != nil {
		return nil, fmt.Errorf("payoneer config currency: %w", err)
	}
	cfg["currency"] = currency
	env := strings.ToLower(strings.TrimSpace(cfg["environment"]))
	if env == "" {
		env = "sandbox"
	}
	if env != "sandbox" && env != "live" {
		return nil, fmt.Errorf("payoneer environment must be sandbox or live")
	}
	cfg["environment"] = env
	apiBase := strings.TrimSpace(cfg["apiBase"])
	if apiBase == "" {
		if env == "live" {
			apiBase = payoneerLiveAPIBase
		} else {
			apiBase = payoneerSandboxAPIBase
		}
	}
	normalizedBase, err := normalizePayoneerAPIBase(apiBase)
	if err != nil {
		return nil, err
	}
	cfg["apiBase"] = normalizedBase
	return &Payoneer{
		instanceID: instanceID,
		config:     cfg,
		httpClient: &http.Client{Timeout: payoneerHTTPTimeout},
	}, nil
}

func normalizePayoneerAPIBase(raw string) (string, error) {
	base := strings.TrimSpace(raw)
	parsed, err := url.Parse(base)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return "", fmt.Errorf("payoneer apiBase must be an HTTPS URL")
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	parsed.RawPath = ""
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	return parsed.String(), nil
}

func (p *Payoneer) Name() string        { return "Payoneer Checkout" }
func (p *Payoneer) ProviderKey() string { return payment.TypePayoneer }
func (p *Payoneer) SupportedTypes() []payment.PaymentType {
	return []payment.PaymentType{payment.TypePayoneer}
}

func (p *Payoneer) MerchantIdentityMetadata() map[string]string {
	if p == nil {
		return nil
	}
	return map[string]string{
		"currency":    p.currency(),
		"environment": p.config["environment"],
	}
}

func (p *Payoneer) currency() string {
	if p == nil {
		return payment.DefaultPaymentCurrency
	}
	currency, err := payment.NormalizePaymentCurrency(p.config["currency"])
	if err != nil {
		return payment.DefaultPaymentCurrency
	}
	return currency
}

func (p *Payoneer) CreatePayment(ctx context.Context, req payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	amount, err := decimal.NewFromString(strings.TrimSpace(req.Amount))
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("payoneer create payment: invalid amount %s", req.Amount)
	}
	currency := p.currency()
	payload := map[string]any{
		"merchant_reference": req.OrderID,
		"order_id":           req.OrderID,
		"amount":             amount,
		"currency":           currency,
		"description":        req.Subject,
		"return_url":         req.ReturnURL,
		"callback_url":       req.NotifyURL,
	}
	var resp payoneerPaymentResponse
	path := configValue(p.config, "createPath", payoneerDefaultCreatePath)
	if err := p.doJSON(ctx, http.MethodPost, path, payload, &resp); err != nil {
		return nil, fmt.Errorf("payoneer create payment: %w", err)
	}
	tradeNo := firstNonEmpty(resp.ID, resp.LongID, resp.TransactionID, resp.PaymentID, resp.OrderID)
	payURL := firstNonEmpty(resp.RedirectURL, resp.CheckoutURL, resp.PayURL, resp.HostedURL)
	if tradeNo == "" {
		return nil, fmt.Errorf("payoneer create payment: missing transaction id")
	}
	if payURL == "" {
		return nil, fmt.Errorf("payoneer create payment: missing checkout URL")
	}
	return &payment.CreatePaymentResponse{
		TradeNo:    tradeNo,
		PayURL:     payURL,
		Currency:   currency,
		PaymentEnv: p.config["environment"],
		ResultType: payment.CreatePaymentResultOrderCreated,
	}, nil
}

func (p *Payoneer) QueryOrder(ctx context.Context, tradeNo string) (*payment.QueryOrderResponse, error) {
	tradeNo = strings.TrimSpace(tradeNo)
	if tradeNo == "" {
		return nil, fmt.Errorf("payoneer query order: missing trade number")
	}
	path := strings.ReplaceAll(configValue(p.config, "queryPath", payoneerDefaultQueryPath), "{trade_no}", url.PathEscape(tradeNo))
	var resp payoneerPaymentResponse
	if err := p.doJSON(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("payoneer query order: %w", err)
	}
	currency := firstNonEmpty(resp.Currency, p.currency())
	return &payment.QueryOrderResponse{
		TradeNo: firstNonEmpty(resp.ID, resp.LongID, resp.TransactionID, resp.PaymentID, tradeNo),
		Status:  payoneerProviderStatus(resp.Status),
		Amount:  resp.Amount.InexactFloat64(),
		Metadata: map[string]string{
			"currency": currency,
			"status":   resp.Status,
		},
	}, nil
}

func (p *Payoneer) VerifyNotification(_ context.Context, rawBody string, headers map[string]string) (*payment.PaymentNotification, error) {
	if err := verifyPayoneerWebhookSignature(rawBody, headers, p.config["webhookSecret"], time.Now()); err != nil {
		return nil, err
	}
	var event payoneerPaymentResponse
	if err := json.Unmarshal([]byte(rawBody), &event); err != nil {
		return nil, fmt.Errorf("payoneer parse webhook: %w", err)
	}
	status := payoneerProviderStatus(event.Status)
	if status != payment.ProviderStatusPaid && status != payment.ProviderStatusFailed {
		return nil, nil
	}
	tradeNo := firstNonEmpty(event.ID, event.LongID, event.TransactionID, event.PaymentID)
	orderID := firstNonEmpty(event.MerchantReference, event.OrderID)
	if tradeNo == "" || orderID == "" {
		return nil, fmt.Errorf("payoneer webhook missing transaction id or merchant reference")
	}
	notificationStatus := payment.ProviderStatusFailed
	if status == payment.ProviderStatusPaid {
		notificationStatus = payment.NotificationStatusSuccess
	}
	currency := firstNonEmpty(event.Currency, p.currency())
	return &payment.PaymentNotification{
		TradeNo: tradeNo,
		OrderID: orderID,
		Amount:  event.Amount.InexactFloat64(),
		Status:  notificationStatus,
		RawData: rawBody,
		Metadata: map[string]string{
			"currency": currency,
			"status":   event.Status,
		},
	}, nil
}

func (p *Payoneer) Refund(context.Context, payment.RefundRequest) (*payment.RefundResponse, error) {
	return nil, fmt.Errorf("payoneer refund is not implemented")
}

func (p *Payoneer) doJSON(ctx context.Context, method, path string, payload any, out any) error {
	endpoint, err := joinPayoneerURL(p.config["apiBase"], path)
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
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.SetBasicAuth(p.config["clientId"], p.config["clientSecret"])
	req.Header.Set("Idempotency-Key", "payoneer-"+configValue(p.config, "instanceId", p.instanceID))
	client := p.httpClient
	if client == nil {
		client = &http.Client{Timeout: payoneerHTTPTimeout}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, payoneerMaxResponseSize))
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, summarizePayoneerResponse(raw))
	}
	if out == nil || len(bytes.TrimSpace(raw)) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}
	return nil
}

func joinPayoneerURL(base, path string) (string, error) {
	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("payoneer API path is required")
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return parsed.String(), nil
}

func verifyPayoneerWebhookSignature(rawBody string, headers map[string]string, secret string, now time.Time) error {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return fmt.Errorf("payoneer webhookSecret not configured")
	}
	signature := firstHeader(headers, "x-payoneer-signature", "payoneer-signature", "x-signature")
	if signature == "" {
		return fmt.Errorf("payoneer notification missing signature header")
	}
	timestamp := firstHeader(headers, "x-payoneer-timestamp", "payoneer-timestamp", "x-timestamp")
	payload := rawBody
	if timestamp != "" {
		ts, err := parsePayoneerWebhookTimestamp(timestamp)
		if err != nil {
			return err
		}
		if diff := now.Sub(ts).Abs(); diff > payoneerWebhookTolerance {
			return fmt.Errorf("payoneer webhook timestamp outside tolerance")
		}
		payload = timestamp + "." + rawBody
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	signature = strings.TrimPrefix(strings.TrimSpace(signature), "sha256=")
	if !hmac.Equal([]byte(strings.ToLower(signature)), []byte(expected)) {
		return fmt.Errorf("payoneer notification invalid signature")
	}
	return nil
}

func parsePayoneerWebhookTimestamp(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, fmt.Errorf("payoneer webhook timestamp is empty")
	}
	if ts, err := time.Parse(time.RFC3339, raw); err == nil {
		return ts, nil
	}
	unix, err := decimal.NewFromString(raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("payoneer webhook timestamp invalid")
	}
	if unix.GreaterThan(decimal.NewFromInt(9999999999)) {
		return time.UnixMilli(unix.IntPart()), nil
	}
	return time.Unix(unix.IntPart(), 0), nil
}

func firstHeader(headers map[string]string, names ...string) string {
	for _, name := range names {
		for k, v := range headers {
			if strings.EqualFold(k, name) {
				return strings.TrimSpace(v)
			}
		}
	}
	return ""
}

func payoneerProviderStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "PAID", "SUCCESS", "SUCCEEDED", "COMPLETED", "CAPTURED", "APPROVED":
		return payment.ProviderStatusPaid
	case "FAILED", "CANCELLED", "CANCELED", "DECLINED", "EXPIRED":
		return payment.ProviderStatusFailed
	case "REFUNDED":
		return payment.ProviderStatusRefunded
	default:
		return payment.ProviderStatusPending
	}
}

func summarizePayoneerResponse(body []byte) string {
	summary := strings.Join(strings.Fields(string(body)), " ")
	if summary == "" {
		return "<empty>"
	}
	if len(summary) > payoneerMaxErrorSummary {
		return summary[:payoneerMaxErrorSummary] + "..."
	}
	return summary
}

func configValue(cfg map[string]string, key, fallback string) string {
	if v := strings.TrimSpace(cfg[key]); v != "" {
		return v
	}
	return fallback
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

type payoneerPaymentResponse struct {
	ID                string          `json:"id"`
	LongID            string          `json:"longId"`
	TransactionID     string          `json:"transaction_id"`
	PaymentID         string          `json:"payment_id"`
	OrderID           string          `json:"order_id"`
	MerchantReference string          `json:"merchant_reference"`
	RedirectURL       string          `json:"redirect_url"`
	CheckoutURL       string          `json:"checkout_url"`
	PayURL            string          `json:"pay_url"`
	HostedURL         string          `json:"hosted_url"`
	Status            string          `json:"status"`
	Amount            decimal.Decimal `json:"amount"`
	Currency          string          `json:"currency"`
}

var (
	_ payment.Provider                 = (*Payoneer)(nil)
	_ payment.MerchantIdentityProvider = (*Payoneer)(nil)
)
