package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/trashwbin/dymall/app/cart/biz/model"
)

type CartRepo struct{}

// NewCartRepo 创建购物车Redis仓储
func NewCartRepo() *CartRepo {
	return &CartRepo{}
}

// GetCart 根据用户ID获取购物车（仅用于首次获取购物车）
func (r *CartRepo) GetCart(ctx context.Context, userID int64) (*model.Cart, error) {
	cache := &CartCache{UserID: userID}
	data, err := RedisClient.Get(ctx, cache.GetKey()).Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("unmarshal cart cache failed: %w", err)
	}

	return cache.ToModel(), nil
}

// SetCart 设置购物车缓存
func (r *CartRepo) SetCart(ctx context.Context, cart *model.Cart) error {
	cache := &CartCache{}
	cache.FromModel(cart)

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("marshal cart cache failed: %w", err)
	}

	return RedisClient.Set(ctx, cache.GetKey(), data, CartExpiration).Err()
}

// GetCartItem 获取购物车商品
func (r *CartRepo) GetCartItem(ctx context.Context, cartID, productID int64) (*model.CartItem, error) {
	cache := &CartItemCache{CartID: cartID, ProductID: productID}
	data, err := RedisClient.Get(ctx, cache.GetKey()).Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("unmarshal cart item cache failed: %w", err)
	}

	return cache.ToModel()
}

// SetCartItem 设置购物车商品缓存
func (r *CartRepo) SetCartItem(ctx context.Context, item *model.CartItem) error {
	cache := &CartItemCache{}
	if err := cache.FromModel(item); err != nil {
		return fmt.Errorf("convert model to cache failed: %w", err)
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("marshal cart item cache failed: %w", err)
	}

	return RedisClient.Set(ctx, cache.GetKey(), data, CartExpiration).Err()
}

// DeleteCartItem 删除购物车商品缓存
func (r *CartRepo) DeleteCartItem(ctx context.Context, cartID, productID int64) error {
	cache := &CartItemCache{CartID: cartID, ProductID: productID}
	return RedisClient.Del(ctx, cache.GetKey()).Err()
}

// GetCartItems 获取购物车商品列表
func (r *CartRepo) GetCartItems(ctx context.Context, cartID int64) ([]*model.CartItem, error) {
	pattern := GetCartItemPattern(cartID)
	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("get cart items keys failed: %w", err)
	}

	var items []*model.CartItem
	for _, key := range keys {
		data, err := RedisClient.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		cache := &CartItemCache{}
		if err := json.Unmarshal(data, cache); err != nil {
			continue
		}

		item, err := cache.ToModel()
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	// 按照ProductID排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].ProductID < items[j].ProductID
	})

	return items, nil
}

// EmptyCart 清空购物车缓存
func (r *CartRepo) EmptyCart(ctx context.Context, cartID int64) error {
	pattern := GetCartItemPattern(cartID)
	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("get cart items keys failed: %w", err)
	}

	if len(keys) > 0 {
		if err := RedisClient.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("delete cart items cache failed: %w", err)
		}
	}

	return nil
}
