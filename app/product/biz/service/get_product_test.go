package service

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

func TestGetProduct_Run(t *testing.T) {
	tests := []struct {
		name    string
		req     *product.GetProductReq
		wantErr bool
		errCode int32
	}{
		{
			name: "正常获取商品",
			req: &product.GetProductReq{
				Id: 1,
			},
			wantErr: false,
		},
		{
			name: "获取不存在的商品",
			req: &product.GetProductReq{
				Id: 99999,
			},
			wantErr: true,
			errCode: 40004,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewGetProductService(context.Background())
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
				assert.Equal(t, tt.req.Id, resp.Product.Id)
			}
		})
	}
}
