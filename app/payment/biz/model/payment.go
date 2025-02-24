package model

import (
	"fmt"
	"time"

	creditcard "github.com/durango/go-credit-card"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

// Payment 支付领域模型
type Payment struct {
	ID        int64
	PaymentID string    // 支付单号
	OrderID   string    // 订单号
	UserID    int64     // 用户ID
	Amount    float64   // 支付金额
	Currency  string    // 货币类型
	Status    PayStatus // 支付状态
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
	ExpireAt  time.Time // 过期时间
	PaidAt    time.Time // 支付时间
}

// PayStatus 支付状态
type PayStatus int32

const (
	PayStatusPending   PayStatus = PayStatus(payment.PaymentStatus_PAYMENT_STATUS_PENDING)   // 待支付
	PayStatusSuccess   PayStatus = PayStatus(payment.PaymentStatus_PAYMENT_STATUS_SUCCESS)   // 支付成功
	PayStatusFailed    PayStatus = PayStatus(payment.PaymentStatus_PAYMENT_STATUS_FAILED)    // 支付失败
	PayStatusCancelled PayStatus = PayStatus(payment.PaymentStatus_PAYMENT_STATUS_CANCELLED) // 已取消
	PayStatusExpired   PayStatus = PayStatus(payment.PaymentStatus_PAYMENT_STATUS_EXPIRED)   // 已过期
)

// CreditCard 信用卡信息
type CreditCard struct {
	Number          string
	CVV             int32
	ExpirationYear  int32
	ExpirationMonth int32
}

// IsValid 验证支付状态是否有效
func (p *Payment) IsValid() bool {
	return p.Status == PayStatusPending
}

// IsPaid 是否已支付
func (p *Payment) IsPaid() bool {
	return p.Status == PayStatusSuccess
}

// IsExpired 是否已过期
func (p *Payment) IsExpired() bool {
	return p.Status == PayStatusExpired || time.Now().After(p.ExpireAt)
}

// CanPay 是否可以支付
func (p *Payment) CanPay() bool {
	return p.Status == PayStatusPending && !p.IsExpired()
}

// ValidateCreditCard 验证信用卡信息
func (c *CreditCard) ValidateCreditCard() (bool, string) {
	// 1. 基本格式验证
	if len(c.Number) < 13 || len(c.Number) > 19 {
		return false, "invalid card number length"
	}

	// 2. 验证卡号是否只包含数字
	for _, r := range c.Number {
		if r < '0' || r > '9' {
			return false, "card number must contain only digits"
		}
	}

	// 3. 创建信用卡对象
	card := creditcard.Card{
		Number: c.Number,
		Cvv:    fmt.Sprintf("%d", c.CVV),
		Month:  fmt.Sprintf("%d", c.ExpirationMonth),
		Year:   fmt.Sprintf("%d", c.ExpirationYear),
	}

	// 4. 验证所有信息
	if err := card.Validate(); err != nil {
		return false, err.Error()
	}

	// 5. 额外验证卡号合法性
	if !card.ValidateNumber() {
		return false, "invalid card number"
	}

	return true, ""
}
