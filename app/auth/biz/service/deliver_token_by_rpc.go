package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/auth/biz/utils"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

type DeliverTokenByRPCService struct {
	ctx     context.Context
	authSvc *AuthorizationService
}

// NewDeliverTokenByRPCService new DeliverTokenByRPCService
func NewDeliverTokenByRPCService(ctx context.Context, authSvc *AuthorizationService) *DeliverTokenByRPCService {
	return &DeliverTokenByRPCService{
		ctx:     ctx,
		authSvc: authSvc,
	}
}

// Run 生成并分发令牌
func (s *DeliverTokenByRPCService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	resp = new(auth.DeliveryResp)

	// 参数校验
	if req.UserId <= 0 || req.Username == "" || req.Role == "" {
		klog.Errorf("参数无效: user_id=%d, username=%s, role=%s", req.UserId, req.Username, req.Role)
		resp.Code = auth.ErrorCode_GenerateTokenError
		resp.Message = "参数无效"
		return resp, nil
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(req.UserId, req.Username, req.Role)
	if err != nil {
		klog.Errorf("生成令牌失败: %v", err)
		resp.Code = auth.ErrorCode_GenerateTokenError
		resp.Message = "生成令牌失败"
		return resp, nil
	}

	// 返回成功响应
	resp.Code = auth.ErrorCode_Success
	resp.Message = "令牌生成成功"
	resp.Token = token

	klog.Infof("成功为用户[%s]生成令牌, role=%s", req.Username, req.Role)
	return resp, nil
}
