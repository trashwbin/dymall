package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/product/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/product/biz/dal/redis"
	"github.com/trashwbin/dymall/app/product/utils"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

type DeleteProductService struct {
	ctx         context.Context
	productRepo *mysql.ProductRepo
	cacheRepo   *redis.ProductRepo
}

// NewDeleteProductService new DeleteProductService
func NewDeleteProductService(ctx context.Context) *DeleteProductService {
	return &DeleteProductService{
		ctx:         ctx,
		productRepo: mysql.NewProductRepo(),
		cacheRepo:   redis.NewProductRepo(),
	}
}

// Run 删除商品
func (s *DeleteProductService) Run(req *product.DeleteProductReq) (resp *product.DeleteProductResp, err error) {
	resp = new(product.DeleteProductResp)
	resp.Success = false // 默认设置为失败

	// 1. 检查商品是否存在
	productModel, err := s.productRepo.GetProduct(req.Id)
	if err != nil {
		klog.CtxErrorf(s.ctx, "商品不存在 - productId: %d, err: %v", req.Id, err)
		return resp, utils.NewBizError(40004, "商品不存在")
	}

	// 2. 使用事务删除商品
	err = s.productRepo.Transaction(func(txRepo *mysql.ProductRepo) error {
		// 2.1 删除商品（软删除）
		if err := txRepo.DeleteProduct(req.Id); err != nil {
			klog.CtxErrorf(s.ctx, "删除商品失败 - productId: %d, err: %v", req.Id, err)
			return utils.NewBizError(50001, "删除商品失败")
		}

		// 2.2 删除商品分类关系
		categories, err := txRepo.GetCategories(req.Id)
		if err != nil {
			klog.CtxWarnf(s.ctx, "获取商品分类失败 - productId: %d, err: %v", req.Id, err)
		} else {
			for _, category := range categories {
				if err := txRepo.RemoveCategory(req.Id, category.ID); err != nil {
					klog.CtxWarnf(s.ctx, "删除商品分类关系失败 - productId: %d, categoryId: %d, err: %v",
						req.Id, category.ID, err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return resp, err
	}

	// 3. 删除商品缓存
	if err := s.cacheRepo.DeleteProduct(s.ctx, req.Id); err != nil {
		klog.CtxWarnf(s.ctx, "删除商品缓存失败 - productId: %d, err: %v", req.Id, err)
	}

	// 4. 删除商品列表缓存
	if err := s.cacheRepo.DeleteProductList(s.ctx, 1, 10, ""); err != nil {
		klog.CtxWarnf(s.ctx, "删除商品列表缓存失败 - err: %v", err)
	}

	resp.Success = true
	klog.CtxInfof(s.ctx, "删除商品成功 - productId: %d, name: %s", productModel.ID, productModel.Name)
	return resp, nil
}
