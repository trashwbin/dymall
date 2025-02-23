package model

import "time"

// Cart 购物车领域模型
type Cart struct {
	ID        int64
	UserID    int64
	Status    CartStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CartItem 购物车商品领域模型
type CartItem struct {
	ID        int64
	CartID    int64
	UserID    int64
	ProductID int64
	Quantity  int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CartStatus 购物车状态
type CartStatus int

const (
	CartStatusNormal CartStatus = 1 // 正常
	CartStatusEmpty  CartStatus = 2 // 已清空
)

// ValidateQuantity 验证商品数量
func (item *CartItem) ValidateQuantity() bool {
	return item.Quantity > 0 && item.Quantity <= 999
}

// IsValid 验证购物车状态是否有效
func (c *Cart) IsValid() bool {
	return c.Status == CartStatusNormal
}
