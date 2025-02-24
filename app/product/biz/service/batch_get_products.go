package service

import (
	"context"

	"github.com/trashwbin/dymall/app/product/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/product/biz/dal/redis"
	"github.com/trashwbin/dymall/app/product/biz/model"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

type BatchGetProductsService struct {
	ctx context.Context
} // NewBatchGetProductsService new BatchGetProductsService
func NewBatchGetProductsService(ctx context.Context) *BatchGetProductsService {
	return &BatchGetProductsService{ctx: ctx}
}

// Run 批量获取商品信息
func (s *BatchGetProductsService) Run(req *product.BatchGetProductsReq) (resp *product.BatchGetProductsResp, err error) {
	resp = new(product.BatchGetProductsResp)

	// 参数校验
	if len(req.Ids) == 0 {
		return resp, nil
	}

	// 去重ID
	uniqueIDs := make([]uint32, 0, len(req.Ids))
	idMap := make(map[uint32]bool)
	for _, id := range req.Ids {
		if !idMap[id] {
			idMap[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	// 1. 先从Redis缓存中获取
	redisRepo := redis.NewProductRepo()
	cacheProducts, cacheMissingIDs, err := redisRepo.BatchGetProducts(s.ctx, uniqueIDs)
	if err != nil {
		// 缓存错误不影响主流程，继续从数据库查询
		cacheMissingIDs = uniqueIDs
	}

	// 如果缓存全部命中，直接返回
	if len(cacheMissingIDs) == 0 {
		resp.Products = convertToProtoProducts(cacheProducts)
		return resp, nil
	}

	// 2. 从MySQL中查询缓存未命中的商品
	mysqlRepo := mysql.NewProductRepo()
	dbProducts, dbMissingIDs, err := mysqlRepo.BatchGetProducts(cacheMissingIDs)
	if err != nil {
		return nil, err
	}

	// 3. 将数据库查询结果更新到缓存
	if len(dbProducts) > 0 {
		go func() {
			_ = redisRepo.BatchSetProducts(context.Background(), dbProducts)
		}()
	}

	// 4. 合并缓存和数据库的结果
	allProducts := append(cacheProducts, dbProducts...)
	resp.Products = convertToProtoProducts(allProducts)
	resp.MissingIds = dbMissingIDs

	return resp, nil
}

// convertToProtoProducts 将领域模型转换为proto消息
func convertToProtoProducts(products []*model.Product) []*product.Product {
	result := make([]*product.Product, len(products))
	for i, p := range products {
		result[i] = &product.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Picture:     p.Picture,
			Price:       p.Price,
			Categories:  p.Categories,
		}
	}
	return result
}
