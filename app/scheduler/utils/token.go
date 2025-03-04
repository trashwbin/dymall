package utils

import (
	"context"

	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/metadata"
)

const (
	TokenKey = "token"
)

// GetTokenFromContext 从上下文中获取令牌
func GetTokenFromContext(ctx context.Context) string {
	// 首先尝试从metadata中获取token
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if tokens := md.Get(TokenKey); len(tokens) > 0 {
			return tokens[0]
		}
	}

	// 如果metadata中没有找到，则从context中获取
	if token, ok := ctx.Value(TokenKey).(string); ok {
		return token
	}
	return ""
}
