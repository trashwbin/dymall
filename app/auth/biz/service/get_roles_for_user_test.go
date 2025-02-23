package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

func TestGetRolesForUser_Run(t *testing.T) {
	ctx := context.Background()
	authSvc, err := NewMockAuthorizationService()
	assert.NoError(t, err)

	// 先添加一个测试角色
	_, err = authSvc.AddRoleForUser(1, "test_role")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		req      *auth.RoleQueryReq
		wantCode auth.ErrorCode
		wantErr  bool
	}{
		{
			name: "正常获取用户角色",
			req: &auth.RoleQueryReq{
				UserId: 1,
			},
			wantCode: auth.ErrorCode_Success,
			wantErr:  false,
		},
		{
			name: "参数无效",
			req: &auth.RoleQueryReq{
				UserId: 0,
			},
			wantCode: auth.ErrorCode_PermissionDenied,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewGetRolesForUserService(ctx, authSvc)
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.wantCode, resp.Code)

			if tt.req.UserId == 1 {
				assert.Contains(t, resp.Roles, "test_role")
			}
		})
	}
}
