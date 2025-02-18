package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/cart/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/cart/biz/dal/redis"
	"github.com/trashwbin/dymall/app/cart/biz/model"
	"github.com/trashwbin/dymall/app/cart/utils"
	cart "github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
)

type EmptyCartService struct {
	ctx       context.Context
	cartRepo  *mysql.CartRepo
	cacheRepo *redis.CartRepo
}

// NewEmptyCartService new EmptyCartService
func NewEmptyCartService(ctx context.Context) *EmptyCartService {
	return &EmptyCartService{
		ctx:       ctx,
		cartRepo:  mysql.NewCartRepo(),
		cacheRepo: redis.NewCartRepo(),
	}
}

// Run 清空购物车
func (s *EmptyCartService) Run(req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	// 1. 检查购物车是否存在
	userCart, err := s.cartRepo.GetCartByUserID(int64(req.UserId), model.CartStatusNormal)
	if err != nil {
		return nil, utils.NewBizError(40001, "购物车不存在")
	}

	// 2. 获取购物车商品列表（用于后续清理缓存）
	items, err := s.cartRepo.GetCartItems(int64(req.UserId))
	if err != nil {
		klog.CtxWarnf(s.ctx, "获取购物车商品列表失败 - userId: %d, err: %v", req.UserId, err)
	}

	// 3. 使用事务清空购物车
	err = s.cartRepo.Transaction(func(txRepo *mysql.CartRepo) error {
		if err := txRepo.EmptyCart(int64(req.UserId)); err != nil {
			return utils.NewBizError(50004, "清空购物车失败")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 4. 清理缓存
	// 注意：即使缓存清理失败也不影响主流程
	for _, item := range items {
		if err := s.cacheRepo.DeleteCartItem(s.ctx, userCart.ID, item.ProductID); err != nil {
			klog.CtxWarnf(s.ctx, "清理购物车商品缓存失败 - userId: %d, cartId: %d, productId: %d, err: %v",
				req.UserId, userCart.ID, item.ProductID, err)
		}
	}

	klog.CtxInfof(s.ctx, "清空购物车成功 - userId: %d, cartId: %d", req.UserId, userCart.ID)
	return &cart.EmptyCartResp{}, nil
}
