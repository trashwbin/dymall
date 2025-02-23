package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

type GetRolesForUserService struct {
	ctx     context.Context
	authSvc *AuthorizationService
}

// NewGetRolesForUserService new GetRolesForUserService
func NewGetRolesForUserService(ctx context.Context, authSvc *AuthorizationService) *GetRolesForUserService {
	return &GetRolesForUserService{
		ctx:     ctx,
		authSvc: authSvc,
	}
}

// Run 获取用户的所有角色
func (s *GetRolesForUserService) Run(req *auth.RoleQueryReq) (resp *auth.RoleQueryResp, err error) {
	resp = new(auth.RoleQueryResp)

	// 参数校验
	if req.UserId <= 0 {
		klog.Errorf("参数无效: user_id=%d", req.UserId)
		resp.Code = auth.ErrorCode_PermissionDenied
		resp.Message = "用户ID无效"
		return resp, nil
	}

	// 获取用户角色列表
	roles, err := s.authSvc.GetRolesForUser(req.UserId)
	if err != nil {
		klog.Errorf("获取用户角色失败: %v", err)
		resp.Code = auth.ErrorCode_PermissionDenied
		resp.Message = err.Error()
		return resp, nil
	}

	resp.Code = auth.ErrorCode_Success
	resp.Message = "获取用户角色成功"
	resp.Roles = roles

	klog.Infof("成功获取用户[%d]的角色列表: %v", req.UserId, roles)
	return resp, nil
}
