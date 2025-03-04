package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/trashwbin/dymall/app/scheduler/infra/rpc"
	"github.com/trashwbin/dymall/app/scheduler/utils"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

const (
	// 服务调用者名称
	OrderService   = "order"
	PaymentService = "payment"
)

// UserAuthMiddleware 用户认证中间件
func UserAuthMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			// 1. 获取token
			token := utils.GetTokenFromContext(ctx)
			if token == "" {
				return utils.NewBizError(401, "missing token")
			}

			// 2. 验证token
			verifyResp, err := rpc.AuthClient.VerifyTokenByRPC(ctx, &auth.VerifyTokenReq{
				Token: token,
			})
			if err != nil {
				klog.CtxErrorf(ctx, "verify token failed: %v", err)
				return utils.NewBizError(401, "invalid token")
			}

			if !verifyResp.IsValid || verifyResp.Code != auth.ErrorCode_Success {
				return utils.NewBizError(401, "invalid token")
			}

			// 3. 调用下一个中间件或处理函数
			return next(ctx, req, resp)
		}
	}
}

// ServiceAuthMiddleware 服务间调用认证中间件
func ServiceAuthMiddleware(allowedServices ...string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			// 3. 调用下一个中间件或处理函数
			return next(ctx, req, resp)
			// 1. 获取调用方信息
			ri := rpcinfo.GetRPCInfo(ctx)
			if ri == nil {
				return utils.NewBizError(403, "unauthorized")
			}

			// 2. 验证调用方服务名称
			caller := ri.From().ServiceName()
			isAllowed := false
			for _, service := range allowedServices {
				if strings.EqualFold(caller, service) {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				klog.CtxErrorf(ctx, "unauthorized service call from: %s", caller)
				return utils.NewBizError(403, "unauthorized service")
			}

			// 3. 调用下一个中间件或处理函数
			return next(ctx, req, resp)
		}
	}
}
