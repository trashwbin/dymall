package sdk

import (
	"github.com/cloudwego/kitex/server"
	"github.com/trashwbin/dymall/app/auth/biz/middleware"
	"github.com/trashwbin/dymall/app/auth/biz/service"
)

// AuthSDK 认证SDK
type AuthSDK struct {
	authSvc *service.AuthorizationService
}

// NewAuthSDK 创建认证SDK实例
func NewAuthSDK() (*AuthSDK, error) {
	// 初始化与auth服务的连接
	// TODO: 实现auth服务客户端的初始化
	return &AuthSDK{}, nil
}

// WithAuthMiddleware 添加认证中间件
func (sdk *AuthSDK) WithAuthMiddleware() server.Option {
	return server.WithMiddleware(middleware.AuthMiddleware(sdk.authSvc))
}

// AddPolicy 添加权限策略
func (sdk *AuthSDK) AddPolicy(role, resource, action string) error {
	// TODO: 实现添加权限策略的RPC调用
	return nil
}

// RemovePolicy 删除权限策略
func (sdk *AuthSDK) RemovePolicy(role, resource, action string) error {
	// TODO: 实现删除权限策略的RPC调用
	return nil
}

// AddRoleForUser 为用户添加角色
func (sdk *AuthSDK) AddRoleForUser(userID int64, role string) error {
	// TODO: 实现添加用户角色的RPC调用
	return nil
}

// RemoveRoleForUser 删除用户角色
func (sdk *AuthSDK) RemoveRoleForUser(userID int64, role string) error {
	// TODO: 实现删除用户角色的RPC调用
	return nil
}

// GetRolesForUser 获取用户角色列表
func (sdk *AuthSDK) GetRolesForUser(userID int64) ([]string, error) {
	// TODO: 实现获取用户角色的RPC调用
	return nil, nil
}
