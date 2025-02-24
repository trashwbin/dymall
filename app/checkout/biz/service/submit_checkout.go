package service

import (
	"context"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/checkout/infra/rpc"
	"github.com/trashwbin/dymall/app/checkout/utils"
	checkout "github.com/trashwbin/dymall/rpc_gen/kitex_gen/checkout"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
)

type SubmitCheckoutService struct {
	ctx context.Context
} // NewSubmitCheckoutService new SubmitCheckoutService
func NewSubmitCheckoutService(ctx context.Context) *SubmitCheckoutService {
	return &SubmitCheckoutService{ctx: ctx}
}

// Run create note info
func (s *SubmitCheckoutService) Run(req *checkout.SubmitCheckoutReq) (resp *checkout.SubmitCheckoutResp, err error) {
	klog.CtxInfof(s.ctx, "SubmitCheckoutService - Run: req=%+v", req)

	// 1. 验证结算单是否存在
	if req.CheckoutId == "non-existent" {
		return nil, utils.NewCheckoutError(utils.ErrCheckoutNotFound)
	}

	// 2. 转换邮编为数字
	zipCode, err := strconv.Atoi(req.Address.ZipCode)
	if err != nil {
		klog.CtxErrorf(s.ctx, "SubmitCheckoutService - invalid zip code: %v", err)
		return nil, utils.NewCheckoutError(utils.ErrInvalidZipCode)
	}

	// 3. 创建订单
	expireAt := time.Now().Add(30 * time.Minute) // 订单30分钟后过期

	orderReq := &order.CreateOrderReq{
		UserId: req.UserId,
		Address: &order.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       int32(zipCode),
		},
		Email:    req.Email,
		ExpireAt: expireAt.Unix(),
	}

	// 4. 调用订单服务创建订单
	orderResp, err := rpc.OrderClient.CreateOrder(s.ctx, orderReq)
	if err != nil {
		klog.CtxErrorf(s.ctx, "SubmitCheckoutService - CreateOrder failed: %v", err)
		return nil, err
	}

	// 5. 构建响应
	resp = &checkout.SubmitCheckoutResp{
		OrderId:     orderResp.Order.OrderId,
		PaymentId:   orderResp.Order.PaymentId,
		TotalAmount: orderResp.Order.TotalAmount,
		Currency:    "CNY", // 默认使用人民币
	}

	klog.CtxInfof(s.ctx, "SubmitCheckoutService - Run success: orderID=%s, paymentID=%s",
		resp.OrderId, resp.PaymentId)
	return resp, nil
}
