package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	order "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
)

func TestUpdateOrder_Run(t *testing.T) {
	// 先创建一个订单
	ctx := context.Background()
	createSvc := NewCreateOrderService(ctx)
	createResp, err := createSvc.Run(&order.CreateOrderReq{
		UserId: 1001,
		Address: &order.Address{
			StreetAddress: "测试街道",
			City:          "测试城市",
			State:         "测试省份",
			Country:       "测试国家",
			ZipCode:       518000,
		},
		Email: "test@example.com",
	})
	assert.NoError(t, err)
	assert.NotNil(t, createResp)
	orderID := createResp.Order.OrderId

	s := NewUpdateOrderService(ctx)
	tests := []struct {
		name    string
		req     *order.UpdateOrderReq
		wantErr bool
	}{
		{
			name: "正常更新订单地址",
			req: &order.UpdateOrderReq{
				OrderId: orderID,
				UserId:  1001,
				Address: &order.Address{
					StreetAddress: "新的测试街道",
					City:          "新的测试城市",
					State:         "新的测试省份",
					Country:       "新的测试国家",
					ZipCode:       518001,
				},
			},
			wantErr: false,
		},
		{
			name: "订单不存在",
			req: &order.UpdateOrderReq{
				OrderId: "not_exist_order",
				UserId:  1001,
				Address: &order.Address{},
			},
			wantErr: true,
		},
		{
			name: "无权更新订单",
			req: &order.UpdateOrderReq{
				OrderId: orderID,
				UserId:  1002, // 不同的用户
				Address: &order.Address{},
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
			assert.Equal(t, tt.req.OrderId, resp.Order.OrderId)
			assert.Equal(t, tt.req.UserId, resp.Order.UserId)
			assert.Equal(t, tt.req.Address.StreetAddress, resp.Order.Address.StreetAddress)
			assert.Equal(t, tt.req.Address.City, resp.Order.Address.City)
			assert.Equal(t, tt.req.Address.State, resp.Order.Address.State)
			assert.Equal(t, tt.req.Address.Country, resp.Order.Address.Country)
			assert.Equal(t, tt.req.Address.ZipCode, resp.Order.Address.ZipCode)
		})
	}
}
