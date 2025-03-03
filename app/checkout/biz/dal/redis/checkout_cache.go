package redis

import (
	"fmt"
	"time"

	"github.com/trashwbin/dymall/app/checkout/biz/model"
)

const (
	// 结算单缓存key前缀
	CheckoutKeyPrefix = "checkout:"
	// 结算单商品缓存key前缀
	CheckoutItemKeyPrefix = "checkout:item:"
	// 结算单缓存过期时间
	CheckoutExpiration = 24 * time.Hour
)

// CheckoutCache Redis结算单缓存数据对象
type CheckoutCache struct {
	ID          string    `json:"id"`
	UserID      int64     `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Currency    string    `json:"currency"`
	Status      int       `json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
	ExpireAt    time.Time `json:"expire_at"`
}

// CheckoutItemCache Redis结算单商品缓存数据对象
type CheckoutItemCache struct {
	ID         int64     `json:"id"`
	CheckoutID string    `json:"checkout_id"`
	ProductID  int64     `json:"product_id"`
	Quantity   int32     `json:"quantity"`
	Price      float64   `json:"price"`
	Currency   string    `json:"currency"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ToModel 转换为领域模型
func (c *CheckoutCache) ToModel() *model.Checkout {
	return &model.Checkout{
		ID:          c.ID,
		UserID:      c.UserID,
		TotalAmount: c.TotalAmount,
		Currency:    c.Currency,
		Status:      model.CheckoutStatus(c.Status),
		UpdatedAt:   c.UpdatedAt,
		ExpireAt:    c.ExpireAt,
	}
}

// FromModel 从领域模型转换
func (c *CheckoutCache) FromModel(m *model.Checkout) {
	c.ID = m.ID
	c.UserID = m.UserID
	c.TotalAmount = m.TotalAmount
	c.Currency = m.Currency
	c.Status = int(m.Status)
	c.UpdatedAt = m.UpdatedAt
	c.ExpireAt = m.ExpireAt
}

// GetKey 获取结算单缓存key
func (c *CheckoutCache) GetKey() string {
	return fmt.Sprintf("%s%s", CheckoutKeyPrefix, c.ID)
}

// ToModel 转换为领域模型
func (ci *CheckoutItemCache) ToModel() (*model.CheckoutItem, error) {
	return &model.CheckoutItem{
		ID:        ci.ID,
		ProductID: ci.ProductID,
		Quantity:  ci.Quantity,
		Price:     ci.Price,
		Currency:  ci.Currency,
	}, nil
}

// FromModel 从领域模型转换
func (ci *CheckoutItemCache) FromModel(m *model.CheckoutItem) error {
	ci.ID = m.ID
	ci.ProductID = m.ProductID
	ci.Quantity = m.Quantity
	ci.Price = m.Price
	ci.Currency = m.Currency
	return nil
}

// GetKey 获取结算单商品缓存key
func (ci *CheckoutItemCache) GetKey() string {
	return fmt.Sprintf("%s%s:%d", CheckoutItemKeyPrefix, ci.CheckoutID, ci.ProductID)
}

// GetCheckoutItemPattern 获取结算单所有商品的key模式
func GetCheckoutItemPattern(checkoutID string) string {
	return fmt.Sprintf("%s%s:*", CheckoutItemKeyPrefix, checkoutID)
}
