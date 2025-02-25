package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	order "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
)

func TestMarkOrderPaid_Run(t *testing.T) {
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

	s := NewMarkOrderPaidService(ctx)
	tests := []struct {
		name    string
		req     *order.MarkOrderPaidReq
		wantErr bool
	}{
		{
			name: "正常标记订单支付完成",
			req: &order.MarkOrderPaidReq{
				OrderId: orderID,
				UserId:  1001,
			},
			wantErr: false,
		},
		{
			name: "订单不存在",
			req: &order.MarkOrderPaidReq{
				OrderId: "not_exist_order",
				UserId:  1001,
			},
			wantErr: true,
		},
		{
			name: "无权标记订单支付完成",
			req: &order.MarkOrderPaidReq{
				OrderId: orderID,
				UserId:  1002, // 不同的用户
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

			// 验证订单状态
			getResp, err := NewGetOrderService(ctx).Run(&order.GetOrderReq{
				OrderId: tt.req.OrderId,
				UserId:  tt.req.UserId,
			})
			if !tt.wantErr {
				assert.NoError(t, err)
				assert.NotNil(t, getResp)
				assert.Equal(t, order.OrderStatus_ORDER_STATUS_PAID, getResp.Order.Status)
				assert.NotZero(t, getResp.Order.UpdatedAt)
			}
		})
	}
}
