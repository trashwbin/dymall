package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/cart/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/cart/biz/dal/redis"
	"github.com/trashwbin/dymall/app/cart/biz/model"
	"github.com/trashwbin/dymall/app/cart/infra/rpc"
	"github.com/trashwbin/dymall/app/cart/utils"
	cart "github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

type AddItemService struct {
	ctx       context.Context
	cartRepo  *mysql.CartRepo
	cacheRepo *redis.CartRepo
}

// NewAddItemService new AddItemService
func NewAddItemService(ctx context.Context) *AddItemService {
	return &AddItemService{
		ctx:       ctx,
		cartRepo:  mysql.NewCartRepo(),
		cacheRepo: redis.NewCartRepo(),
	}
}

// Run 添加商品到购物车
func (s *AddItemService) Run(req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	// 1. 检查商品是否存在
	productResp, err := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{
		Id: req.Item.ProductId,
	})
	if err != nil {
		return nil, err
	}
	if productResp == nil || productResp.Product == nil {
		return nil, utils.NewBizError(40004, "商品不存在")
	}

	// 使用事务处理购物车操作
	err = s.cartRepo.Transaction(func(txRepo *mysql.CartRepo) error {
		// 2. 检查购物车是否存在，不存在则创建
		cart, err := txRepo.GetCartByUserID(int64(req.UserId), model.CartStatusNormal)
		if err != nil {
			klog.CtxInfof(s.ctx, "用户购物车不存在，创建新购物车 - userId: %d", req.UserId)
			// 创建新购物车
			cart = &model.Cart{
				UserID: int64(req.UserId),
				Status: model.CartStatusNormal,
			}
			cart, err = txRepo.CreateCart(cart)
			if err != nil {
				return utils.NewBizError(50001, "创建购物车失败")
			}
			// 更新缓存
			if err := s.cacheRepo.SetCart(s.ctx, cart); err != nil {
				klog.CtxWarnf(s.ctx, "更新购物车缓存失败 - userId: %d, cartId: %d, err: %v",
					req.UserId, cart.ID, err)
			}
			klog.CtxInfof(s.ctx, "创建购物车成功 - userId: %d, cartId: %d", req.UserId, cart.ID)
		}

		// 3. 检查商品是否已在购物车中
		cartItem, err := txRepo.GetCartItem(cart.ID, int64(req.Item.ProductId))
		if err == nil {
			// 更新数量
			oldQuantity := cartItem.Quantity
			cartItem.Quantity += req.Item.Quantity
			if err := txRepo.UpdateCartItem(cartItem); err != nil {
				return utils.NewBizError(50002, "更新购物车商品失败")
			}
			// 更新缓存
			if err := s.cacheRepo.SetCartItem(s.ctx, cartItem); err != nil {
				klog.CtxWarnf(s.ctx, "更新购物车商品缓存失败 - cartId: %d, productId: %d, err: %v",
					cart.ID, cartItem.ProductID, err)
			}
			klog.CtxInfof(s.ctx, "更新购物车商品数量 - userId: %d, cartId: %d, productId: %d, oldQuantity: %d, newQuantity: %d",
				req.UserId, cart.ID, req.Item.ProductId, oldQuantity, cartItem.Quantity)
		} else {
			// 添加新商品
			cartItem = &model.CartItem{
				CartID:    cart.ID,
				UserID:    int64(req.UserId),
				ProductID: int64(req.Item.ProductId),
				Quantity:  req.Item.Quantity,
			}
			var err error
			cartItem, err = txRepo.CreateCartItem(cartItem)
			if err != nil {
				return utils.NewBizError(50003, "添加购物车商品失败")
			}
			// 更新缓存
			if err := s.cacheRepo.SetCartItem(s.ctx, cartItem); err != nil {
				klog.CtxWarnf(s.ctx, "更新购物车商品缓存失败 - cartId: %d, productId: %d, err: %v",
					cart.ID, cartItem.ProductID, err)
			}
			klog.CtxInfof(s.ctx, "添加新商品到购物车 - userId: %d, cartId: %d, productId: %d, quantity: %d",
				req.UserId, cart.ID, req.Item.ProductId, cartItem.Quantity)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &cart.AddItemResp{}, nil
}
