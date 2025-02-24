package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

func TestCreatePayment_Run(t *testing.T) {
	tests := []struct {
		name    string
		req     *payment.CreatePaymentReq
		wantErr bool
	}{
		{
			name: "normal case",
			req: &payment.CreatePaymentReq{
				OrderId:  "test_order_001",
				UserId:   1001,
				Amount:   99.9,
				Currency: "CNY",
				ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
			},
			wantErr: false,
		},
		{
			name: "invalid amount",
			req: &payment.CreatePaymentReq{
				OrderId:  "test_order_002",
				UserId:   1001,
				Amount:   -1,
				Currency: "CNY",
				ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
			},
			wantErr: true,
		},
		{
			name: "invalid expire time",
			req: &payment.CreatePaymentReq{
				OrderId:  "test_order_003",
				UserId:   1001,
				Amount:   99.9,
				Currency: "CNY",
				ExpireAt: time.Now().Add(-24 * time.Hour).Unix(),
			},
			wantErr: true,
		},
	}

	ctx := context.Background()
	s := NewCreatePaymentService(ctx)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.Run(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, resp.Payment.PaymentId)
			assert.Equal(t, tt.req.OrderId, resp.Payment.OrderId)
			assert.Equal(t, tt.req.UserId, resp.Payment.UserId)
			assert.Equal(t, tt.req.Amount, resp.Payment.Amount)
			assert.Equal(t, tt.req.Currency, resp.Payment.Currency)
			assert.Equal(t, payment.PaymentStatus_PAYMENT_STATUS_PENDING, resp.Payment.Status)
		})
	}
}
