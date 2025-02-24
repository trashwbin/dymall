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

type CreateProductService struct {
	ctx         context.Context
	productRepo *mysql.ProductRepo
	cacheRepo   *redis.ProductRepo
}

// NewCreateProductService new CreateProductService
func NewCreateProductService(ctx context.Context) *CreateProductService {
	return &CreateProductService{
		ctx:         ctx,
		productRepo: mysql.NewProductRepo(),
		cacheRepo:   redis.NewProductRepo(),
	}
}

// Run 创建商品
func (s *CreateProductService) Run(req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	resp = new(product.CreateProductResp)

	// 1. 参数验证
	if req.Name == "" {
		return nil, utils.NewBizError(40001, "商品名称不能为空")
	}
	if req.Price < 0 {
		return nil, utils.NewBizError(40002, "商品价格不能为负数")
	}

	// 2. 构建商品模型
	now := time.Now()
	productModel := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  req.Categories,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 3. 验证商品信息
	if !productModel.IsValid() {
		return nil, utils.NewBizError(40003, "商品信息不合法")
	}

	// 4. 使用事务创建商品
	err = s.productRepo.Transaction(func(txRepo *mysql.ProductRepo) error {
		// 4.1 创建商品
		var err error
		productModel, err = txRepo.CreateProduct(productModel)
		if err != nil {
			klog.CtxErrorf(s.ctx, "创建商品失败 - name: %s, err: %v", req.Name, err)
			return utils.NewBizError(50001, "创建商品失败")
		}

		// 4.2 处理商品分类
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

	// 5. 更新缓存
	if err := s.cacheRepo.SetProduct(s.ctx, productModel); err != nil {
		klog.CtxWarnf(s.ctx, "更新商品缓存失败 - productId: %d, err: %v", productModel.ID, err)
	}

	// 6. 删除商品列表缓存（因为新增了商品，所以需要清理列表缓存）
	if err := s.cacheRepo.DeleteProductList(s.ctx, 1, 10, ""); err != nil {
		klog.CtxWarnf(s.ctx, "删除商品列表缓存失败 - err: %v", err)
	}

	resp.Product = converter.ToProductProto(productModel)
	klog.CtxInfof(s.ctx, "创建商品成功 - productId: %d, name: %s", productModel.ID, req.Name)
	return resp, nil
}
