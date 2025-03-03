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

	// 3. 收集所有商品ID
	productIDs := make([]uint32, 0, len(cartResp.Cart.Items))
	quantityMap := make(map[uint32]int32)
	for _, item := range cartResp.Cart.Items {
		// 3.1 验证商品数量
		if item.Quantity <= 0 {
			klog.CtxErrorf(s.ctx, "CreateCheckoutService - invalid quantity: %d", item.Quantity)
			return nil, utils.NewCheckoutError(utils.ErrInvalidQuantity)
		}
		productIDs = append(productIDs, uint32(item.ProductId))
		quantityMap[uint32(item.ProductId)] = item.Quantity
	}

	// 4. 批量获取商品信息
	productsResp, err := rpc.ProductClient.BatchGetProducts(s.ctx, &product.BatchGetProductsReq{
		Ids: productIDs,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "CreateCheckoutService - BatchGetProducts failed: %v", err)
		return nil, err
	}

	// 5. 验证是否所有商品都存在
	if len(productsResp.MissingIds) > 0 {
		klog.CtxErrorf(s.ctx, "CreateCheckoutService - some products not found: %v", productsResp.MissingIds)
		return nil, utils.NewCheckoutError(utils.ErrProductNotFound)
	}

	// 6. 计算总价并验证库存
	var totalAmount float32
	currency := "CNY" // 默认使用人民币

	for _, p := range productsResp.Products {
		quantity := quantityMap[p.Id]
		// 6.1 计算商品总价
		totalAmount += p.Price * float32(quantity)
	}

	// 7. 创建结算单
	checkoutID := uuid.New().String()

	// 8. 构建响应
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
