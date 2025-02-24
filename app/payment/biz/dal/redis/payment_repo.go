package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/trashwbin/dymall/app/payment/biz/model"
)

type PaymentRepo struct{}

func NewPaymentRepo() *PaymentRepo {
	return &PaymentRepo{}
}

// GetPayment 获取支付单缓存
func (r *PaymentRepo) GetPayment(ctx context.Context, paymentID string) (*model.Payment, error) {
	data, err := RedisClient.Get(ctx, GetPaymentKey(paymentID)).Bytes()
	if err != nil {
		return nil, err
	}

	cache := &PaymentCache{}
	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("unmarshal payment cache failed: %w", err)
	}

	return cache.ToModel(), nil
}

// SetPayment 设置支付单缓存
func (r *PaymentRepo) SetPayment(ctx context.Context, payment *model.Payment) error {
	cache := &PaymentCache{}
	cache.FromModel(payment)

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("marshal payment cache failed: %w", err)
	}

	// 设置支付单缓存
	if err := RedisClient.Set(ctx, GetPaymentKey(payment.PaymentID), data, PaymentExpiration).Err(); err != nil {
		return err
	}

	// 设置订单支付单缓存
	if err := RedisClient.Set(ctx, GetOrderPaymentKey(payment.OrderID), data, PaymentExpiration).Err(); err != nil {
		return err
	}

	return nil
}

// DeletePayment 删除支付单缓存
func (r *PaymentRepo) DeletePayment(ctx context.Context, payment *model.Payment) error {
	// 删除支付单缓存
	if err := RedisClient.Del(ctx, GetPaymentKey(payment.PaymentID)).Err(); err != nil {
		return err
	}

	// 删除订单支付单缓存
	if err := RedisClient.Del(ctx, GetOrderPaymentKey(payment.OrderID)).Err(); err != nil {
		return err
	}

	return nil
}

// GetPaymentByOrderID 根据订单号获取支付单缓存
func (r *PaymentRepo) GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	data, err := RedisClient.Get(ctx, GetOrderPaymentKey(orderID)).Bytes()
	if err != nil {
		return nil, err
	}

	cache := &PaymentCache{}
	if err := json.Unmarshal(data, cache); err != nil {
		return nil, fmt.Errorf("unmarshal payment cache failed: %w", err)
	}

	return cache.ToModel(), nil
}

// SetUserPaymentIDs 设置用户支付单ID列表缓存
func (r *PaymentRepo) SetUserPaymentIDs(ctx context.Context, userID int64, paymentIDs []string) error {
	data, err := json.Marshal(paymentIDs)
	if err != nil {
		return fmt.Errorf("marshal payment ids failed: %w", err)
	}

	return RedisClient.Set(ctx, GetUserPaymentKey(userID), data, PaymentExpiration).Err()
}

// GetUserPaymentIDs 获取用户支付单ID列表缓存
func (r *PaymentRepo) GetUserPaymentIDs(ctx context.Context, userID int64) ([]string, error) {
	data, err := RedisClient.Get(ctx, GetUserPaymentKey(userID)).Bytes()
	if err != nil {
		return nil, err
	}

	var paymentIDs []string
	if err := json.Unmarshal(data, &paymentIDs); err != nil {
		return nil, fmt.Errorf("unmarshal payment ids failed: %w", err)
	}

	return paymentIDs, nil
}

// DeleteUserPaymentIDs 删除用户支付单ID列表缓存
func (r *PaymentRepo) DeleteUserPaymentIDs(ctx context.Context, userID int64) error {
	return RedisClient.Del(ctx, GetUserPaymentKey(userID)).Err()
}
