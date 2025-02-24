package service

import (
	"context"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

type GetPaymentService struct {
	ctx context.Context
} // NewGetPaymentService new GetPaymentService
func NewGetPaymentService(ctx context.Context) *GetPaymentService {
	return &GetPaymentService{ctx: ctx}
}

// Run create note info
func (s *GetPaymentService) Run(req *payment.GetPaymentReq) (resp *payment.GetPaymentResp, err error) {
	// Finish your business logic.

	return
}
