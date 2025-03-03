package service

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/product/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/product/biz/dal/redis"
	"github.com/trashwbin/dymall/app/product/biz/model"
	"github.com/trashwbin/dymall/app/product/biz/service/converter"
	"github.com/trashwbin/dymall/app/product/utils"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

type UpdateProductService struct {
	ctx         context.Context
	productRepo *mysql.ProductRepo
	cacheRepo   *redis.ProductRepo
}

// NewUpdateProductService new UpdateProductService
func NewUpdateProductService(ctx context.Context) *UpdateProductService {
	return &UpdateProductService{
		ctx:         ctx,
		productRepo: mysql.NewProductRepo(),
		cacheRepo:   redis.NewProductRepo(),
	}
}

// Run 更新商品信息
func (s *UpdateProductService) Run(req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	resp = new(product.UpdateProductResp)

	// 1. 参数验证
	if req.Id == 0 {
		return nil, utils.NewBizError(40001, "商品ID不能为空")
	}
	if req.Name == "" {
		return nil, utils.NewBizError(40002, "商品名称不能为空")
	}
	if req.Price < 0 {
		return nil, utils.NewBizError(40003, "商品价格不能为负数")
	}

	// 2. 检查商品是否存在
	oldProduct, err := s.productRepo.GetProduct(req.Id)
	if err != nil {
		klog.CtxErrorf(s.ctx, "商品不存在 - productId: %d, err: %v", req.Id, err)
		return nil, utils.NewBizError(40004, "商品不存在")
	}

	// 3. 构建更新后的商品模型
	productModel := &model.Product{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  req.Categories,
		CreatedAt:   oldProduct.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	// 4. 验证商品信息
	if !productModel.IsValid() {
		return nil, utils.NewBizError(40005, "商品信息不合法")
	}

	// 5. 使用事务更新商品
	err = s.productRepo.Transaction(func(txRepo *mysql.ProductRepo) error {
		// 5.1 更新商品基本信息
		if err := txRepo.UpdateProduct(productModel); err != nil {
			klog.CtxErrorf(s.ctx, "更新商品失败 - productId: %d, err: %v", req.Id, err)
			return utils.NewBizError(50001, "更新商品失败")
		}

		// 5.2 更新商品分类关系
		// 5.2.1 获取原有分类
		oldCategories, err := txRepo.GetCategories(req.Id)
		if err != nil {
			klog.CtxWarnf(s.ctx, "获取商品原有分类失败 - productId: %d, err: %v", req.Id, err)
		}

		// 5.2.2 删除所有原有分类关系
		for _, category := range oldCategories {
			if err := txRepo.RemoveCategory(req.Id, category.ID); err != nil {
				klog.CtxWarnf(s.ctx, "删除商品分类关系失败 - productId: %d, categoryId: %d, err: %v",
					req.Id, category.ID, err)
			}
		}

		// 5.2.3 添加新的分类关系
		for _, categoryName := range req.Categories {
			// 获取或创建分类
			category, err := txRepo.GetOrCreateCategory(categoryName)
			if err != nil {
				klog.CtxWarnf(s.ctx, "处理商品分类失败 - productId: %d, category: %s, err: %v",
					productModel.ID, categoryName, err)
				continue
			}

			// 添加商品分类关系
			if err := txRepo.AddProductCategory(productModel.ID, category.ID); err != nil {
				klog.CtxWarnf(s.ctx, "添加商品分类关系失败 - productId: %d, categoryId: %d, err: %v",
					productModel.ID, category.ID, err)
			} else {
				klog.CtxInfof(s.ctx, "添加商品分类成功 - productId: %d, category: %s",
					productModel.ID, categoryName)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 6. 更新缓存
	if err := s.cacheRepo.SetProduct(s.ctx, productModel); err != nil {
		klog.CtxWarnf(s.ctx, "更新商品缓存失败 - productId: %d, err: %v", productModel.ID, err)
	}

	// 7. 删除商品列表缓存（因为更新了商品，所以需要清理列表缓存）
	if err := s.cacheRepo.DeleteProductList(s.ctx, 1, 10, ""); err != nil {
		klog.CtxWarnf(s.ctx, "删除商品列表缓存失败 - err: %v", err)
	}

	resp.Product = converter.ToProductProto(productModel)
	klog.CtxInfof(s.ctx, "更新商品成功 - productId: %d, name: %s", productModel.ID, req.Name)
	return resp, nil
}
