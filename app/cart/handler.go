package main

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/cart/biz/service"
	"github.com/trashwbin/dymall/app/cart/infra/rpc"
	"github.com/trashwbin/dymall/app/cart/utils"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
	cart "github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
)

// CartServiceImpl implements the last service interface defined in the IDL.
type CartServiceImpl struct{}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// verifyToken 验证令牌和权限
func verifyToken(ctx context.Context, token string, requiredRole string) (*auth.VerifyResp, error) {
	resp, err := rpc.AuthClient.VerifyTokenByRPC(ctx, &auth.VerifyTokenReq{
		Token: token,
	})
	if err != nil {
		klog.CtxErrorf(ctx, "验证令牌失败: %v", err)
		return nil, utils.NewBizError(40001, "验证令牌失败")
	}
	if resp.Code != auth.ErrorCode_Success {
		klog.CtxErrorf(ctx, "令牌无效: code=%v, message=%v", resp.Code, resp.Message)
		return nil, utils.NewBizError(40001, resp.Message)
	}

	// 验证角色
	if resp.Role != requiredRole && resp.Role != RoleAdmin {
		klog.CtxErrorf(ctx, "权限不足: role=%v, requiredRole=%v", resp.Role, requiredRole)
		return nil, utils.NewBizError(40002, "权限不足")
	}

	return resp, nil
}

// AddItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	// 验证令牌和权限
	// ctx = metainfo.WithValue(ctx, "token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoidXNlciIsImlzcyI6ImF1dGgtc2VydmljZSIsImV4cCI6MTc0MTExMzE5MywibmJmIjoxNzQxMDI2NzkzLCJpYXQiOjE3NDEwMjY3OTN9.FFX_iZcs-Fcmc5RUtDmD9coUITHuQQGtjvGlj2vgRBA")
	token := utils.GetTokenFromContext(ctx)
	verifyResp, err := verifyToken(ctx, token, RoleUser)
	if err != nil {
		return nil, err
	}

	// 验证用户身份
	if verifyResp.UserId != int64(req.UserId) && verifyResp.Role != RoleAdmin {
		return nil, utils.NewBizError(40003, "无权访问该用户的购物车")
	}

	resp, err = service.NewAddItemService(ctx).Run(req)
	return resp, err
}

// GetCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	// 验证令牌和权限
	token := utils.GetTokenFromContext(ctx)
	verifyResp, err := verifyToken(ctx, token, RoleUser)
	if err != nil {
		return nil, err
	}

	// 验证用户身份
	if verifyResp.UserId != int64(req.UserId) && verifyResp.Role != RoleAdmin {
		return nil, utils.NewBizError(40003, "无权访问该用户的购物车")
	}

	resp, err = service.NewGetCartService(ctx).Run(req)
	return resp, err
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	// 验证令牌和权限
	token := utils.GetTokenFromContext(ctx)
	verifyResp, err := verifyToken(ctx, token, RoleUser)
	if err != nil {
		return nil, err
	}

	// 验证用户身份
	if verifyResp.UserId != int64(req.UserId) && verifyResp.Role != RoleAdmin {
		return nil, utils.NewBizError(40003, "无权访问该用户的购物车")
	}

	resp, err = service.NewEmptyCartService(ctx).Run(req)
	return resp, err
}
