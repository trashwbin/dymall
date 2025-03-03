package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/product/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/product/biz/dal/redis"
	"github.com/trashwbin/dymall/app/product/biz/service/converter"
	"github.com/trashwbin/dymall/app/product/utils"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

type GetProductService struct {
	ctx         context.Context
	productRepo *mysql.ProductRepo
	cacheRepo   *redis.ProductRepo
}

// NewGetProductService new GetProductService
func NewGetProductService(ctx context.Context) *GetProductService {
	return &GetProductService{
		ctx:         ctx,
		productRepo: mysql.NewProductRepo(),
		cacheRepo:   redis.NewProductRepo(),
	}
}

// Run 获取商品信息
func (s *GetProductService) Run(req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	resp = new(product.GetProductResp)

	// 1. 尝试从缓存获取商品信息
	productModel, err := s.cacheRepo.GetProduct(s.ctx, req.Id)
	if err == nil {
		// 缓存命中，获取商品分类信息
		categories, err := s.cacheRepo.GetCategories(s.ctx, req.Id)
		if err == nil {
			// 分类信息也在缓存中
			categoryNames := make([]string, len(categories))
			for i, category := range categories {
				categoryNames[i] = category.Name
			}
			productModel.Categories = categoryNames
			resp.Product = converter.ToProductProto(productModel)
			klog.CtxInfof(s.ctx, "从缓存获取商品信息成功 - productId: %d", req.Id)
			return resp, nil
		}
	}

	// 2. 缓存未命中，从数据库获取商品信息
	productModel, err = s.productRepo.GetProduct(req.Id)
	if err != nil {
		klog.CtxErrorf(s.ctx, "获取商品信息失败 - productId: %d, err: %v", req.Id, err)
		return nil, utils.NewBizError(40004, "商品不存在")
	}

	// 3. 获取商品分类信息
	categories, err := s.productRepo.GetCategories(req.Id)
	if err != nil {
		klog.CtxWarnf(s.ctx, "获取商品分类信息失败 - productId: %d, err: %v", req.Id, err)
	} else {
		categoryNames := make([]string, len(categories))
		for i, category := range categories {
			categoryNames[i] = category.Name
		}
		productModel.Categories = categoryNames
	}

	// 4. 更新缓存
	if err := s.cacheRepo.SetProduct(s.ctx, productModel); err != nil {
		klog.CtxWarnf(s.ctx, "更新商品缓存失败 - productId: %d, err: %v", req.Id, err)
	}

	resp.Product = converter.ToProductProto(productModel)
	klog.CtxInfof(s.ctx, "获取商品信息成功 - productId: %d", req.Id)
	return resp, nil
}
