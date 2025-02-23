package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

type AddRoleForUserService struct {
	ctx     context.Context
	authSvc *AuthorizationService
}

// NewAddRoleForUserService new AddRoleForUserService
func NewAddRoleForUserService(ctx context.Context, authSvc *AuthorizationService) *AddRoleForUserService {
	return &AddRoleForUserService{
		ctx:     ctx,
		authSvc: authSvc,
	}
}

// Run 为用户添加角色
func (s *AddRoleForUserService) Run(req *auth.RoleBindingReq) (resp *auth.PolicyResp, err error) {
	resp = new(auth.PolicyResp)

	// 参数校验
	if req.UserId <= 0 || req.Role == "" {
		klog.Errorf("参数无效: user_id=%d, role=%s", req.UserId, req.Role)
		resp.Code = auth.ErrorCode_PermissionDenied
		resp.Message = "参数无效"
		resp.Success = false
		return resp, nil
	}

	// 为用户添加角色
	success, err := s.authSvc.AddRoleForUser(req.UserId, req.Role)
	if err != nil {
		klog.Errorf("为用户添加角色失败: %v", err)
		resp.Code = auth.ErrorCode_PermissionDenied
		resp.Message = err.Error()
		resp.Success = false
		return resp, nil
	}

	resp.Code = auth.ErrorCode_Success
	resp.Message = "为用户添加角色成功"
	resp.Success = success

	klog.Infof("成功为用户[%d]添加角色: %s", req.UserId, req.Role)
	return resp, nil
}
