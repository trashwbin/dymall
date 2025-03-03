package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	TokenKey = "token"
)

// GetTokenFromMetadata 从metadata中获取token
func GetTokenFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	// 从metadata中获取token值
	tokens := md.Get(TokenKey)
	if len(tokens) > 0 {
		return tokens[0]
	}
	return ""
}
