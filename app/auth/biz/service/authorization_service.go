package service

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/auth/conf"
	"gorm.io/gorm"
)

type AuthorizationService struct {
	enforcer *casbin.Enforcer
}

var authService *AuthorizationService

// NewAuthorizationService 创建授权服务单例
func NewAuthorizationService(db *gorm.DB) (*AuthorizationService, error) {
	if authService != nil {
		return authService, nil
	}

	// 初始化 Gorm 适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		klog.Errorf("初始化Casbin Adapter失败: %v", err)
		return nil, err
	}

	// 获取配置
	config := conf.GetConf()

	// 创建enforcer
	enforcer, err := casbin.NewEnforcer(config.Casbin.ModelPath, adapter)
	if err != nil {
		klog.Errorf("初始化Casbin Enforcer失败: %v", err)
		return nil, err
	}

	authService = &AuthorizationService{
		enforcer: enforcer,
	}

	// 初始化默认权限
	if err := authService.InitializeDefaultPolicies(); err != nil {
		klog.Errorf("初始化默认权限失败: %v", err)
		return nil, err
	}

	return authService, nil
}

// InitializeDefaultPolicies 初始化默认的权限策略
func (s *AuthorizationService) InitializeDefaultPolicies() error {
	// 定义服务和操作
	policies := []struct {
		Role     string
		Resource string
		Action   string
	}{
		// 管理员权限
		{"admin", "*", "*"}, // 管理员可以访问所有资源的所有操作

		// 用户服务权限
		{"user", "user", "get"},     // 获取用户信息
		{"user", "user", "update"},  // 更新用户信息
		{"guest", "user", "create"}, // 创建用户（注册）
		{"guest", "user", "login"},  // 登录

		// 商品服务权限
		{"user", "product", "get"},     // 查询商品
		{"user", "product", "list"},    // 批量查询商品
		{"admin", "product", "create"}, // 创建商品
		{"admin", "product", "update"}, // 更新商品
		{"admin", "product", "delete"}, // 删除商品

		// 购物车服务权限
		{"user", "cart", "create"}, // 创建购物车
		{"user", "cart", "clear"},  // 清空购物车
		{"user", "cart", "get"},    // 获取购物车信息

		// 订单服务权限
		{"user", "order", "create"},  // 创建订单
		{"user", "order", "update"},  // 更新订单
		{"user", "order", "get"},     // 获取订单信息
		{"admin", "order", "cancel"}, // 取消订单

		// 支付服务权限
		{"user", "payment", "pay"},    // 支付
		{"user", "payment", "cancel"}, // 取消支付
	}

	// 添加权限策略
	for _, p := range policies {
		_, err := s.enforcer.AddPolicy(p.Role, p.Resource, p.Action)
		if err != nil {
			klog.Errorf("添加权限策略失败: role=%s, resource=%s, action=%s, err=%v",
				p.Role, p.Resource, p.Action, err)
			return err
		}
	}

	klog.Info("成功初始化默认权限策略")
	return nil
}

// CheckPermission 检查权限
func (s *AuthorizationService) CheckPermission(sub, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, obj, act)
}

// AddPolicy 添加策略
func (s *AuthorizationService) AddPolicy(sub, obj, act string) (bool, error) {
	return s.enforcer.AddPolicy(sub, obj, act)
}

// RemovePolicy 删除策略
func (s *AuthorizationService) RemovePolicy(sub, obj, act string) (bool, error) {
	return s.enforcer.RemovePolicy(sub, obj, act)
}

// AddRoleForUser 为用户添加角色
func (s *AuthorizationService) AddRoleForUser(userID int64, role string) (bool, error) {
	return s.enforcer.AddGroupingPolicy(fmt.Sprintf("user:%d", userID), role)
}

// RemoveRoleForUser 删除用户的角色
func (s *AuthorizationService) RemoveRoleForUser(userID int64, role string) (bool, error) {
	return s.enforcer.RemoveGroupingPolicy(fmt.Sprintf("user:%d", userID), role)
}

// GetRolesForUser 获取用户的所有角色
func (s *AuthorizationService) GetRolesForUser(userID int64) ([]string, error) {
	return s.enforcer.GetRolesForUser(fmt.Sprintf("user:%d", userID))
}
