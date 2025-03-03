package middleware

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/product/infra/rpc"
	"github.com/trashwbin/dymall/app/product/utils"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

// RequireAdmin 检查管理员权限
func RequireAdmin(ctx context.Context) error {
	return verifyAuth(ctx, true, true)
}

// RequireUser 检查用户权限
func RequireUser(ctx context.Context) error {
	return verifyAuth(ctx, true, false)
}

// verifyAuth 验证权限
func verifyAuth(ctx context.Context, required bool, requireAdmin bool) error {
	// 从上下文获取token
	token := utils.GetTokenFromContext(ctx)
	if token == "" {
		if required {
			klog.CtxErrorf(ctx, "未提供认证token")
			return utils.NewBizError(40100, "请先登录")
		}
		// 不需要认证，直接放行
		return nil
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

	return nil
}

// AuthMiddleware 权限验证中间件（保留用于向后兼容）
func AuthMiddleware(required bool, requireAdmin bool) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			if err := verifyAuth(ctx, required, requireAdmin); err != nil {
				return err
			}
			return next(ctx, req, resp)
		}
	}
}
