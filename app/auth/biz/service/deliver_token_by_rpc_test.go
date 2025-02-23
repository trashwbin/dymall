package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

func TestDeliverTokenByRPC_Run(t *testing.T) {
	ctx := context.Background()
	authSvc, err := NewMockAuthorizationService()
	assert.NoError(t, err)

	// 为测试用户添加角色
	_, err = authSvc.AddRoleForUser(1, "user")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		req      *auth.DeliverTokenReq
		wantCode auth.ErrorCode
		wantRole string
		wantErr  bool
	}{
		{
			name: "正常生成令牌",
			req: &auth.DeliverTokenReq{
				UserId: 1,
			},
			wantCode: auth.ErrorCode_Success,
			wantRole: "user",
			wantErr:  false,
		},
		{
			name: "用户ID无效",
			req: &auth.DeliverTokenReq{
				UserId: 0,
			},
			wantCode: auth.ErrorCode_GenerateTokenError,
			wantRole: "",
			wantErr:  false,
		},
		{
			name: "用户无角色时使用guest角色",
			req: &auth.DeliverTokenReq{
				UserId: 2,
			},
			wantCode: auth.ErrorCode_Success,
			wantRole: "guest",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewDeliverTokenByRPCService(ctx, authSvc)
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.wantCode, resp.Code)

			if tt.wantCode == auth.ErrorCode_Success {
				assert.NotEmpty(t, resp.Token)
				assert.Equal(t, tt.wantRole, resp.Role)
			}
		})
	}
}
