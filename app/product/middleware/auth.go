package middleware

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/product/infra/rpc"
	"github.com/trashwbin/dymall/app/product/utils"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

// AuthMiddleware 权限验证中间件
func AuthMiddleware(required bool, requireAdmin bool) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			// 从上下文获取token
			token := utils.GetTokenFromContext(ctx)
			if token == "" {
				if required {
					klog.CtxErrorf(ctx, "未提供认证token")
					return utils.NewBizError(40100, "请先登录")
				}
				// 不需要认证，直接放行
				return next(ctx, req, resp)
			}

			// 验证token
			verifyResp, err := rpc.AuthClient.VerifyTokenByRPC(ctx, &auth.VerifyTokenReq{
				Token: token,
			})
			if err != nil {
				klog.CtxErrorf(ctx, "验证token失败 - err: %v", err)
				return utils.NewBizError(40101, "token验证失败")
			}
			if verifyResp.Code != auth.ErrorCode_Success {
				klog.CtxErrorf(ctx, "token无效 - code: %v, message: %s", verifyResp.Code, verifyResp.Message)
				return utils.NewBizError(40102, "token无效")
			}

			// 检查角色权限
			if requireAdmin && verifyResp.Role != "admin" {
				klog.CtxErrorf(ctx, "权限不足 - userId: %d, role: %s", verifyResp.UserId, verifyResp.Role)
				return utils.NewBizError(40103, "权限不足，需要管理员权限")
			}

			// token验证通过，继续处理请求
			return next(ctx, req, resp)
		}
	}
}
