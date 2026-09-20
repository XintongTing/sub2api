//go:build unit

package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/stretchr/testify/require"
)

func TestSunrateCreatePaymentUsesCashierPurchaseRequest(t *testing.T) {
	t.Parallel()

	var payload map[string]string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v5/createcashier", r.URL.Path)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		_, _ = w.Write([]byte(`{"respCode":"00","parameter":{"payUrl":"https://cashier.example.com/pay/session"}}`))
	}))
	defer server.Close()

	provider, err := NewSunrate("1", map[string]string{
		"merchantId": "merchant_123", "signatureKey": "secret", "apiBase": server.URL,
		"filters": "truemoney,rabbit_line_pay,kplus,promptpay",
	})
	require.NoError(t, err)
	provider.httpClient = server.Client()

	resp, err := provider.CreatePayment(context.Background(), payment.CreatePaymentRequest{
		OrderID: "sub2_order", Amount: "100", Subject: "Wallet top-up",
		NotifyURL: "https://merchant.example.com/api/v1/payment/webhook/sunrate",
		ReturnURL: "https://merchant.example.com/payment/result", ClientIP: "203.0.113.10",
	})
	require.NoError(t, err)
	require.Equal(t, "https://cashier.example.com/pay/session", resp.PayURL)
	require.Equal(t, "PURC", payload["transType"])
	require.Equal(t, "sunrate_cashier", payload["paymentMethod"])
	require.Equal(t, "100.00", payload["orderAmount"])
	require.Equal(t, "THB", payload["orderCurrency"])
	require.Equal(t, "truemoney,rabbit_line_pay,kplus,promptpay", payload["filters"])
	require.NotEmpty(t, payload["signature"])
}
