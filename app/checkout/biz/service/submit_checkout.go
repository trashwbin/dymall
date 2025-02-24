package service

import (
	"context"
	checkout "github.com/trashwbin/dymall/rpc_gen/kitex_gen/checkout"
)

type SubmitCheckoutService struct {
	ctx context.Context
} // NewSubmitCheckoutService new SubmitCheckoutService
func NewSubmitCheckoutService(ctx context.Context) *SubmitCheckoutService {
	return &SubmitCheckoutService{ctx: ctx}
}

// Run create note info
func (s *SubmitCheckoutService) Run(req *checkout.SubmitCheckoutReq) (resp *checkout.SubmitCheckoutResp, err error) {
	// Finish your business logic.

	return
}
