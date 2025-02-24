package main

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/trashwbin/dymall/app/checkout/biz/service"
	"github.com/trashwbin/dymall/app/checkout/middleware"
	checkout "github.com/trashwbin/dymall/rpc_gen/kitex_gen/checkout"
)

// CheckoutServiceImpl implements the last service interface defined in the IDL.
type CheckoutServiceImpl struct{}

// CreateCheckout implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) CreateCheckout(ctx context.Context, req *checkout.CreateCheckoutReq) (resp *checkout.CreateCheckoutResp, err error) {
	// 需要用户认证
	if err := endpoint.Chain(
		middleware.UserAuthMiddleware(),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*checkout.CreateCheckoutReq)
		response, err := service.NewCreateCheckoutService(ctx).Run(r)
		if err == nil && response != nil {
			resp.(*checkout.CreateCheckoutResp).CheckoutId = response.CheckoutId
			resp.(*checkout.CreateCheckoutResp).Items = response.Items
			resp.(*checkout.CreateCheckoutResp).TotalAmount = response.TotalAmount
			resp.(*checkout.CreateCheckoutResp).Currency = response.Currency
		}
		return err
	})(ctx, req, &resp); err != nil {
		return nil, err
	}
	return resp, err
}

// SubmitCheckout implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) SubmitCheckout(ctx context.Context, req *checkout.SubmitCheckoutReq) (resp *checkout.SubmitCheckoutResp, err error) {
	// 需要用户认证
	if err := endpoint.Chain(
		middleware.UserAuthMiddleware(),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*checkout.SubmitCheckoutReq)
		response, err := service.NewSubmitCheckoutService(ctx).Run(r)
		if err == nil && response != nil {
			resp.(*checkout.SubmitCheckoutResp).OrderId = response.OrderId
			resp.(*checkout.SubmitCheckoutResp).PaymentId = response.PaymentId
			resp.(*checkout.SubmitCheckoutResp).TotalAmount = response.TotalAmount
			resp.(*checkout.SubmitCheckoutResp).Currency = response.Currency
		}
		return err
	})(ctx, req, &resp); err != nil {
		return nil, err
	}
	return resp, err
}
