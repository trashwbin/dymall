package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/checkout/utils"
	checkout "github.com/trashwbin/dymall/rpc_gen/kitex_gen/checkout"
)

func TestCreateCheckout_Run(t *testing.T) {
	tests := []struct {
		name    string
		userID  uint32
		wantErr error
		check   func(*testing.T, *checkout.CreateCheckoutResp)
	}{
		{
			name:    "购物车为空",
			userID:  1002,
			wantErr: utils.NewCheckoutError(utils.ErrCartEmpty),
		},
		{
			name:   "正常结算",
			userID: 1001,
			check: func(t *testing.T, resp *checkout.CreateCheckoutResp) {
				assert.NotEmpty(t, resp.CheckoutId)
				assert.Equal(t, float32(399.7), resp.TotalAmount) // (99.9 * 2) + (199.9 * 1)
				assert.Equal(t, "CNY", resp.Currency)
				assert.Len(t, resp.Items, 2)
			},
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewCreateCheckoutService(ctx)
			resp, err := s.Run(&checkout.CreateCheckoutReq{
				UserId: tt.userID,
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
