package main

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/auth/biz/service"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct {
	authSvc *service.AuthorizationService
}

// NewAuthServiceImpl creates a new AuthServiceImpl with AuthorizationService
func NewAuthServiceImpl(authSvc *service.AuthorizationService) *AuthServiceImpl {
	return &AuthServiceImpl{authSvc: authSvc}
}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	klog.Infof("收到令牌分发请求: user_id=%d, username=%s, role=%s", req.UserId, req.Username, req.Role)
	resp, err = service.NewDeliverTokenByRPCService(ctx, s.authSvc).Run(req)
	if err != nil {
		klog.Errorf("令牌分发失败: %v", err)
	}
	return resp, err
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	klog.Infof("收到令牌验证请求: token=%s", req.Token)
	resp, err = service.NewVerifyTokenByRPCService(ctx, s.authSvc).Run(req)
	if err != nil {
		klog.Errorf("令牌验证失败: %v", err)
	}
	return resp, err
}

// AddPolicy implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) AddPolicy(ctx context.Context, req *auth.PolicyReq) (resp *auth.PolicyResp, err error) {
	return service.NewAddPolicyService(ctx, s.authSvc).Run(req)
}

// RemovePolicy implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RemovePolicy(ctx context.Context, req *auth.PolicyReq) (resp *auth.PolicyResp, err error) {
	return service.NewRemovePolicyService(ctx, s.authSvc).Run(req)
}

// AddRoleForUser implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) AddRoleForUser(ctx context.Context, req *auth.RoleBindingReq) (resp *auth.PolicyResp, err error) {
	return service.NewAddRoleForUserService(ctx, s.authSvc).Run(req)
}

// RemoveRoleForUser implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RemoveRoleForUser(ctx context.Context, req *auth.RoleBindingReq) (resp *auth.PolicyResp, err error) {
	return service.NewRemoveRoleForUserService(ctx, s.authSvc).Run(req)
}

// GetRolesForUser implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) GetRolesForUser(ctx context.Context, req *auth.RoleQueryReq) (resp *auth.RoleQueryResp, err error) {
	return service.NewGetRolesForUserService(ctx, s.authSvc).Run(req)
}
