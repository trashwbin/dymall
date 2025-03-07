package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

type AddPolicyService struct {
	ctx     context.Context
	authSvc *AuthorizationService
}

// NewAddPolicyService new AddPolicyService
func NewAddPolicyService(ctx context.Context, authSvc *AuthorizationService) *AddPolicyService {
	return &AddPolicyService{
		ctx:     ctx,
		authSvc: authSvc,
	}
}

// Run 添加权限策略
func (s *AddPolicyService) Run(req *auth.PolicyReq) (resp *auth.PolicyResp, err error) {
	resp = new(auth.PolicyResp)

	// 参数校验
	if req.Role == "" || req.Resource == "" || req.Action == "" {
		klog.Errorf("参数无效: role=%s, resource=%s, action=%s", req.Role, req.Resource, req.Action)
		resp.Code = auth.ErrorCode_PermissionDenied
		resp.Message = "参数无效"
		resp.Success = false
		return resp, nil
	}

	// 添加权限策略
	success, err := s.authSvc.AddPolicy(req.Role, req.Resource, req.Action)
	if err != nil {
		klog.Errorf("添加权限策略失败: %v", err)
		resp.Code = auth.ErrorCode_PermissionDenied
		resp.Message = err.Error()
		resp.Success = false
		return resp, nil
	}

	resp.Code = auth.ErrorCode_Success
	resp.Message = "添加权限策略成功"
	resp.Success = success

	klog.Infof("成功添加权限策略: role=%s, resource=%s, action=%s", req.Role, req.Resource, req.Action)
	return resp, nil
}
