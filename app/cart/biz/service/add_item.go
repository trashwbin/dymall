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
	ctx context.Context
}

// NewAddItemService new AddItemService
func NewAddItemService(ctx context.Context) *AddItemService {
	return &AddItemService{ctx: ctx}
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

	// 2. 检查购物车是否存在，不存在则创建
	var cartDO mysql.CartDO
	result := mysql.DB.Where("user_id = ? AND status = ?", req.UserId, model.CartStatusNormal).First(&cartDO)
	if result.Error != nil {
		// 创建新购物车
		cartDO = mysql.CartDO{
			UserID: int64(req.UserId),
			Status: int(model.CartStatusNormal),
		}
		if err := mysql.DB.Create(&cartDO).Error; err != nil {
			return nil, utils.NewBizError(50001, "创建购物车失败")
		}
	}

	// 3. 检查商品是否已在购物车中
	var cartItemDO mysql.CartItemDO
	result = mysql.DB.Where("cart_id = ? AND product_id = ?", cartDO.ID, req.Item.ProductId).First(&cartItemDO)
	if result.Error == nil {
		// 更新数量
		cartItemDO.Quantity += req.Item.Quantity
		if err := mysql.DB.Save(&cartItemDO).Error; err != nil {
			return nil, utils.NewBizError(50002, "更新购物车商品失败")
		}
	} else {
		// 添加新商品
		cartItemDO = mysql.CartItemDO{
			CartID:    cartDO.ID,
			UserID:    int64(req.UserId),
			ProductID: int64(req.Item.ProductId),
			Quantity:  req.Item.Quantity,
		}
		if err := mysql.DB.Create(&cartItemDO).Error; err != nil {
			return nil, utils.NewBizError(50003, "添加购物车商品失败")
		}
	}

	// 4. 更新缓存
	cartItem := cartItemDO.ToModel()
	cartRepo := redis.NewCartRepo()
	if err := cartRepo.SetCartItem(s.ctx, cartItem); err != nil {
		// 缓存错误不影响主流程
		klog.CtxWarnf(s.ctx, "更新购物车商品缓存失败 - userId: %d, cartId: %d, productId: %d, err: %v",
			req.UserId, cartDO.ID, req.Item.ProductId, err)
	}

	return &cart.AddItemResp{}, nil
}
