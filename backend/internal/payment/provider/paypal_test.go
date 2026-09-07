//go:build unit

package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/stretchr/testify/require"
)

func TestNewPayPalDefaultsToSandboxTHB(t *testing.T) {
	t.Parallel()

	prov, err := NewPayPal("1", map[string]string{
		"clientId":     "cid",
		"clientSecret": "secret",
		"webhookId":    "wh_123",
	})
	require.NoError(t, err)
	require.Equal(t, payment.TypePayPal, prov.ProviderKey())
	require.Equal(t, []payment.PaymentType{payment.TypePayPal}, prov.SupportedTypes())
	require.Equal(t, "THB", prov.config["currency"])
	require.Equal(t, "sandbox", prov.config["environment"])
	require.Equal(t, paypalSandboxAPIBase, prov.config["apiBase"])
}

func TestPayPalCreatePaymentUsesHostedRedirect(t *testing.T) {
	t.Parallel()

	var payload paypalCreateOrderRequest
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/oauth2/token":
			user, pass, ok := r.BasicAuth()
			require.True(t, ok)
			require.Equal(t, "cid", user)
			require.Equal(t, "secret", pass)
			_, _ = w.Write([]byte(`{"access_token":"token","expires_in":3600}`))
		case "/v2/checkout/orders":
			require.Equal(t, "Bearer token", r.Header.Get("Authorization"))
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(body, &payload))
			_, _ = w.Write([]byte(`{"id":"order_123","status":"CREATED","links":[{"rel":"approve","href":"https://www.paypal.com/checkoutnow?token=order_123"}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	prov := mustTestPayPalProvider(t, server)
	resp, err := prov.CreatePayment(context.Background(), payment.CreatePaymentRequest{
		OrderID:   "sub2_order",
		Amount:    "199",
		Subject:   "OneAPI 199 THB",
		NotifyURL: "https://merchant.example.com/api/v1/payment/webhook/paypal",
		ReturnURL: "https://merchant.example.com/payment/result",
	})
	require.NoError(t, err)
	require.Equal(t, "order_123", resp.TradeNo)
	require.Equal(t, "https://www.paypal.com/checkoutnow?token=order_123", resp.PayURL)
	require.Equal(t, "THB", resp.Currency)
	require.Equal(t, payment.CreatePaymentResultOrderCreated, resp.ResultType)
	require.Equal(t, "CAPTURE", payload.Intent)
	require.Equal(t, "sub2_order", payload.PurchaseUnits[0].CustomID)
	require.Equal(t, "THB", payload.PurchaseUnits[0].Amount.CurrencyCode)
	require.Equal(t, "199.00", payload.PurchaseUnits[0].Amount.Value)
}

func TestPayPalQueryOrderCapturesApprovedOrder(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/oauth2/token":
			_, _ = w.Write([]byte(`{"access_token":"token","expires_in":3600}`))
		case "/v2/checkout/orders/order_123":
			_, _ = w.Write([]byte(`{"id":"order_123","status":"APPROVED","purchase_units":[{"custom_id":"sub2_order","amount":{"currency_code":"THB","value":"699.00"}}]}`))
		case "/v2/checkout/orders/order_123/capture":
			_, _ = w.Write([]byte(`{"id":"order_123","status":"COMPLETED","purchase_units":[{"custom_id":"sub2_order","payments":{"captures":[{"id":"cap_123","status":"COMPLETED","amount":{"currency_code":"THB","value":"699.00"}}]}}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	prov := mustTestPayPalProvider(t, server)
	resp, err := prov.QueryOrder(context.Background(), "order_123")
	require.NoError(t, err)
	require.Equal(t, "order_123", resp.TradeNo)
	require.Equal(t, payment.ProviderStatusPaid, resp.Status)
	require.InDelta(t, 699, resp.Amount, 0.0001)
	require.Equal(t, "THB", resp.Metadata["currency"])
	require.Equal(t, "cap_123", resp.Metadata["capture_id"])
	require.Equal(t, "sub2_order", resp.Metadata["order_id"])
}

func TestPayPalVerifyNotificationUsesPayPalVerification(t *testing.T) {
	t.Parallel()

	raw := `{"id":"evt_123","event_type":"CHECKOUT.ORDER.APPROVED","resource":{"id":"order_123","status":"APPROVED","purchase_units":[{"custom_id":"sub2_order"}]}}`
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/oauth2/token":
			_, _ = w.Write([]byte(`{"access_token":"token","expires_in":3600}`))
		case "/v1/notifications/verify-webhook-signature":
			var payload paypalVerifyWebhookRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "wh_123", payload.WebhookID)
			require.JSONEq(t, raw, string(payload.WebhookEvent))
			_, _ = w.Write([]byte(`{"verification_status":"SUCCESS"}`))
		case "/v2/checkout/orders/order_123":
			_, _ = w.Write([]byte(`{"id":"order_123","status":"COMPLETED","purchase_units":[{"custom_id":"sub2_order","payments":{"captures":[{"id":"cap_123","status":"COMPLETED","amount":{"currency_code":"THB","value":"199.00"}}]}}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	prov := mustTestPayPalProvider(t, server)
	n, err := prov.VerifyNotification(context.Background(), raw, signedPayPalHeaders())
	require.NoError(t, err)
	require.NotNil(t, n)
	require.Equal(t, "order_123", n.TradeNo)
	require.Equal(t, "sub2_order", n.OrderID)
	require.Equal(t, payment.NotificationStatusSuccess, n.Status)
	require.InDelta(t, 199, n.Amount, 0.0001)
	require.Equal(t, "THB", n.Metadata["currency"])
}

func TestPayPalRefundCapture(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/oauth2/token":
			_, _ = w.Write([]byte(`{"access_token":"token","expires_in":3600}`))
		case "/v2/checkout/orders/order_123":
			_, _ = w.Write([]byte(`{"id":"order_123","status":"COMPLETED","purchase_units":[{"payments":{"captures":[{"id":"cap_123","status":"COMPLETED","amount":{"currency_code":"THB","value":"199.00"}}]}}]}`))
		case "/v2/payments/captures/cap_123/refund":
			var payload paypalRefundRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "THB", payload.Amount.CurrencyCode)
			require.Equal(t, "99.00", payload.Amount.Value)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"refund_123","status":"COMPLETED"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	prov := mustTestPayPalProvider(t, server)
	resp, err := prov.Refund(context.Background(), payment.RefundRequest{
		TradeNo: "order_123",
		OrderID: "sub2_order",
		Amount:  "99.00",
	})
	require.NoError(t, err)
	require.Equal(t, "refund_123", resp.RefundID)
	require.Equal(t, payment.ProviderStatusSuccess, resp.Status)
}

func mustTestPayPalProvider(t *testing.T, server *httptest.Server) *PayPal {
	t.Helper()
	prov, err := NewPayPal("1", map[string]string{
		"clientId":     "cid",
		"clientSecret": "secret",
		"webhookId":    "wh_123",
		"apiBase":      server.URL,
	})
	require.NoError(t, err)
	prov.httpClient = server.Client()
	return prov
}

func signedPayPalHeaders() map[string]string {
	return map[string]string{
		"paypal-auth-algo":         "SHA256withRSA",
		"paypal-cert-url":          "https://api-m.sandbox.paypal.com/certs/test",
		"paypal-transmission-id":   "transmission_123",
		"paypal-transmission-sig":  "signature",
		"paypal-transmission-time": "2026-06-29T00:00:00Z",
	}
}
