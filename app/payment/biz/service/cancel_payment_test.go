package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/payment/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/payment/biz/model"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

func TestCancelPayment_Run(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		req     *payment.CancelPaymentReq
		wantErr bool
	}{
		{
			name: "normal case",
			setup: func() {
				testPayment.Status = model.PayStatusPending
				_ = mysql.NewPaymentRepo().UpdatePayment(testPayment)
			},
			req: &payment.CancelPaymentReq{
				PaymentId: testPayment.PaymentID,
				UserId:    uint32(testPayment.UserID),
			},
			wantErr: false,
		},
		{
			name:  "payment not found",
			setup: func() {},
			req: &payment.CancelPaymentReq{
				PaymentId: "non_existent_payment",
				UserId:    1001,
			},
			wantErr: true,
		},
		{
			name:  "user id not match",
			setup: func() {},
			req: &payment.CancelPaymentReq{
				PaymentId: testPayment.PaymentID,
				UserId:    9999,
			},
			wantErr: true,
		},
		{
			name: "payment already paid",
			setup: func() {
				testPayment.Status = model.PayStatusSuccess
				_ = mysql.NewPaymentRepo().UpdatePayment(testPayment)
			},
			req: &payment.CancelPaymentReq{
				PaymentId: testPayment.PaymentID,
				UserId:    uint32(testPayment.UserID),
			},
			wantErr: true,
		},
	}

	ctx := context.Background()
	s := NewCancelPaymentService(ctx)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := s.Run(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)

			// 验证支付单状态
			payment, err := mysql.NewPaymentRepo().GetPaymentByID(tt.req.PaymentId)
			assert.NoError(t, err)
			assert.Equal(t, model.PayStatusCancelled, payment.Status)
		})
	}
}
