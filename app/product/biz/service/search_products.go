package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/product/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/product/biz/service/converter"
	"github.com/trashwbin/dymall/app/product/utils"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

type SearchProductsService struct {
	ctx         context.Context
	productRepo *mysql.ProductRepo
}

// NewSearchProductsService new SearchProductsService
func NewSearchProductsService(ctx context.Context) *SearchProductsService {
	return &SearchProductsService{
		ctx:         ctx,
		productRepo: mysql.NewProductRepo(),
	}
}

// Run 搜索商品
func (s *SearchProductsService) Run(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	resp = new(product.SearchProductsResp)

	// 参数验证
	if req.Query == "" {
		return nil, utils.NewBizError(40001, "搜索关键词不能为空")
	}

	// 从数据库搜索商品
	products, err := s.productRepo.SearchProducts(req.Query)
	if err != nil {
		klog.CtxErrorf(s.ctx, "搜索商品失败 - query: %s, err: %v", req.Query, err)
		return nil, utils.NewBizError(50001, "搜索商品失败")
	}

	// 转换并返回结果
	resp.Results = converter.ToProductProtos(products)

	klog.CtxInfof(s.ctx, "搜索商品成功 - query: %s, count: %d", req.Query, len(products))
	return resp, nil
}
