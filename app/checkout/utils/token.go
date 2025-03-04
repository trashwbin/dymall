package utils

import (
	"context"
	"fmt"

	"github.com/bytedance/gopkg/cloud/metainfo"
)

const (
	TokenKey = "token"
)

// GetTokenFromContext 从上下文中获取令牌
func GetTokenFromContext(ctx context.Context) string {
	temp, ok1 := metainfo.GetValue(ctx, "token")
	if ok1 {
		return temp
	}

	// 从context中获取
	if token, ok := ctx.Value(TokenKey).(string); ok {
		fmt.Println("从Context获取到token:", token)
		return token
	}

	fmt.Println("未找到token")
	return ""
}
