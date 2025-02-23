package service

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/cart/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/cart/biz/model"
	cart "github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
)

func TestEmptyCart_Run(t *testing.T) {
	// 初始化测试环境
	// 清理测试数据
	mysql.DB.Exec("DELETE FROM carts WHERE user_id = ?", 1001)
	mysql.DB.Exec("DELETE FROM cart_items WHERE user_id = ?", 1001)

	ctx := context.Background()
	addItemService := NewAddItemService(ctx)
	emptyCartService := NewEmptyCartService(ctx)

	// 1. 先添加一些商品到购物车
	addItemReqs := []*cart.AddItemReq{
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

	// 添加商品
	for _, req := range addItemReqs {
		resp, err := addItemService.Run(req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	}

	// 验证商品添加成功
	var cartItems []mysql.CartItemDO
	err := mysql.DB.Where("user_id = ?", 1001).Find(&cartItems).Error
	assert.NoError(t, err)
	assert.Len(t, cartItems, 2)

	// 2. 测试清空购物车
	tests := []struct {
		name    string
		req     *cart.EmptyCartReq
		wantErr bool
		errCode int32
	}{
		{
			name: "购物车不存在",
			req: &cart.EmptyCartReq{
				UserId: 9999,
			},
			wantErr: true,
			errCode: 40001,
		},
		{
			name: "正常清空购物车",
			req: &cart.EmptyCartReq{
				UserId: 1001,
			},
			wantErr: false,
		},
		{
			name: "重复清空购物车",
			req: &cart.EmptyCartReq{
				UserId: 1001,
			},
			wantErr: true,
			errCode: 40001, // 购物车不存在（因为状态已经是Empty）
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := emptyCartService.Run(tt.req)
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

				// 验证购物车状态
				var cartDO mysql.CartDO
				err := mysql.DB.Where("user_id = ?", tt.req.UserId).First(&cartDO).Error
				assert.NoError(t, err)
				assert.Equal(t, int(model.CartStatusEmpty), cartDO.Status)

				// 验证购物车商品已被软删除
				var count int64
				err = mysql.DB.Unscoped().Model(&mysql.CartItemDO{}).
					Where("user_id = ? AND deleted_at IS NOT NULL", tt.req.UserId).
					Count(&count).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(2), count)
			}
		})
	}
}
