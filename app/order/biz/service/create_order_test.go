package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	order "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
)

func TestCreateOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCreateOrderService(ctx)

	tests := []struct {
		name    string
		req     *order.CreateOrderReq
		wantErr bool
	}{
		{
			name: "正常创建订单",
			req: &order.CreateOrderReq{
				UserId: 1001,
				Address: &order.Address{
					StreetAddress: "测试街道",
					City:          "测试城市",
					State:         "测试省份",
					Country:       "测试国家",
					ZipCode:       518000,
				},
				Email:    "test@example.com",
				ExpireAt: time.Now().Add(30 * time.Minute).Unix(),
			},
			wantErr: false,
		},
		{
			name: "购物车为空",
			req: &order.CreateOrderReq{
				UserId: 1002, // 使用不同的用户ID，模拟购物车为空的情况
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.Run(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, resp.Order.OrderId)
			assert.Equal(t, tt.req.UserId, resp.Order.UserId)
			assert.Equal(t, tt.req.Email, resp.Order.Email)
			assert.Equal(t, order.OrderStatus_ORDER_STATUS_PENDING, resp.Order.Status)
		})
	}
}
