package middleware

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/checkout/infra/rpc"
	"github.com/trashwbin/dymall/app/checkout/utils"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

// UserAuthMiddleware 用户认证中间件
func UserAuthMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			// 1. 获取token
			token := utils.GetTokenFromContext(ctx)
			// if token == "" {
			// 	return utils.NewBizError(401, "missing token")
			// }

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
