package sdk

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/server"
	"github.com/trashwbin/dymall/app/auth/biz/middleware"
	"github.com/trashwbin/dymall/app/auth/biz/service"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth/authservice"
)

// AuthSDK 认证SDK
type AuthSDK struct {
	client  authservice.Client
	authSvc *service.AuthorizationService
}

// NewAuthSDK 创建认证SDK实例
func NewAuthSDK(opts ...client.Option) (*AuthSDK, error) {
	// 初始化与auth服务的连接
	c, err := authservice.NewClient("auth", opts...)
	if err != nil {
		return nil, err
	}
	return &AuthSDK{client: c}, nil
}

// WithAuthMiddleware 添加认证中间件
func (sdk *AuthSDK) WithAuthMiddleware() server.Option {
	return server.WithMiddleware(middleware.AuthMiddleware(sdk.authSvc))
}

// AddPolicy 添加权限策略（需要管理员token）
func (sdk *AuthSDK) AddPolicy(ctx context.Context, role, resource, action, adminToken string) error {
	req := &auth.PolicyReq{
		Role:     role,
		Resource: resource,
		Action:   action,
	}
	opt := callopt.WithTag("token", adminToken)
	_, err := sdk.client.AddPolicy(ctx, req, opt)
	return err
}

// RemovePolicy 删除权限策略（需要管理员token）
func (sdk *AuthSDK) RemovePolicy(ctx context.Context, role, resource, action, adminToken string) error {
	req := &auth.PolicyReq{
		Role:     role,
		Resource: resource,
		Action:   action,
	}
	opt := callopt.WithTag("token", adminToken)
	_, err := sdk.client.RemovePolicy(ctx, req, opt)
	return err
}

// AddRoleForUser 为用户添加角色（内部服务调用）
func (sdk *AuthSDK) AddRoleForUser(ctx context.Context, userID int64, role string) error {
	req := &auth.RoleBindingReq{
		UserId: userID,
		Role:   role,
	}
	_, err := sdk.client.AddRoleForUser(ctx, req)
	return err
}

// RemoveRoleForUser 删除用户角色（内部服务调用）
func (sdk *AuthSDK) RemoveRoleForUser(ctx context.Context, userID int64, role string) error {
	req := &auth.RoleBindingReq{
		UserId: userID,
		Role:   role,
	}
	_, err := sdk.client.RemoveRoleForUser(ctx, req)
	return err
}

// GetRolesForUser 获取用户角色列表（内部服务调用）
func (sdk *AuthSDK) GetRolesForUser(ctx context.Context, userID int64) ([]string, error) {
	req := &auth.RoleQueryReq{
		UserId: userID,
	}
	resp, err := sdk.client.GetRolesForUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Roles, nil
}
