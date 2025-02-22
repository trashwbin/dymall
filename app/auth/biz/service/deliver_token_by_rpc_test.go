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

	tests := []struct {
		name     string
		req      *auth.DeliverTokenReq
		wantCode auth.ErrorCode
		wantErr  bool
	}{
		{
			name: "正常生成令牌",
			req: &auth.DeliverTokenReq{
				UserId:   1,
				Username: "test_user",
				Role:     "user",
			},
			wantCode: auth.ErrorCode_Success,
			wantErr:  false,
		},
		{
			name: "参数无效",
			req: &auth.DeliverTokenReq{
				UserId:   0,
				Username: "",
				Role:     "",
			},
			wantCode: auth.ErrorCode_GenerateTokenError,
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

			if tt.req.UserId > 0 {
				assert.NotEmpty(t, resp.Token)
			}
		})
	}
}
