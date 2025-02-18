package service

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/cart/biz/dal/mysql"
	cart "github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

// MockProductClient 模拟商品服务客户端
type MockProductClient struct{}

func (m *MockProductClient) GetProduct(ctx context.Context, req *product.GetProductReq, opts ...callopt.Option) (r *product.GetProductResp, err error) {
	// 模拟商品服务响应
	if req.Id == 1001 {
		return &product.GetProductResp{
			Product: &product.Product{
				Id:          1001,
				Name:        "测试商品",
				Description: "这是一个测试商品",
				Price:       99.9,
			},
		}, nil
	}
	return &product.GetProductResp{}, nil
}

func (m *MockProductClient) ListProducts(ctx context.Context, req *product.ListProductsReq, opts ...callopt.Option) (r *product.ListProductsResp, err error) {
	return nil, nil
}

func (m *MockProductClient) SearchProducts(ctx context.Context, req *product.SearchProductsReq, opts ...callopt.Option) (r *product.SearchProductsResp, err error) {
	return nil, nil
}

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
	assert.Len(t, cartItems, 1)
	assert.Equal(t, int32(5), cartItems[0].Quantity) // 2 + 3
}
