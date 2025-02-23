package middleware

import (
	"context"
	"os"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/trashwbin/dymall/app/auth/biz/service"

	"github.com/trashwbin/dymall/app/auth/biz/utils"
)

// 自定义上下文键类型
type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

// 白名单路径
var whiteList = map[string]bool{
	"DeliverTokenByRPC": true,
	"VerifyTokenByRPC":  true,
	"Register":          true,
	"Login":             true,
}

// 内部服务调用路径
var internalServiceMethods = map[string]bool{
	"AddRoleForUser":    true,
	"GetRolesForUser":   true,
	"RemoveRoleForUser": true,
}

// 允许的内部服务列表
var allowedInternalServices = map[string]bool{
	"user":    true, // 用户服务
	"product": true, // 商品服务
	"order":   true, // 订单服务
}

// isDevEnvironment 判断是否为开发环境
func isDevEnvironment() bool {
	return os.Getenv("GO_ENV") == "dev" || os.Getenv("GO_ENV") == "test"
}

// AuthMiddleware 认证中间件
func AuthMiddleware(authSvc *service.AuthorizationService) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			// 获取请求路径和来源服务
			ri := rpcinfo.GetRPCInfo(ctx)
			if ri == nil {
				return next(ctx, req, resp)
			}

			method := ri.To().Method()
			fromService := ri.From().ServiceName()
			klog.Infof("收到请求: method=%s, from_service=%s", method, fromService)

			// 检查是否是内部服务方法
			if internalServiceMethods[method] {
				// 在开发环境中，允许没有服务名的请求（用于测试）
				if isDevEnvironment() && fromService == "" {
					klog.Warnf("开发环境: 允许直接调用内部服务方法 method=%s", method)
					return next(ctx, req, resp)
				}

				if !allowedInternalServices[fromService] {
					klog.Errorf("非法的服务调用: method=%s, from_service=%s", method, fromService)
					return nil
				}
				return next(ctx, req, resp)
			}

			// 检查是否在白名单中
			if whiteList[method] {
				return next(ctx, req, resp)
			}

			// 获取认证令牌
			token := ""
			if md, ok := rpcinfo.GetRPCInfo(ctx).From().Tag("token"); ok {
				token = md
			}

			if token == "" {
				klog.Errorf("未提供认证令牌: method=%s", method)
				return nil
			}

			// 验证令牌
			claims, err := utils.ParseToken(token)
			if err != nil {
				klog.Errorf("令牌验证失败: %v", err)
				return err
			}

			// 只有管理员可以管理权限策略
			if method == "AddPolicy" || method == "RemovePolicy" {
				if claims.Role != "admin" {
					klog.Errorf("权限不足: 只有管理员可以管理权限策略")
					return nil
				}
			}

			// 将用户信息添加到上下文
			ctx = context.WithValue(ctx, userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, roleKey, claims.Role)

			return next(ctx, req, resp)
		}
	}
}
