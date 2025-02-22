package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

func TestRemovePolicy_Run(t *testing.T) {
	ctx := context.Background()
	authSvc, err := NewMockAuthorizationService()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		req      *auth.PolicyReq
		wantCode auth.ErrorCode
		wantErr  bool
	}{
		{
			name: "正常删除权限策略",
			req: &auth.PolicyReq{
				Role:     "user",
				Resource: "product",
				Action:   "read",
			},
			wantCode: auth.ErrorCode_Success,
			wantErr:  false,
		},
		{
			name: "参数无效",
			req: &auth.PolicyReq{
				Role:     "",
				Resource: "",
				Action:   "",
			},
			wantCode: auth.ErrorCode_PermissionDenied,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewRemovePolicyService(ctx, authSvc)
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.wantCode, resp.Code)
		})
	}
}
