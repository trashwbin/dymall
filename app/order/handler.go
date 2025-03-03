package main

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/trashwbin/dymall/app/order/biz/service"
	"github.com/trashwbin/dymall/app/order/middleware"
	order "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// CreateOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *order.CreateOrderReq) (resp *order.CreateOrderResp, err error) {
	// 只允许 checkout 服务调用
	if err := endpoint.Chain(
		middleware.ServiceAuthMiddleware(middleware.CheckoutService),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*order.CreateOrderReq)
		response, err := service.NewCreateOrderService(ctx).Run(r)
		if err == nil && response != nil {
			resp.(*order.CreateOrderResp).Order = response.Order
		}
		return err
	})(ctx, req, &resp); err != nil {
		return nil, err
	}
	return resp, err
}

// UpdateOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrder(ctx context.Context, req *order.UpdateOrderReq) (resp *order.UpdateOrderResp, err error) {
	// 需要用户认证
	if err := endpoint.Chain(
		middleware.UserAuthMiddleware(),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*order.UpdateOrderReq)
		response, err := service.NewUpdateOrderService(ctx).Run(r)
		if err == nil && response != nil {
			resp.(*order.UpdateOrderResp).Order = response.Order
		}
		return err
	})(ctx, req, &resp); err != nil {
		return nil, err
	}
	return resp, err
}

// CancelOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, req *order.CancelOrderReq) (resp *order.CancelOrderResp, err error) {
	// 需要用户认证
	if err := endpoint.Chain(
		middleware.UserAuthMiddleware(),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*order.CancelOrderReq)
		response, err := service.NewCancelOrderService(ctx).Run(r)
		if err == nil && response != nil {
			// CancelOrderResp 是空结构体，无需赋值
		}
		return err
	})(ctx, req, &resp); err != nil {
		return nil, err
	}
	return resp, err
}

// GetOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetOrder(ctx context.Context, req *order.GetOrderReq) (resp *order.GetOrderResp, err error) {
	// 需要用户认证
	if err := endpoint.Chain(
		middleware.UserAuthMiddleware(),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*order.GetOrderReq)
		response, err := service.NewGetOrderService(ctx).Run(r)
		if err == nil && response != nil {
			resp.(*order.GetOrderResp).Order = response.Order
		}
		return err
	})(ctx, req, &resp); err != nil {
		return nil, err
	}
	return resp, err
}

// ListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	// 需要用户认证
	if err := endpoint.Chain(
		middleware.UserAuthMiddleware(),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*order.ListOrderReq)
		response, err := service.NewListOrderService(ctx).Run(r)
		if err == nil && response != nil {
			resp.(*order.ListOrderResp).Orders = response.Orders
		}
		return err
	})(ctx, req, &resp); err != nil {
		return nil, err
	}
	return resp, err
}

// MarkOrderPaid implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	// 只允许 payment 服务调用
	if err := endpoint.Chain(
		middleware.ServiceAuthMiddleware(middleware.PaymentService),
	)(func(ctx context.Context, req, resp interface{}) error {
		var r = req.(*order.MarkOrderPaidReq)
		response, err := service.NewMarkOrderPaidService(ctx).Run(r)
		if err == nil && response != nil {
			// MarkOrderPaidResp 是空结构体，无需赋值
		}
		return err
	})(ctx, req, &resp); err != nil {
		return nil, err
	}
	return resp, err
}
