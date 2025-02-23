package service

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/cart/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/cart/biz/dal/redis"
	"github.com/trashwbin/dymall/app/cart/biz/model"
	cart "github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
)

func TestGetCart_Run(t *testing.T) {
	// 初始化测试环境
	// 清理测试数据
	mysql.DB.Exec("DELETE FROM carts WHERE user_id = ?", 1001)
	mysql.DB.Exec("DELETE FROM cart_items WHERE user_id = ?", 1001)

	// 清理Redis缓存
	ctx := context.Background()
	redis.RedisClient.FlushDB(ctx)

	addItemService := NewAddItemService(ctx)
	getCartService := NewGetCartService(ctx)

	tests := []struct {
		name    string
		prepare func(t *testing.T)
		req     *cart.GetCartReq
		wantErr bool
		errCode int32
		check   func(t *testing.T, resp *cart.GetCartResp)
	}{
		{
			name: "购物车不存在",
			req: &cart.GetCartReq{
				UserId: 9999,
			},
			wantErr: true,
			errCode: 40001,
		},
		{
			name: "空购物车",
			prepare: func(t *testing.T) {
				// 创建一个空购物车
				cart := &model.Cart{
					UserID: 1001,
					Status: model.CartStatusNormal,
				}
				_, err := mysql.NewCartRepo().CreateCart(cart)
				assert.NoError(t, err)
			},
			req: &cart.GetCartReq{
				UserId: 1001,
			},
			wantErr: false,
			check: func(t *testing.T, resp *cart.GetCartResp) {
				assert.NotNil(t, resp.Cart)
				assert.Equal(t, uint32(1001), resp.Cart.UserId)
				assert.Empty(t, resp.Cart.Items)
			},
		},
		{
			name: "购物车有商品",
			prepare: func(t *testing.T) {
				// 添加测试商品
				addReqs := []*cart.AddItemReq{
					{
						UserId: 1001,
						Item: &cart.CartItem{
							ProductId: 1001,
							Quantity:  2,
						},
					},
					{
						UserId: 1001,
						Item: &cart.CartItem{
							ProductId: 1002,
							Quantity:  3,
						},
					},
				}

				for _, req := range addReqs {
					resp, err := addItemService.Run(req)
					assert.NoError(t, err)
					assert.NotNil(t, resp)
				}
			},
			req: &cart.GetCartReq{
				UserId: 1001,
			},
			wantErr: false,
			check: func(t *testing.T, resp *cart.GetCartResp) {
				assert.NotNil(t, resp.Cart)
				assert.Equal(t, uint32(1001), resp.Cart.UserId)
				assert.Len(t, resp.Cart.Items, 2)

				// 验证商品1
				assert.Equal(t, uint32(1001), resp.Cart.Items[0].ProductId)
				assert.Equal(t, int32(2), resp.Cart.Items[0].Quantity)

				// 验证商品2
				assert.Equal(t, uint32(1002), resp.Cart.Items[1].ProductId)
				assert.Equal(t, int32(3), resp.Cart.Items[1].Quantity)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理测试数据
			mysql.DB.Exec("DELETE FROM carts WHERE user_id = ?", 1001)
			mysql.DB.Exec("DELETE FROM cart_items WHERE user_id = ?", 1001)

			// 准备测试数据
			if tt.prepare != nil {
				tt.prepare(t)
			}

			// 执行测试
			resp, err := getCartService.Run(tt.req)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errCode > 0 {
					bizErr, ok := err.(*kerrors.BizStatusError)
					assert.True(t, ok)
					assert.Equal(t, tt.errCode, bizErr.BizStatusCode())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.check != nil {
					tt.check(t, resp)
				}
			}
		})
	}
}
