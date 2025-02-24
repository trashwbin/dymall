package model

import "time"

// CheckoutStatus 结算单状态
type CheckoutStatus int

const (
	CheckoutStatusPending   CheckoutStatus = 1 // 待提交
	CheckoutStatusSubmitted CheckoutStatus = 2 // 已提交
	CheckoutStatusExpired   CheckoutStatus = 3 // 已过期
)

// Checkout 结算单领域模型
type Checkout struct {
	ID          string
	UserID      int64
	TotalAmount float64
	Currency    string
	Status      CheckoutStatus
	Items       []*CheckoutItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpireAt    time.Time
}

// CheckoutItem 结算单商品领域模型
type CheckoutItem struct {
	ID        int64
	ProductID int64
	Quantity  int32
	Price     float64
	Currency  string
}

// Address 地址领域模型
type Address struct {
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       string
}

// IsValid 验证结算单状态是否有效
func (c *Checkout) IsValid() bool {
	return c.Status == CheckoutStatusPending
}

// IsExpired 检查结算单是否过期
func (c *Checkout) IsExpired() bool {
	return time.Now().After(c.ExpireAt)
}

// CalculateTotalAmount 计算总金额
func (c *Checkout) CalculateTotalAmount() float64 {
	var total float64
	for _, item := range c.Items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}
