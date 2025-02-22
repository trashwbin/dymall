package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

type RemoveRoleForUserService struct {
	ctx     context.Context
	authSvc *AuthorizationService
}

// NewRemoveRoleForUserService new RemoveRoleForUserService
func NewRemoveRoleForUserService(ctx context.Context, authSvc *AuthorizationService) *RemoveRoleForUserService {
	return &RemoveRoleForUserService{
		ctx:     ctx,
		authSvc: authSvc,
	}
}

// Run 删除用户角色
func (s *RemoveRoleForUserService) Run(req *auth.RoleBindingReq) (resp *auth.PolicyResp, err error) {
	resp = new(auth.PolicyResp)

	// 参数校验
	if req.UserId <= 0 || req.Role == "" {
		klog.Errorf("参数无效: user_id=%d, role=%s", req.UserId, req.Role)
		resp.Code = auth.ErrorCode_PermissionDenied
		resp.Message = "参数无效"
		resp.Success = false
		return resp, nil
	}

	// 删除用户角色
	success, err := s.authSvc.RemoveRoleForUser(req.UserId, req.Role)
	if err != nil {
		klog.Errorf("删除用户角色失败: %v", err)
		resp.Code = auth.ErrorCode_PermissionDenied
		resp.Message = err.Error()
		resp.Success = false
		return resp, nil
	}

	resp.Code = auth.ErrorCode_Success
	resp.Message = "删除用户角色成功"
	resp.Success = success

	klog.Infof("成功删除用户[%d]的角色: %s", req.UserId, req.Role)
	return resp, nil
}
