package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/trashwbin/dymall/app/checkout/biz/model"
)

type CheckoutRepo struct{}

// NewCheckoutRepo 创建结算单Redis仓储
func NewCheckoutRepo() *CheckoutRepo {
	return &CheckoutRepo{}
}

// GetCheckout 获取结算单
func (r *CheckoutRepo) GetCheckout(ctx context.Context, id string) (*model.Checkout, error) {
	cache := &CheckoutCache{ID: id}
	data, err := RedisClient.Get(ctx, cache.GetKey()).Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("unmarshal checkout cache failed: %w", err)
	}

	checkout := cache.ToModel()

	// 获取结算单商品
	items, err := r.GetCheckoutItems(ctx, id)
	if err != nil {
		return nil, err
	}
	checkout.Items = items

	return checkout, nil
}

// SetCheckout 设置结算单缓存
func (r *CheckoutRepo) SetCheckout(ctx context.Context, checkout *model.Checkout) error {
	cache := &CheckoutCache{}
	cache.FromModel(checkout)

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("marshal checkout cache failed: %w", err)
	}

	// 设置结算单缓存
	if err := RedisClient.Set(ctx, cache.GetKey(), data, CheckoutExpiration).Err(); err != nil {
		return err
	}

	// 设置结算单商品缓存
	for _, item := range checkout.Items {
		if err := r.SetCheckoutItem(ctx, checkout.ID, item); err != nil {
			return err
		}
	}

	return nil
}

// GetCheckoutItem 获取结算单商品
func (r *CheckoutRepo) GetCheckoutItem(ctx context.Context, checkoutID string, productID int64) (*model.CheckoutItem, error) {
	cache := &CheckoutItemCache{CheckoutID: checkoutID, ProductID: productID}
	data, err := RedisClient.Get(ctx, cache.GetKey()).Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("unmarshal checkout item cache failed: %w", err)
	}

	return cache.ToModel()
}

// SetCheckoutItem 设置结算单商品缓存
func (r *CheckoutRepo) SetCheckoutItem(ctx context.Context, checkoutID string, item *model.CheckoutItem) error {
	cache := &CheckoutItemCache{CheckoutID: checkoutID}
	if err := cache.FromModel(item); err != nil {
		return fmt.Errorf("convert model to cache failed: %w", err)
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("marshal checkout item cache failed: %w", err)
	}

	return RedisClient.Set(ctx, cache.GetKey(), data, CheckoutExpiration).Err()
}

// DeleteCheckoutItem 删除结算单商品缓存
func (r *CheckoutRepo) DeleteCheckoutItem(ctx context.Context, checkoutID string, productID int64) error {
	cache := &CheckoutItemCache{CheckoutID: checkoutID, ProductID: productID}
	return RedisClient.Del(ctx, cache.GetKey()).Err()
}

// GetCheckoutItems 获取结算单商品列表
func (r *CheckoutRepo) GetCheckoutItems(ctx context.Context, checkoutID string) ([]*model.CheckoutItem, error) {
	pattern := GetCheckoutItemPattern(checkoutID)
	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("get checkout items keys failed: %w", err)
	}

	var items []*model.CheckoutItem
	for _, key := range keys {
		data, err := RedisClient.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		cache := &CheckoutItemCache{}
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

// DeleteCheckout 删除结算单缓存
func (r *CheckoutRepo) DeleteCheckout(ctx context.Context, id string) error {
	// 删除结算单商品缓存
	pattern := GetCheckoutItemPattern(id)
	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("get checkout items keys failed: %w", err)
	}

	if len(keys) > 0 {
		if err := RedisClient.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("delete checkout items cache failed: %w", err)
		}
	}

	// 删除结算单缓存
	cache := &CheckoutCache{ID: id}
	return RedisClient.Del(ctx, cache.GetKey()).Err()
}
