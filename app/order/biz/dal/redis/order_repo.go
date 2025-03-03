package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/trashwbin/dymall/app/order/biz/model"
)

type OrderRepo struct{}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{}
}

// GetOrder 获取订单缓存
func (r *OrderRepo) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	cache := &OrderCache{OrderID: orderID}
	data, err := RedisClient.Get(ctx, cache.GetKey()).Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("unmarshal order cache failed: %w", err)
	}

	return cache.ToModel(), nil
}

// SetOrder 设置订单缓存
func (r *OrderRepo) SetOrder(ctx context.Context, order *model.Order) error {
	cache := &OrderCache{}
	cache.FromModel(order)

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("marshal order cache failed: %w", err)
	}

	return RedisClient.Set(ctx, cache.GetKey(), data, OrderExpiration).Err()
}

// DeleteOrder 删除订单缓存
func (r *OrderRepo) DeleteOrder(ctx context.Context, orderID string) error {
	cache := &OrderCache{OrderID: orderID}
	return RedisClient.Del(ctx, cache.GetKey()).Err()
}

// GetUserOrders 获取用户订单列表缓存
func (r *OrderRepo) GetUserOrders(ctx context.Context, userID int64) ([]*model.Order, error) {
	key := GetUserOrdersKey(userID)
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var orderIDs []string
	if err := json.Unmarshal(data, &orderIDs); err != nil {
		return nil, fmt.Errorf("unmarshal user orders cache failed: %w", err)
	}

	orders := make([]*model.Order, 0, len(orderIDs))
	for _, orderID := range orderIDs {
		order, err := r.GetOrder(ctx, orderID)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// SetUserOrders 设置用户订单列表缓存
func (r *OrderRepo) SetUserOrders(ctx context.Context, userID int64, orderIDs []string) error {
	key := GetUserOrdersKey(userID)
	data, err := json.Marshal(orderIDs)
	if err != nil {
		return fmt.Errorf("marshal user orders cache failed: %w", err)
	}

	return RedisClient.Set(ctx, key, data, OrderExpiration).Err()
}

// DeleteUserOrders 删除用户订单列表缓存
func (r *OrderRepo) DeleteUserOrders(ctx context.Context, userID int64) error {
	key := GetUserOrdersKey(userID)
	return RedisClient.Del(ctx, key).Err()
}
