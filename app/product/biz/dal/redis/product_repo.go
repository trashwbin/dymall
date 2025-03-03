package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/trashwbin/dymall/app/product/biz/model"
)

type ProductRepo struct{}

// NewProductRepo 创建商品Redis仓储
func NewProductRepo() *ProductRepo {
	return &ProductRepo{}
}

// GetProduct 获取商品缓存
func (r *ProductRepo) GetProduct(ctx context.Context, id uint32) (*model.Product, error) {
	cache := &ProductCache{ID: id}
	data, err := RedisClient.Get(ctx, cache.GetKey()).Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("unmarshal product cache failed: %w", err)
	}

	return cache.ToModel(), nil
}

// SetProduct 设置商品缓存
func (r *ProductRepo) SetProduct(ctx context.Context, product *model.Product) error {
	cache := &ProductCache{}
	cache.FromModel(product)

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("marshal product cache failed: %w", err)
	}

	return RedisClient.Set(ctx, cache.GetKey(), data, ProductExpiration).Err()
}

// DeleteProduct 删除商品缓存
func (r *ProductRepo) DeleteProduct(ctx context.Context, id uint32) error {
	cache := &ProductCache{ID: id}
	return RedisClient.Del(ctx, cache.GetKey()).Err()
}

// GetProductList 获取商品列表缓存
func (r *ProductRepo) GetProductList(ctx context.Context, page int32, pageSize int64, categoryName string) ([]*model.Product, error) {
	key := GetProductListKey(page, pageSize, categoryName)
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var caches []*ProductCache
	if err := json.Unmarshal(data, &caches); err != nil {
		return nil, fmt.Errorf("unmarshal product list cache failed: %w", err)
	}

	products := make([]*model.Product, len(caches))
	for i, cache := range caches {
		products[i] = cache.ToModel()
	}
	return products, nil
}

// SetProductList 设置商品列表缓存
func (r *ProductRepo) SetProductList(ctx context.Context, products []*model.Product, page int32, pageSize int64, categoryName string) error {
	caches := make([]*ProductCache, len(products))
	for i, product := range products {
		cache := &ProductCache{}
		cache.FromModel(product)
		caches[i] = cache
	}

	data, err := json.Marshal(caches)
	if err != nil {
		return fmt.Errorf("marshal product list cache failed: %w", err)
	}

	key := GetProductListKey(page, pageSize, categoryName)
	return RedisClient.Set(ctx, key, data, ProductExpiration).Err()
}

// DeleteProductList 删除商品列表缓存
func (r *ProductRepo) DeleteProductList(ctx context.Context, page int32, pageSize int64, categoryName string) error {
	key := GetProductListKey(page, pageSize, categoryName)
	return RedisClient.Del(ctx, key).Err()
}

// GetCategory 获取分类缓存
func (r *ProductRepo) GetCategory(ctx context.Context, id uint32) (*model.Category, error) {
	cache := &CategoryCache{ID: id}
	data, err := RedisClient.Get(ctx, cache.GetKey()).Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("unmarshal category cache failed: %w", err)
	}

	return cache.ToModel(), nil
}

// SetCategory 设置分类缓存
func (r *ProductRepo) SetCategory(ctx context.Context, category *model.Category) error {
	cache := &CategoryCache{}
	cache.FromModel(category)

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("marshal category cache failed: %w", err)
	}

	return RedisClient.Set(ctx, cache.GetKey(), data, ProductExpiration).Err()
}

// DeleteCategory 删除分类缓存
func (r *ProductRepo) DeleteCategory(ctx context.Context, id uint32) error {
	cache := &CategoryCache{ID: id}
	return RedisClient.Del(ctx, cache.GetKey()).Err()
}

// GetCategories 获取商品分类列表缓存
func (r *ProductRepo) GetCategories(ctx context.Context, productID uint32) ([]*model.Category, error) {
	// 获取商品缓存，从中提取分类信息
	product, err := r.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	// 根据分类名称获取分类信息
	categories := make([]*model.Category, len(product.Categories))
	for i, categoryName := range product.Categories {
		categories[i] = &model.Category{
			Name: categoryName,
		}
	}
	return categories, nil
}

// BatchGetProducts 批量获取商品缓存
func (r *ProductRepo) BatchGetProducts(ctx context.Context, ids []uint32) ([]*model.Product, []uint32, error) {
	// 构建所有商品的缓存key
	keys := make([]string, len(ids))
	idToIndex := make(map[uint32]int)
	for i, id := range ids {
		keys[i] = fmt.Sprintf("%s%d", ProductKeyPrefix, id)
		idToIndex[id] = i
	}

	// 批量获取缓存
	results, err := RedisClient.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("batch get products from cache failed: %w", err)
	}

	products := make([]*model.Product, 0, len(ids))
	missingIDs := make([]uint32, 0)

	// 处理结果
	for i, result := range results {
		if result == nil {
			missingIDs = append(missingIDs, ids[i])
			continue
		}

		data, ok := result.(string)
		if !ok {
			missingIDs = append(missingIDs, ids[i])
			continue
		}

		var cache ProductCache
		if err := json.Unmarshal([]byte(data), &cache); err != nil {
			missingIDs = append(missingIDs, ids[i])
			continue
		}

		products = append(products, cache.ToModel())
	}

	return products, missingIDs, nil
}

// BatchSetProducts 批量设置商品缓存
func (r *ProductRepo) BatchSetProducts(ctx context.Context, products []*model.Product) error {
	if len(products) == 0 {
		return nil
	}

	// 使用管道批量设置缓存
	pipe := RedisClient.Pipeline()
	for _, product := range products {
		cache := &ProductCache{}
		cache.FromModel(product)

		data, err := json.Marshal(cache)
		if err != nil {
			return fmt.Errorf("marshal product cache failed: %w", err)
		}

		pipe.Set(ctx, cache.GetKey(), data, ProductExpiration)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("batch set products to cache failed: %w", err)
	}

	return nil
}
