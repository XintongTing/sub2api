package provider

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/shopspring/decimal"
)

const (
	sunrateTestAPIBase = "https://test-api.xunhuiacq.com"
	sunrateLiveAPIBase = "https://api.acq.sunrate.com"
	sunrateHTTPTimeout = 15 * time.Second
)

type Sunrate struct { instanceID string; config map[string]string; httpClient *http.Client }

func NewSunrate(instanceID string, config map[string]string) (*Sunrate, error) {
	for _, k := range []string{"merchantId", "signatureKey"} {
		if strings.TrimSpace(config[k]) == "" { return nil, fmt.Errorf("sunrate config missing required key: %s", k) }
	}
	cfg := cloneStringMap(config)
	env := strings.ToLower(strings.TrimSpace(cfg["environment"])); if env == "" { env = "sandbox" }
	if env != "sandbox" && env != "live" { return nil, fmt.Errorf("sunrate environment must be sandbox or live") }
	cfg["environment"] = env
	if strings.TrimSpace(cfg["currency"]) == "" { cfg["currency"] = "THB" }
	currency, err := payment.NormalizePaymentCurrency(cfg["currency"]); if err != nil { return nil, fmt.Errorf("sunrate config currency: %w", err) }; cfg["currency"] = currency
	base := strings.TrimSpace(cfg["apiBase"]); if base == "" { if env == "live" { base = sunrateLiveAPIBase } else { base = sunrateTestAPIBase } }
	u, err := url.Parse(base); if err != nil || u.Scheme != "https" || u.Host == "" { return nil, fmt.Errorf("sunrate apiBase must be an HTTPS URL") }
	u.RawQuery = ""; u.Fragment = ""; u.Path = strings.TrimRight(u.Path, "/"); cfg["apiBase"] = u.String()
	return &Sunrate{instanceID: instanceID, config: cfg, httpClient: &http.Client{Timeout: sunrateHTTPTimeout}}, nil
}

func (s *Sunrate) Name() string { return "SUNRATE Cashier" }
func (s *Sunrate) ProviderKey() string { return payment.TypeSunrate }
func (s *Sunrate) SupportedTypes() []payment.PaymentType { return []payment.PaymentType{payment.TypeSunrate} }
func (s *Sunrate) MerchantIdentityMetadata() map[string]string { return map[string]string{"merchant_id": strings.TrimSpace(s.config["merchantId"]), "currency": s.config["currency"], "environment": s.config["environment"]} }

func (s *Sunrate) CreatePayment(ctx context.Context, req payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	amount, err := decimal.NewFromString(strings.TrimSpace(req.Amount)); if err != nil || amount.LessThanOrEqual(decimal.Zero) { return nil, fmt.Errorf("sunrate create payment: invalid amount %s", req.Amount) }
	front, back := strings.TrimSpace(req.ReturnURL), strings.TrimSpace(req.NotifyURL); if front == "" || back == "" { return nil, fmt.Errorf("sunrate create payment: return and notify URLs are required") }
	osType := "WINDOWS"; if req.IsMobile { osType = "OTHER" }
	goods := strings.TrimSpace(req.Subject); if goods == "" { goods = "Service" }; goods = strings.NewReplacer("#", "", "&", "", "+", "").Replace(goods)
	detailRaw, _ := json.Marshal([]map[string]string{{"goods_id": req.OrderID, "goods_name": goods, "quantity": "1", "price": payment.FormatAmountForCurrency(amount.InexactFloat64(), s.config["currency"])}})
	fields := map[string]string{"paymentMethod":"sunrate_cashier", "filters": strings.TrimSpace(s.config["filters"]), "orderNum": req.OrderID, "orderAmount": payment.FormatAmountForCurrency(amount.InexactFloat64(), s.config["currency"]), "orderCurrency": s.config["currency"], "frontURL": front, "backURL": back, "merID": s.config["merchantId"], "goodsInfo": goods, "detailInfo": base64.StdEncoding.EncodeToString(detailRaw), "transTime": time.Now().Format("20060102150405"), "osType": osType, "userIP": req.ClientIP, "userID": firstNonEmpty(req.OrderID, s.config["userID"]), "signType":"SHA256"}
	for k, v := range fields { if strings.TrimSpace(v) == "" { delete(fields, k) } }
	fields["signature"] = sunrateSign(fields, s.config["signatureKey"])
	var resp sunrateCreateResponse; if err := s.post(ctx, "/api/v5/createcashier", fields, &resp); err != nil { return nil, fmt.Errorf("sunrate create payment: %w", err) }
	if resp.RespCode != "00" { return nil, fmt.Errorf("sunrate create payment: %s (%s)", resp.RespMsg, resp.RespCode) }
	if strings.TrimSpace(resp.Parameter.PayURL) == "" { return nil, fmt.Errorf("sunrate create payment: missing payUrl") }
	return &payment.CreatePaymentResponse{TradeNo: req.OrderID, PayURL: resp.Parameter.PayURL, Currency: s.config["currency"], PaymentEnv: s.config["environment"], ResultType: payment.CreatePaymentResultOrderCreated}, nil
}

func (s *Sunrate) QueryOrder(ctx context.Context, tradeNo string) (*payment.QueryOrderResponse, error) {
	if strings.TrimSpace(tradeNo) == "" { return nil, fmt.Errorf("sunrate query order: missing order number") }
	fields := map[string]string{"paymentMethod":"sunrate_cashier", "orderNum": tradeNo, "merID": s.config["merchantId"], "transTime": time.Now().Format("20060102150405"), "signType":"SHA256"}; fields["signature"] = sunrateSign(fields, s.config["signatureKey"])
	var resp sunrateNotifyResponse; if err := s.post(ctx, "/api/v5/order_query", fields, &resp); err != nil { return nil, fmt.Errorf("sunrate query order: %w", err) }
	status := payment.ProviderStatusPending; if resp.RespCode == "00" { status = payment.ProviderStatusPaid } else if resp.RespCode != "01" { status = payment.ProviderStatusFailed }
	amount, _ := decimal.NewFromString(resp.OrderAmount)
	return &payment.QueryOrderResponse{TradeNo: firstNonEmpty(resp.TransID, resp.OrderNum, tradeNo), Status: status, Amount: amount.InexactFloat64(), Metadata: map[string]string{"currency": resp.OrderCurrency, "payment_brand": resp.PaymentBrand, "resp_code": resp.RespCode}}, nil
}

func (s *Sunrate) VerifyNotification(_ context.Context, rawBody string, _ map[string]string) (*payment.PaymentNotification, error) {
	var n sunrateNotifyResponse; if err := json.Unmarshal([]byte(rawBody), &n); err != nil { return nil, fmt.Errorf("sunrate parse webhook: %w", err) }
	if strings.TrimSpace(n.Signature) == "" || !equalSignature(n.Signature, sunrateSign(n.signingFields(), s.config["signatureKey"])) { return nil, fmt.Errorf("sunrate notification invalid signature") }
	if n.MerID != "" && n.MerID != s.config["merchantId"] { return nil, fmt.Errorf("sunrate notification merchant mismatch") }
	if n.RespCode != "00" && n.RespCode != "01" { return nil, nil }
	status := payment.ProviderStatusFailed; if n.RespCode == "00" { status = payment.NotificationStatusSuccess }
	amount, _ := decimal.NewFromString(n.OrderAmount)
	return &payment.PaymentNotification{TradeNo: n.TransID, OrderID: n.OrderNum, Amount: amount.InexactFloat64(), Status: status, RawData: rawBody, Metadata: map[string]string{"currency": n.OrderCurrency, "payment_brand": n.PaymentBrand, "resp_code": n.RespCode}}, nil
}

func (s *Sunrate) Refund(context.Context, payment.RefundRequest) (*payment.RefundResponse, error) { return nil, fmt.Errorf("sunrate refund is not implemented") }

func (s *Sunrate) post(ctx context.Context, path string, payload map[string]string, out any) error {
	raw, _ := json.Marshal(payload); req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(s.config["apiBase"], "/")+path, bytes.NewReader(raw)); if err != nil { return err }; req.Header.Set("Accept", "application/json"); req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req); if err != nil { return err }; defer resp.Body.Close(); body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)); if err != nil { return err }; if resp.StatusCode < 200 || resp.StatusCode >= 300 { return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body)) }; if out != nil { return json.Unmarshal(body, out) }; return nil
}

func sunrateSign(fields map[string]string, key string) string { keys := make([]string, 0, len(fields)); for k := range fields { if k != "signature" { keys = append(keys, k) } }; sort.Strings(keys); var b strings.Builder; for i, k := range keys { if i > 0 { b.WriteByte('&') }; b.WriteString(k); b.WriteByte('='); b.WriteString(fields[k]) }; b.WriteString(key); sum := sha256.Sum256([]byte(b.String())); return hex.EncodeToString(sum[:]) }
func equalSignature(a,b string) bool { return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b)) }

type sunrateCreateResponse struct { RespCode string `json:"respCode"`; RespMsg string `json:"respMsg"`; OrderNum string `json:"orderNum"`; TransID string `json:"transID"`; Parameter struct { PayURL string `json:"payUrl"` } `json:"parameter"` }
type sunrateNotifyResponse struct { TransType string `json:"transType"`; OrderNum string `json:"orderNum"`; OrderAmount string `json:"orderAmount"`; OrderCurrency string `json:"orderCurrency"`; MerID string `json:"merID"`; RespCode string `json:"respCode"`; RespMsg string `json:"respMsg"`; TransID string `json:"transID"`; PaymentBrand string `json:"paymentBrand"`; TransTime string `json:"transTime"`; GwTime string `json:"gwTime"`; SignType string `json:"signType"`; Signature string `json:"signature"` }
func (n sunrateNotifyResponse) signingFields() map[string]string { m:=map[string]string{"transType":n.TransType,"orderNum":n.OrderNum,"orderAmount":n.OrderAmount,"orderCurrency":n.OrderCurrency,"merID":n.MerID,"respCode":n.RespCode,"respMsg":n.RespMsg,"transID":n.TransID,"paymentBrand":n.PaymentBrand,"transTime":n.TransTime,"gwTime":n.GwTime,"signType":n.SignType}; for k,v:=range m { if v=="" { delete(m,k) } }; return m }
