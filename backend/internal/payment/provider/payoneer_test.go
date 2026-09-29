//go:build unit

package provider

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/stretchr/testify/require"
)

func TestNewPayoneerDefaultsToSandboxTHB(t *testing.T) {
	t.Parallel()

	prov, err := NewPayoneer("1", map[string]string{
		"clientId":      "cid",
		"clientSecret":  "secret",
		"webhookSecret": "whsec",
	})
	require.NoError(t, err)
	require.Equal(t, payment.TypePayoneer, prov.ProviderKey())
	require.Equal(t, []payment.PaymentType{payment.TypePayoneer}, prov.SupportedTypes())
	require.Equal(t, "THB", prov.config["currency"])
	require.Equal(t, "sandbox", prov.config["environment"])
	require.Equal(t, payoneerSandboxAPIBase, prov.config["apiBase"])
}

func TestPayoneerCreatePaymentUsesHostedRedirect(t *testing.T) {
	t.Parallel()

	var payload map[string]any
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/checkout/payment", r.URL.Path)
		user, pass, ok := r.BasicAuth()
		require.True(t, ok)
		require.Equal(t, "cid", user)
		require.Equal(t, "secret", pass)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(body, &payload))
		_, _ = w.Write([]byte(`{"id":"pay_123","redirect_url":"https://checkout.payoneer.example/session/pay_123","amount":199,"currency":"THB","status":"PENDING"}`))
	}))
	defer server.Close()

	prov := mustTestPayoneerProvider(t, server)
	resp, err := prov.CreatePayment(context.Background(), payment.CreatePaymentRequest{
		OrderID:   "sub2_order",
		Amount:    "199",
		Subject:   "OneAPI 199 THB",
		NotifyURL: "https://merchant.example.com/api/v1/payment/webhook/payoneer",
		ReturnURL: "https://merchant.example.com/payment/result",
	})
	require.NoError(t, err)
	require.Equal(t, "pay_123", resp.TradeNo)
	require.Equal(t, "https://checkout.payoneer.example/session/pay_123", resp.PayURL)
	require.Equal(t, "THB", resp.Currency)
	require.Equal(t, payment.CreatePaymentResultOrderCreated, resp.ResultType)
	require.Equal(t, "sub2_order", payload["merchant_reference"])
	require.Equal(t, "THB", payload["currency"])
}

func TestPayoneerVerifyNotificationRequiresValidSignature(t *testing.T) {
	t.Parallel()

	prov, err := NewPayoneer("1", map[string]string{
		"clientId":      "cid",
		"clientSecret":  "secret",
		"webhookSecret": "whsec",
	})
	require.NoError(t, err)

	raw := `{"id":"pay_123","merchant_reference":"sub2_abc","amount":199,"currency":"THB","status":"COMPLETED"}`
	timestamp := time.Now().UTC().Format(time.RFC3339)
	headers := signedPayoneerHeaders(raw, timestamp, "whsec")

	n, err := prov.VerifyNotification(context.Background(), raw, headers)
	require.NoError(t, err)
	require.NotNil(t, n)
	require.Equal(t, "pay_123", n.TradeNo)
	require.Equal(t, "sub2_abc", n.OrderID)
	require.Equal(t, payment.NotificationStatusSuccess, n.Status)
	require.InDelta(t, 199, n.Amount, 0.0001)
	require.Equal(t, "THB", n.Metadata["currency"])

	headers["x-payoneer-signature"] = strings.Repeat("0", 64)
	_, err = prov.VerifyNotification(context.Background(), raw, headers)
	require.ErrorContains(t, err, "invalid signature")
}

func TestPayoneerQueryOrderMapsCompleted(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/checkout/payment/pay_123", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"pay_123","amount":699,"currency":"THB","status":"COMPLETED"}`))
	}))
	defer server.Close()

	prov := mustTestPayoneerProvider(t, server)
	resp, err := prov.QueryOrder(context.Background(), "pay_123")
	require.NoError(t, err)
	require.Equal(t, "pay_123", resp.TradeNo)
	require.Equal(t, payment.ProviderStatusPaid, resp.Status)
	require.InDelta(t, 699, resp.Amount, 0.0001)
	require.Equal(t, "THB", resp.Metadata["currency"])
}

func mustTestPayoneerProvider(t *testing.T, server *httptest.Server) *Payoneer {
	t.Helper()
	prov, err := NewPayoneer("1", map[string]string{
		"clientId":      "cid",
		"clientSecret":  "secret",
		"webhookSecret": "whsec",
		"apiBase":       server.URL,
	})
	require.NoError(t, err)
	prov.httpClient = server.Client()
	return prov
}

func signedPayoneerHeaders(rawBody, timestamp, secret string) map[string]string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestamp + "." + rawBody))
	return map[string]string{
		"x-payoneer-timestamp": timestamp,
		"x-payoneer-signature": hex.EncodeToString(mac.Sum(nil)),
	}
}
