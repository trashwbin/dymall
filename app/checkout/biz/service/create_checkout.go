package service

import (
	"context"
	checkout "github.com/trashwbin/dymall/rpc_gen/kitex_gen/checkout"
)

type CreateCheckoutService struct {
	ctx context.Context
} // NewCreateCheckoutService new CreateCheckoutService
func NewCreateCheckoutService(ctx context.Context) *CreateCheckoutService {
	return &CreateCheckoutService{ctx: ctx}
}

// Run create note info
func (s *CreateCheckoutService) Run(req *checkout.CreateCheckoutReq) (resp *checkout.CreateCheckoutResp, err error) {
	// Finish your business logic.

	return
}
