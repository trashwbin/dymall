package main

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/trashwbin/dymall/app/payment/biz/service"
	"github.com/trashwbin/dymall/app/payment/middleware"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

// PaymentServiceImpl implements the last service interface defined in the IDL.
type PaymentServiceImpl struct{}

// CreatePayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, req *payment.CreatePaymentReq) (resp *payment.CreatePaymentResp, err error) {
	resp = &payment.CreatePaymentResp{} // 初始化响应对象
	// 只允许 order 服务调用
	if err := endpoint.Chain(
		middleware.ServiceAuthMiddleware(middleware.OrderService),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*payment.CreatePaymentReq)
		response, err := service.NewCreatePaymentService(ctx).Run(r)
		if err == nil && response != nil {
			resp.(*payment.CreatePaymentResp).Payment = response.Payment
		}
		return err
	})(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, err
}

// ProcessPayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) ProcessPayment(ctx context.Context, req *payment.ProcessPaymentReq) (resp *payment.ProcessPaymentResp, err error) {
	resp = &payment.ProcessPaymentResp{} // 初始化响应对象
	// 需要用户认证
	if err := endpoint.Chain(
		middleware.UserAuthMiddleware(),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*payment.ProcessPaymentReq)
		response, err := service.NewProcessPaymentService(ctx).Run(r)
		if err == nil && response != nil {
			resp.(*payment.ProcessPaymentResp).Payment = response.Payment
		}
		return err
	})(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, err
}

// CancelPayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) CancelPayment(ctx context.Context, req *payment.CancelPaymentReq) (resp *payment.CancelPaymentResp, err error) {
	resp = &payment.CancelPaymentResp{} // 初始化响应对象
	// 允许用户或定时服务调用
	if err := endpoint.Chain(
		middleware.ServiceAuthMiddleware(middleware.SchedulerService),
		middleware.UserAuthMiddleware(),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*payment.CancelPaymentReq)
		response, err := service.NewCancelPaymentService(ctx).Run(r)
		if err == nil && response != nil {
			// CancelPaymentResp 是空结构体，无需赋值
		}
		return err
	})(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, err
}

// GetPayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) GetPayment(ctx context.Context, req *payment.GetPaymentReq) (resp *payment.GetPaymentResp, err error) {
	resp = &payment.GetPaymentResp{} // 初始化响应对象
	// 需要用户认证
	if err := endpoint.Chain(
		middleware.UserAuthMiddleware(),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*payment.GetPaymentReq)
		response, err := service.NewGetPaymentService(ctx).Run(r)
		if err == nil && response != nil {
			resp.(*payment.GetPaymentResp).Payment = response.Payment
		}
		return err
	})(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, err
}
