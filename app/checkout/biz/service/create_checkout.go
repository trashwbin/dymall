package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/uuid"
	"github.com/trashwbin/dymall/app/checkout/infra/rpc"
	"github.com/trashwbin/dymall/app/checkout/utils"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
	checkout "github.com/trashwbin/dymall/rpc_gen/kitex_gen/checkout"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

type CreateCheckoutService struct {
	ctx context.Context
} // NewCreateCheckoutService new CreateCheckoutService
func NewCreateCheckoutService(ctx context.Context) *CreateCheckoutService {
	return &CreateCheckoutService{ctx: ctx}
}

// Run create note info
func (s *CreateCheckoutService) Run(req *checkout.CreateCheckoutReq) (resp *checkout.CreateCheckoutResp, err error) {
	klog.CtxInfof(s.ctx, "CreateCheckoutService - Run: req=%+v", req)

	// 1. 获取用户购物车
	cartResp, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "CreateCheckoutService - GetCart failed: %v", err)
		return nil, err
	}

	// 2. 验证购物车商品
	if len(cartResp.Cart.Items) == 0 {
		klog.CtxErrorf(s.ctx, "CreateCheckoutService - cart is empty")
		return nil, utils.NewCheckoutError(utils.ErrCartEmpty)
	}

	// 3. 获取商品详情并计算总价
	var totalAmount float32
	currency := "CNY" // 默认使用人民币

	for _, item := range cartResp.Cart.Items {
		// 3.1 验证商品数量
		if item.Quantity <= 0 {
			klog.CtxErrorf(s.ctx, "CreateCheckoutService - invalid quantity: %d", item.Quantity)
			return nil, utils.NewCheckoutError(utils.ErrInvalidQuantity)
		}

		// 3.2 获取商品信息
		productResp, err := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{
			ProductId: uint32(item.ProductId),
		})
		if err != nil {
			klog.CtxErrorf(s.ctx, "CreateCheckoutService - GetProduct failed: %v", err)
			return nil, err
		}

		// 3.3 验证商品库存
		if productResp.Product.Stock < item.Quantity {
			klog.CtxErrorf(s.ctx, "CreateCheckoutService - insufficient stock: want=%d, got=%d",
				item.Quantity, productResp.Product.Stock)
			return nil, utils.NewCheckoutError(utils.ErrInsufficientStock)
		}

		// 3.4 计算商品总价
		totalAmount += productResp.Product.Price * float32(item.Quantity)
	}

	// 4. 创建结算单
	checkoutID := uuid.New().String()

	// 5. 构建响应
	resp = &checkout.CreateCheckoutResp{
		CheckoutId:  checkoutID,
		Items:       cartResp.Cart.Items,
		TotalAmount: totalAmount,
		Currency:    currency,
	}

	klog.CtxInfof(s.ctx, "CreateCheckoutService - Run success: checkoutID=%s, totalAmount=%f",
		checkoutID, totalAmount)
	return resp, nil
}
