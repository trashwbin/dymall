package service

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

func TestListProducts_Run(t *testing.T) {
	tests := []struct {
		name    string
		req     *product.ListProductsReq
		wantErr bool
		errCode int32
	}{
		{
			name: "正常获取商品列表",
			req: &product.ListProductsReq{
				Page:     1,
				PageSize: 10,
			},
			wantErr: false,
		},
		{
			name: "按分类获取商品列表",
			req: &product.ListProductsReq{
				Page:         1,
				PageSize:     10,
				CategoryName: "电子产品",
			},
			wantErr: false,
		},
		{
			name: "页码为0自动修正为1",
			req: &product.ListProductsReq{
				Page:     0,
				PageSize: 10,
			},
			wantErr: false,
		},
		{
			name: "每页数量为0自动修正为10",
			req: &product.ListProductsReq{
				Page:     1,
				PageSize: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewListProductsService(context.Background())
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				bizErr, ok := err.(*kerrors.BizStatusError)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, bizErr.BizStatusCode())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Products)
			}
		})
	}
}
