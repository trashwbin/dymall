package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/auth/biz/utils"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

func TestVerifyTokenByRPC_Run(t *testing.T) {
	ctx := context.Background()
	authSvc, err := NewMockAuthorizationService()
	assert.NoError(t, err)

	// 生成一个有效的测试令牌
	validToken, err := utils.GenerateToken(1, "test_user", "user")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		req      *auth.VerifyTokenReq
		wantCode auth.ErrorCode
		wantErr  bool
	}{
		{
			name: "正常验证令牌",
			req: &auth.VerifyTokenReq{
				Token: validToken,
			},
			wantCode: auth.ErrorCode_Success,
			wantErr:  false,
		},
		{
			name: "令牌为空",
			req: &auth.VerifyTokenReq{
				Token: "",
			},
			wantCode: auth.ErrorCode_TokenInvalid,
			wantErr:  false,
		},
		{
			name: "无效令牌",
			req: &auth.VerifyTokenReq{
				Token: "invalid_token",
			},
			wantCode: auth.ErrorCode_TokenInvalid,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewVerifyTokenByRPCService(ctx, authSvc)
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.wantCode, resp.Code)

			if tt.req.Token == validToken {
				assert.True(t, resp.IsValid)
				assert.Equal(t, int64(1), resp.UserId)
				assert.Equal(t, "test_user", resp.Username)
				assert.Equal(t, "user", resp.Role)
			}
		})
	}
}
