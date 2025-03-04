package service

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

func TestUpdateProduct_Run(t *testing.T) {
	tests := []struct {
		name    string
		req     *product.UpdateProductReq
		wantErr bool
		errCode int32
	}{
		{
			name: "正常更新商品",
			req: &product.UpdateProductReq{
				Id:          36,
				Name:        "更新后的商品",
				Description: "这是更新后的商品描述",
				Picture:     "http://example.com/updated.jpg",
				Price:       199.9,
				Categories:  []string{"更新分类", "电子产品"},
			},
			wantErr: false,
		},
		{
			name: "商品ID为空",
			req: &product.UpdateProductReq{
				Name:  "测试商品",
				Price: 99.9,
			},
			wantErr: true,
			errCode: 40001,
		},
		{
			name: "商品名称为空",
			req: &product.UpdateProductReq{
				Id:    36,
				Price: 99.9,
			},
			wantErr: true,
			errCode: 40002,
		},
		{
			name: "商品价格为负",
			req: &product.UpdateProductReq{
				Id:    36,
				Name:  "测试商品",
				Price: -1,
			},
			wantErr: true,
			errCode: 40003,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUpdateProductService(context.Background())
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				bizErr, ok := err.(*kerrors.BizStatusError)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, bizErr.BizStatusCode())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Product)
				assert.Equal(t, tt.req.Name, resp.Product.Name)
				assert.Equal(t, tt.req.Description, resp.Product.Description)
				assert.Equal(t, tt.req.Picture, resp.Product.Picture)
				assert.Equal(t, tt.req.Price, resp.Product.Price)
				assert.Equal(t, tt.req.Categories, resp.Product.Categories)
			}
		})
	}
}
