package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

func TestAddRoleForUser_Run(t *testing.T) {
	ctx := context.Background()
	authSvc, err := NewMockAuthorizationService()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		req      *auth.RoleBindingReq
		wantCode auth.ErrorCode
		wantErr  bool
	}{
		{
			name: "正常添加用户角色",
			req: &auth.RoleBindingReq{
				UserId: 1,
				Role:   "user",
			},
			wantCode: auth.ErrorCode_Success,
			wantErr:  false,
		},
		{
			name: "参数无效",
			req: &auth.RoleBindingReq{
				UserId: 0,
				Role:   "",
			},
			wantCode: auth.ErrorCode_PermissionDenied,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAddRoleForUserService(ctx, authSvc)
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
