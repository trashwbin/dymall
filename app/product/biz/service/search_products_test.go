package service

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

func TestSearchProducts_Run(t *testing.T) {
	tests := []struct {
		name    string
		req     *product.SearchProductsReq
		wantErr bool
		errCode int32
	}{
		{
			name: "正常搜索商品",
			req: &product.SearchProductsReq{
				Query: "测试",
			},
			wantErr: false,
		},
		{
			name: "搜索关键词为空",
			req: &product.SearchProductsReq{
				Query: "",
			},
			wantErr: true,
			errCode: 40001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSearchProductsService(context.Background())
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				bizErr, ok := err.(*kerrors.BizStatusError)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, bizErr.BizStatusCode())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Results)
			}
		})
	}
}
