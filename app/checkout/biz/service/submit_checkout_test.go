package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/checkout/utils"
	checkout "github.com/trashwbin/dymall/rpc_gen/kitex_gen/checkout"
)

func TestSubmitCheckout_Run(t *testing.T) {
	// 创建一个通用的有效地址
	validAddress := &checkout.Address{
		StreetAddress: "测试街道",
		City:          "测试城市",
		State:         "测试省份",
		Country:       "中国",
		ZipCode:       "100000",
	}

	tests := []struct {
		name       string
		checkoutID string
		userID     uint32
		address    *checkout.Address
		email      string
		firstname  string
		lastname   string
		wantErr    error
		check      func(*testing.T, *checkout.SubmitCheckoutResp)
	}{
		{
			name:       "结算单不存在",
			checkoutID: "non-existent",
			userID:     1002,
			address:    validAddress,
			email:      "test@example.com",
			firstname:  "Test",
			lastname:   "User",
			wantErr:    utils.NewCheckoutError(utils.ErrCheckoutNotFound),
		},
		{
			name:       "正常提交结算",
			checkoutID: "test-checkout-id",
			userID:     1001,
			address:    validAddress,
			email:      "test@example.com",
			firstname:  "Test",
			lastname:   "User",
			check: func(t *testing.T, resp *checkout.SubmitCheckoutResp) {
				assert.NotEmpty(t, resp.OrderId)
				assert.NotEmpty(t, resp.PaymentId)
				assert.Equal(t, float32(399.7), resp.TotalAmount)
				assert.Equal(t, "CNY", resp.Currency)
			},
		},
		{
			name:       "无效的邮编",
			checkoutID: "test-checkout-id",
			userID:     1001,
			address: &checkout.Address{
				StreetAddress: "测试街道",
				City:          "测试城市",
				State:         "测试省份",
				Country:       "中国",
				ZipCode:       "invalid",
			},
			email:     "test@example.com",
			firstname: "Test",
			lastname:  "User",
			wantErr:   utils.NewCheckoutError(utils.ErrInvalidZipCode),
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSubmitCheckoutService(ctx)
			resp, err := s.Run(&checkout.SubmitCheckoutReq{
				CheckoutId: tt.checkoutID,
				UserId:     tt.userID,
				Address:    tt.address,
				Email:      tt.email,
				Firstname:  tt.firstname,
				Lastname:   tt.lastname,
			})

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			if tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}
