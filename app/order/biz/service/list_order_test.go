package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/order/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/order/biz/dal/redis"
	order "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
)

func TestListOrder_Run(t *testing.T) {
	// 清理之前的测试数据
	db := mysql.DB
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 使用事务清理数据
	tx.Exec("DELETE FROM order_items WHERE order_id IN (SELECT order_id FROM orders WHERE user_id = 1001)")
	tx.Exec("DELETE FROM order_addresses WHERE order_id IN (SELECT order_id FROM orders WHERE user_id = 1001)")
	tx.Exec("DELETE FROM orders WHERE user_id = 1001")
	tx.Commit()

	// 清理Redis缓存
	redisRepo := redis.NewOrderRepo()
	redisRepo.DeleteUserOrders(context.Background(), 1001)

	// 先创建几个订单
	ctx := context.Background()
	createSvc := NewCreateOrderService(ctx)

	// 为用户1001创建2个订单
	for i := 0; i < 2; i++ {
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
	}

	s := NewListOrderService(ctx)
	tests := []struct {
		name     string
		req      *order.ListOrderReq
		wantErr  bool
		wantSize int
	}{
		{
			name: "正常获取订单列表",
			req: &order.ListOrderReq{
				UserId: 1001,
			},
			wantErr:  false,
			wantSize: 2,
		},
		{
			name: "用户没有订单",
			req: &order.ListOrderReq{
				UserId: 1002,
			},
			wantErr:  false,
			wantSize: 0,
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
			assert.Len(t, resp.Orders, tt.wantSize)

			// 验证订单属于正确的用户
			for _, order := range resp.Orders {
				assert.Equal(t, tt.req.UserId, order.UserId)
			}
		})
	}

	// 使用事务清理测试数据
	tx = db.Begin()
	tx.Exec("DELETE FROM order_items WHERE order_id IN (SELECT order_id FROM orders WHERE user_id = 1001)")
	tx.Exec("DELETE FROM order_addresses WHERE order_id IN (SELECT order_id FROM orders WHERE user_id = 1001)")
	tx.Exec("DELETE FROM orders WHERE user_id = 1001")
	tx.Commit()

	// 清理Redis缓存
	redisRepo.DeleteUserOrders(context.Background(), 1001)
}
