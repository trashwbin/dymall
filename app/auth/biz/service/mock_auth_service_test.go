package service

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

// NewMockAuthorizationService 创建用于测试的授权服务
func NewMockAuthorizationService() (*AuthorizationService, error) {
	// 使用内存适配器
	enforcer, err := casbin.NewEnforcer("../../resources/rbac_model.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %v", err)
	}

	return &AuthorizationService{
		enforcer: enforcer,
	}, nil
}
