package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

func TestRemoveRoleForUser_Run(t *testing.T) {
	ctx := context.Background()
	authSvc, err := NewMockAuthorizationService()
	assert.NoError(t, err)

	// 先添加一个测试角色
	_, err = authSvc.AddRoleForUser(1, "test_role")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		req      *auth.RoleBindingReq
		wantCode auth.ErrorCode
		wantErr  bool
	}{
		{
			name: "正常删除用户角色",
			req: &auth.RoleBindingReq{
				UserId: 1,
				Role:   "test_role",
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
			s := NewRemoveRoleForUserService(ctx, authSvc)
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.wantCode, resp.Code)

			if tt.req.UserId == 1 {
				// 验证角色已被删除
				roles, err := authSvc.GetRolesForUser(1)
				assert.NoError(t, err)
				assert.NotContains(t, roles, "test_role")
			}
		})
	}
}
