package main

import (
	"context"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
	"github.com/trashwbin/dymall/app/payment/biz/service"
)

// PaymentServiceImpl implements the last service interface defined in the IDL.
type PaymentServiceImpl struct{}

// CreatePayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, req *payment.CreatePaymentReq) (resp *payment.CreatePaymentResp, err error) {
	resp, err = service.NewCreatePaymentService(ctx).Run(req)

	return resp, err
}

// ProcessPayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) ProcessPayment(ctx context.Context, req *payment.ProcessPaymentReq) (resp *payment.ProcessPaymentResp, err error) {
	resp, err = service.NewProcessPaymentService(ctx).Run(req)

	return resp, err
}

// CancelPayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) CancelPayment(ctx context.Context, req *payment.CancelPaymentReq) (resp *payment.CancelPaymentResp, err error) {
	resp, err = service.NewCancelPaymentService(ctx).Run(req)

	return resp, err
}

// GetPayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) GetPayment(ctx context.Context, req *payment.GetPaymentReq) (resp *payment.GetPaymentResp, err error) {
	resp, err = service.NewGetPaymentService(ctx).Run(req)

	return resp, err
}
