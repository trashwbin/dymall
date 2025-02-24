package main

import (
	"context"
	checkout "github.com/trashwbin/dymall/rpc_gen/kitex_gen/checkout"
	"github.com/trashwbin/dymall/app/checkout/biz/service"
)

// CheckoutServiceImpl implements the last service interface defined in the IDL.
type CheckoutServiceImpl struct{}

// CreateCheckout implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) CreateCheckout(ctx context.Context, req *checkout.CreateCheckoutReq) (resp *checkout.CreateCheckoutResp, err error) {
	resp, err = service.NewCreateCheckoutService(ctx).Run(req)

	return resp, err
}

// SubmitCheckout implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) SubmitCheckout(ctx context.Context, req *checkout.SubmitCheckoutReq) (resp *checkout.SubmitCheckoutResp, err error) {
	resp, err = service.NewSubmitCheckoutService(ctx).Run(req)

	return resp, err
}
