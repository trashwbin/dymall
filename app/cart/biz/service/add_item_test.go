package service

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/cart/biz/dal/mysql"
	cart "github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
)

func TestAddItem_Run(t *testing.T) {
	// 初始化测试环境
	// 清理测试数据
	mysql.DB.Exec("DELETE FROM carts WHERE user_id = ?", 1001)
	mysql.DB.Exec("DELETE FROM cart_items WHERE user_id = ?", 1001)

	ctx := context.Background()
	s := NewAddItemService(ctx)

	tests := []struct {
		name    string
		req     *cart.AddItemReq
		wantErr bool
		errCode int32
	}{
		{
			name: "正常添加商品",
			req: &cart.AddItemReq{
				UserId: 1001,
				Item: &cart.CartItem{
					ProductId: 1001,
					Quantity:  2,
				},
			},
			wantErr: false,
		},
		{
			name: "正常添加商品2",
			req: &cart.AddItemReq{
				UserId: 1001,
				Item: &cart.CartItem{
					ProductId: 1002,
					Quantity:  3,
				},
			},
			wantErr: false,
		},
		{
			name: "商品不存在",
			req: &cart.AddItemReq{
				UserId: 1001,
				Item: &cart.CartItem{
					ProductId: 9999,
					Quantity:  1,
				},
			},
			wantErr: true,
			errCode: 40004,
		},
		{
			name: "添加到已有商品",
			req: &cart.AddItemReq{
				UserId: 1001,
				Item: &cart.CartItem{
					ProductId: 1001,
					Quantity:  3,
				},
			},
			wantErr: false,
		},
	}

	// 遍历测试用例
	for _, tt := range tests {
		// 为每个测试用例创建一个子测试
		t.Run(tt.name, func(t *testing.T) {
			// 执行服务的Run方法并获取响应和错误
			resp, err := s.Run(tt.req)
			// 根据预期的错误状态进行断言
			if tt.wantErr {
				// 如果期望有错误，断言错误存在
				assert.Error(t, err)
				// 检查错误码
				if tt.errCode > 0 {
					// 如果期望特定的错误码，断言错误类型并检查错误码
					bizErr, ok := err.(*kerrors.BizStatusError)
					assert.True(t, ok)
					assert.Equal(t, tt.errCode, bizErr.BizStatusCode())
				}
			} else {
				// 如果不期望有错误，断言无错误且响应非空
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}

	// 验证最终结果
	var cartItems []mysql.CartItemDO
	err := mysql.DB.Where("user_id = ?", 1001).Find(&cartItems).Error
	assert.NoError(t, err)
	assert.Len(t, cartItems, 2)
	assert.Equal(t, int32(5), cartItems[0].Quantity) // 2 + 3
	assert.Equal(t, int32(3), cartItems[1].Quantity) // 3
}
