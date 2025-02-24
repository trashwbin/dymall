package service

import (
	"context"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

type CreatePaymentService struct {
	ctx context.Context
} // NewCreatePaymentService new CreatePaymentService
func NewCreatePaymentService(ctx context.Context) *CreatePaymentService {
	return &CreatePaymentService{ctx: ctx}
}

// Run create note info
func (s *CreatePaymentService) Run(req *payment.CreatePaymentReq) (resp *payment.CreatePaymentResp, err error) {
	// Finish your business logic.

	return
}
