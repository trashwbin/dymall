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

type ListProductsService struct {
	ctx         context.Context
	productRepo *mysql.ProductRepo
	cacheRepo   *redis.ProductRepo
}

// NewListProductsService new ListProductsService
func NewListProductsService(ctx context.Context) *ListProductsService {
	return &ListProductsService{
		ctx:         ctx,
		productRepo: mysql.NewProductRepo(),
		cacheRepo:   redis.NewProductRepo(),
	}
}

// Run 获取商品列表
func (s *ListProductsService) Run(req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	resp = new(product.ListProductsResp)

	// 参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 1. 尝试从缓存获取商品列表
	products, err := s.cacheRepo.GetProductList(s.ctx, req.Page, req.PageSize, req.CategoryName)
	if err == nil {
		// 缓存命中
		resp.Products = converter.ToProductProtos(products)
		klog.CtxInfof(s.ctx, "从缓存获取商品列表成功 - page: %d, pageSize: %d, category: %s",
			req.Page, req.PageSize, req.CategoryName)
		return resp, nil
	}

	// 2. 缓存未命中，从数据库获取
	products, err = s.productRepo.ListProducts(req.Page, req.PageSize, req.CategoryName)
	if err != nil {
		klog.CtxErrorf(s.ctx, "获取商品列表失败 - page: %d, pageSize: %d, category: %s, err: %v",
			req.Page, req.PageSize, req.CategoryName, err)
		return nil, utils.NewBizError(50001, "获取商品列表失败")
	}

	// 3. 更新缓存
	if err := s.cacheRepo.SetProductList(s.ctx, products, req.Page, req.PageSize, req.CategoryName); err != nil {
		klog.CtxWarnf(s.ctx, "更新商品列表缓存失败 - page: %d, pageSize: %d, category: %s, err: %v",
			req.Page, req.PageSize, req.CategoryName, err)
	}

	// 4. 转换并返回结果
	resp.Products = converter.ToProductProtos(products)

	klog.CtxInfof(s.ctx, "获取商品列表成功 - page: %d, pageSize: %d, category: %s, count: %d",
		req.Page, req.PageSize, req.CategoryName, len(products))
	return resp, nil
}
