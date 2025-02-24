package redis

import (
	"fmt"
	"time"

	"github.com/trashwbin/dymall/app/payment/biz/model"
)

const (
	// 支付单缓存key前缀
	PaymentKeyPrefix = "payment:id:"
	// 订单支付单缓存key前缀
	OrderPaymentKeyPrefix = "payment:order:"
	// 用户支付单列表缓存key前缀
	UserPaymentKeyPrefix = "payment:user:"
	// 支付单缓存过期时间
	PaymentExpiration = 24 * time.Hour
)

// PaymentCache Redis支付单缓存数据对象
type PaymentCache struct {
	ID        int64     `json:"id"`
	PaymentID string    `json:"payment_id"`
	OrderID   string    `json:"order_id"`
	UserID    int64     `json:"user_id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpireAt  time.Time `json:"expire_at"`
	PaidAt    time.Time `json:"paid_at"`
}

// ToModel 转换为领域模型
func (c *PaymentCache) ToModel() *model.Payment {
	return &model.Payment{
		ID:        c.ID,
		PaymentID: c.PaymentID,
		OrderID:   c.OrderID,
		UserID:    c.UserID,
		Amount:    c.Amount,
		Currency:  c.Currency,
		Status:    model.PayStatus(c.Status),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		ExpireAt:  c.ExpireAt,
		PaidAt:    c.PaidAt,
	}
}

// FromModel 从领域模型转换
func (c *PaymentCache) FromModel(m *model.Payment) {
	c.ID = m.ID
	c.PaymentID = m.PaymentID
	c.OrderID = m.OrderID
	c.UserID = m.UserID
	c.Amount = m.Amount
	c.Currency = m.Currency
	c.Status = int(m.Status)
	c.CreatedAt = m.CreatedAt
	c.UpdatedAt = m.UpdatedAt
	c.ExpireAt = m.ExpireAt
	c.PaidAt = m.PaidAt
}

// GetPaymentKey 获取支付单缓存key
func GetPaymentKey(paymentID string) string {
	return fmt.Sprintf("%s%s", PaymentKeyPrefix, paymentID)
}

// GetOrderPaymentKey 获取订单支付单缓存key
func GetOrderPaymentKey(orderID string) string {
	return fmt.Sprintf("%s%s", OrderPaymentKeyPrefix, orderID)
}

// GetUserPaymentKey 获取用户支付单列表缓存key
func GetUserPaymentKey(userID int64) string {
	return fmt.Sprintf("%s%d", UserPaymentKeyPrefix, userID)
}
