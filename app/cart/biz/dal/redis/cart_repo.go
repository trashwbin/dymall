package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/trashwbin/dymall/app/cart/biz/model"
)

type CartRepo struct{}

// NewCartRepo 创建购物车Redis仓储
func NewCartRepo() *CartRepo {
	return &CartRepo{}
}

// SetCart 设置购物车缓存
func (r *CartRepo) SetCart(ctx context.Context, cart *model.Cart) error {
	var cache CartCache
	cache.FromModel(cart)

	key := GetCartKey(cart.UserID)
	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("marshal cart cache failed: %w", err)
	}

	return RedisClient.Set(ctx, key, data, CartExpiration).Err()
}

// GetCart 获取购物车缓存
func (r *CartRepo) GetCart(ctx context.Context, userID int64) (*model.Cart, error) {
	key := GetCartKey(userID)
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var cache CartCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("unmarshal cart cache failed: %w", err)
	}

	return cache.ToModel(), nil
}

// SetCartItem 设置购物车商品缓存
func (r *CartRepo) SetCartItem(ctx context.Context, item *model.CartItem) error {
	var cache CartItemCache
	if err := cache.FromModel(item); err != nil {
		return fmt.Errorf("convert model to cache failed: %w", err)
	}

	key := GetCartItemKey(item.CartID, item.ProductID)
	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("marshal cart item cache failed: %w", err)
	}

	return RedisClient.Set(ctx, key, data, CartExpiration).Err()
}

// GetCartItem 获取购物车商品缓存
func (r *CartRepo) GetCartItem(ctx context.Context, cartID, productID int64) (*model.CartItem, error) {
	key := GetCartItemKey(cartID, productID)
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var cache CartItemCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("unmarshal cart item cache failed: %w", err)
	}

	return cache.ToModel()
}

// DeleteCartItem 删除购物车商品缓存
func (r *CartRepo) DeleteCartItem(ctx context.Context, cartID, productID int64) error {
	key := GetCartItemKey(cartID, productID)
	return RedisClient.Del(ctx, key).Err()
}

// GetCartItems 获取用户购物车所有商品
func (r *CartRepo) GetCartItems(ctx context.Context, cartID int64) ([]*model.CartItem, error) {
	pattern := fmt.Sprintf("%s%d:*", CartItemKeyPrefix, cartID)
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

		var cache CartItemCache
		if err := json.Unmarshal(data, &cache); err != nil {
			continue
		}

		item, err := cache.ToModel()
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, nil
}
