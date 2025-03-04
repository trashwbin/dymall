package service

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

func TestDeleteProduct_Run(t *testing.T) {
	tests := []struct {
		name    string
		req     *product.DeleteProductReq
		wantErr bool
		errCode int32
	}{
		{
			name: "正常删除商品",
			req: &product.DeleteProductReq{
				Id: 35,
			},
			wantErr: false,
		},
		{
			name: "删除不存在的商品",
			req: &product.DeleteProductReq{
				Id: 99999,
			},
			wantErr: true,
			errCode: 40004,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewDeleteProductService(context.Background())
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				bizErr, ok := err.(*kerrors.BizStatusError)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, bizErr.BizStatusCode())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)
			}
		})
	}
}
