package service

import (
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

	return authService, nil
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
func (s *AuthorizationService) AddRoleForUser(user, role string) (bool, error) {
	return s.enforcer.AddGroupingPolicy(user, role)
}

// RemoveRoleForUser 删除用户的角色
func (s *AuthorizationService) RemoveRoleForUser(user, role string) (bool, error) {
	return s.enforcer.RemoveGroupingPolicy(user, role)
}

// GetRolesForUser 获取用户的所有角色
func (s *AuthorizationService) GetRolesForUser(user string) ([]string, error) {
	return s.enforcer.GetRolesForUser(user)
}
