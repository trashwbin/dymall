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
	if req.UserId <= 0 {
		klog.Errorf("参数无效: user_id=%d", req.UserId)
		resp.Code = auth.ErrorCode_GenerateTokenError
		resp.Message = "用户ID无效"
		return resp, nil
	}

	// 获取用户角色
	roles, err := s.authSvc.GetRolesForUser(req.UserId)
	if err != nil {
		klog.Errorf("获取用户角色失败: %v", err)
		resp.Code = auth.ErrorCode_GenerateTokenError
		resp.Message = "获取用户角色失败"
		return resp, nil
	}

	// 如果用户没有任何角色，默认赋予guest角色
	role := "guest"
	if len(roles) > 0 {
		role = roles[0] // 使用第一个角色
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(req.UserId, role)

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
	resp.Role = role
	klog.Infof("成功为用户[%d]生成令牌, role=%s", req.UserId, role)
	return resp, nil
}
