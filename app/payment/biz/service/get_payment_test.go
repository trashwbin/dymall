package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

func TestGetPayment_Run(t *testing.T) {
	tests := []struct {
		name    string
		req     *payment.GetPaymentReq
		wantErr bool
	}{
		{
			name: "normal case",
			req: &payment.GetPaymentReq{
				PaymentId: testPayment.PaymentID,
				UserId:    uint32(testPayment.UserID),
			},
			wantErr: false,
		},
		{
			name: "payment not found",
			req: &payment.GetPaymentReq{
				PaymentId: "non_existent_payment",
				UserId:    1001,
			},
			wantErr: true,
		},
		{
			name: "user id not match",
			req: &payment.GetPaymentReq{
				PaymentId: testPayment.PaymentID,
				UserId:    9999,
			},
			wantErr: true,
		},
	}

	ctx := context.Background()
	s := NewGetPaymentService(ctx)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.Run(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, testPayment.PaymentID, resp.Payment.PaymentId)
			assert.Equal(t, testPayment.OrderID, resp.Payment.OrderId)
			assert.Equal(t, uint32(testPayment.UserID), resp.Payment.UserId)
			assert.Equal(t, float32(testPayment.Amount), resp.Payment.Amount)
			assert.Equal(t, testPayment.Currency, resp.Payment.Currency)
		})
	}
}
