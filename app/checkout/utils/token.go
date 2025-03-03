package utils

import (
	"context"
)

const (
	TokenKey = "token"
)

// GetTokenFromContext 从上下文中获取令牌
func GetTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(TokenKey).(string); ok {
		return token
	}
	return ""
}
