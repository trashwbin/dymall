package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/payment/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/payment/biz/model"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

func TestProcessPayment_Run(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		req     *payment.ProcessPaymentReq
		wantErr bool
	}{
		{
			name: "normal case",
			setup: func() {
				testPayment.Status = model.PayStatusPending
				_ = mysql.NewPaymentRepo().UpdatePayment(testPayment)
			},
			req: &payment.ProcessPaymentReq{
				PaymentId: testPayment.PaymentID,
				UserId:    uint32(testPayment.UserID),
				CreditCard: &payment.CreditCardInfo{
					CreditCardNumber:          "4532756279624064",
					CreditCardCvv:             123,
					CreditCardExpirationYear:  int32(time.Now().Year() + 1),
					CreditCardExpirationMonth: 12,
				},
			},
			wantErr: false,
		},
		{
			name:  "invalid payment id",
			setup: func() {},
			req: &payment.ProcessPaymentReq{
				PaymentId: "non_existent_payment",
				UserId:    1001,
				CreditCard: &payment.CreditCardInfo{
					CreditCardNumber:          "4532756279624064",
					CreditCardCvv:             123,
					CreditCardExpirationYear:  int32(time.Now().Year() + 1),
					CreditCardExpirationMonth: 12,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid credit card",
			setup: func() {
				testPayment.Status = model.PayStatusPending
				_ = mysql.NewPaymentRepo().UpdatePayment(testPayment)
			},
			req: &payment.ProcessPaymentReq{
				PaymentId: testPayment.PaymentID,
				UserId:    uint32(testPayment.UserID),
				CreditCard: &payment.CreditCardInfo{
					CreditCardNumber:          "invalid_card_number",
					CreditCardCvv:             123,
					CreditCardExpirationYear:  int32(time.Now().Year() + 1),
					CreditCardExpirationMonth: 12,
				},
			},
			wantErr: true,
		},
		{
			name: "expired credit card",
			setup: func() {
				testPayment.Status = model.PayStatusPending
				_ = mysql.NewPaymentRepo().UpdatePayment(testPayment)
			},
			req: &payment.ProcessPaymentReq{
				PaymentId: testPayment.PaymentID,
				UserId:    uint32(testPayment.UserID),
				CreditCard: &payment.CreditCardInfo{
					CreditCardNumber:          "4532756279624064",
					CreditCardCvv:             123,
					CreditCardExpirationYear:  int32(time.Now().Year() - 1),
					CreditCardExpirationMonth: 12,
				},
			},
			wantErr: true,
		},
	}

	ctx := context.Background()
	s := NewProcessPaymentService(ctx)

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
			assert.Equal(t, testPayment.PaymentID, resp.Payment.PaymentId)
			assert.Equal(t, payment.PaymentStatus_PAYMENT_STATUS_SUCCESS, resp.Payment.Status)
			assert.NotZero(t, resp.Payment.PaidAt)
		})
	}
}
