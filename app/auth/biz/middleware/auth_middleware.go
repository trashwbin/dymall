package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/trashwbin/dymall/app/auth/biz/service"

	"github.com/trashwbin/dymall/app/auth/biz/utils"
)

// 白名单路径
var whiteList = map[string]bool{
	"DeliverTokenByRPC": true,
	"VerifyTokenByRPC":  true,
	"Register":          true,
	"Login":             true,
}

// AuthMiddleware 认证中间件
func AuthMiddleware(authSvc *service.AuthorizationService) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			// 获取请求路径
			ri := rpcinfo.GetRPCInfo(ctx)
			if ri == nil {
				return next(ctx, req, resp)
			}

			method := ri.To().Method()
			klog.Infof("收到请求: method=%s", method)

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

			// 检查权限
			service := strings.Split(method, ".")[0]
			action := strings.Split(method, ".")[1]
			allowed, err := authSvc.CheckPermission(claims.Role, service, action)
			if err != nil {
				klog.Errorf("权限检查失败: %v", err)
				return err
			}

			if !allowed {
				klog.Errorf("权限不足: %s 无权访问 %s", claims.Role, method)
				return nil
			}

			// 将用户信息添加到上下文
			ctx = context.WithValue(ctx, "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)
			ctx = context.WithValue(ctx, "role", claims.Role)

			return next(ctx, req, resp)
		}
	}
}
