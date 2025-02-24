package service

import (
	"context"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

type ProcessPaymentService struct {
	ctx context.Context
} // NewProcessPaymentService new ProcessPaymentService
func NewProcessPaymentService(ctx context.Context) *ProcessPaymentService {
	return &ProcessPaymentService{ctx: ctx}
}

// Run create note info
func (s *ProcessPaymentService) Run(req *payment.ProcessPaymentReq) (resp *payment.ProcessPaymentResp, err error) {
	// Finish your business logic.

	return
}
