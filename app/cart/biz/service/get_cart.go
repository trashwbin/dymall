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

type GetCartService struct {
	ctx       context.Context
	cartRepo  *mysql.CartRepo
	cacheRepo *redis.CartRepo
}

// NewGetCartService new GetCartService
func NewGetCartService(ctx context.Context) *GetCartService {
	return &GetCartService{
		ctx:       ctx,
		cartRepo:  mysql.NewCartRepo(),
		cacheRepo: redis.NewCartRepo(),
	}
}

// Run 获取购物车信息
func (s *GetCartService) Run(req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	// 1. 尝试从缓存获取购物车
	userCart, err := s.cacheRepo.GetCart(s.ctx, int64(req.UserId))
	if err != nil || userCart == nil || userCart.Status != model.CartStatusNormal {
		// 缓存未命中或购物车状态不正常，从数据库获取
		userCart, err = s.cartRepo.GetCartByUserID(int64(req.UserId), model.CartStatusNormal)
		if err != nil {
			return nil, utils.NewBizError(40001, "购物车不存在")
		}
		// 更新缓存
		if err := s.cacheRepo.SetCart(s.ctx, userCart); err != nil {
			klog.CtxWarnf(s.ctx, "更新购物车缓存失败 - userId: %d, err: %v", req.UserId, err)
		}
	}

	// 2. 尝试从缓存获取购物车商品列表
	items, err := s.cacheRepo.GetCartItems(s.ctx, userCart.ID)
	if err != nil || len(items) == 0 {
		// 缓存未命中，从数据库获取
		items, err = s.cartRepo.GetCartItems(userCart.ID)
		if err != nil {
			klog.CtxWarnf(s.ctx, "获取购物车商品列表失败 - userId: %d, cartId: %d, err: %v",
				req.UserId, userCart.ID, err)
			return nil, utils.NewBizError(50001, "获取购物车商品列表失败")
		}
		// 批量更新缓存
		for _, item := range items {
			if err := s.cacheRepo.SetCartItem(s.ctx, item); err != nil {
				klog.CtxWarnf(s.ctx, "更新购物车商品缓存失败 - cartId: %d, productId: %d, err: %v",
					userCart.ID, item.ProductID, err)
			}
		}
	}

	// 3. 调用商品服务获取商品详情
	cartItems := make([]*cart.CartItem, 0, len(items))
	for _, item := range items {
		productResp, err := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{
			Id: uint32(item.ProductID),
		})
		if err != nil {
			klog.CtxWarnf(s.ctx, "获取商品信息失败 - productId: %d, err: %v", item.ProductID, err)
			continue
		}
		if productResp == nil || productResp.Product == nil {
			klog.CtxWarnf(s.ctx, "商品不存在 - productId: %d", item.ProductID)
			continue
		}

		cartItems = append(cartItems, &cart.CartItem{
			ProductId: uint32(item.ProductID),
			Quantity:  item.Quantity,
		})
	}

	resp = &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: uint32(req.UserId),
			Items:  cartItems,
		},
	}

	return resp, nil
}
