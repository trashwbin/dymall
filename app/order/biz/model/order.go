package model

import "time"

// OrderStatus 订单状态
type OrderStatus int

const (
	OrderStatusPending  OrderStatus = 1 // 待支付
	OrderStatusPaid     OrderStatus = 2 // 已支付
	OrderStatusCanceled OrderStatus = 3 // 已取消
	OrderStatusExpired  OrderStatus = 4 // 已过期
)

// Order 订单领域模型
type Order struct {
	ID          int64
	OrderID     string  // 订单号
	UserID      int64   // 用户ID
	TotalAmount float64 // 总金额
	Currency    string  // 货币类型
	Status      OrderStatus
	PaymentID   string    // 支付单ID
	Address     *Address  // 收货地址
	Email       string    // 用户邮箱
	ExpireAt    time.Time // 过期时间
	PaidAt      time.Time // 支付时间
	CreatedAt   time.Time
	UpdatedAt   time.Time
	OrderItems  []*OrderItem // 订单商品列表
}

// OrderItem 订单商品领域模型
type OrderItem struct {
	ID        int64
	OrderID   string  // 订单号
	ProductID int64   // 商品ID
	Quantity  int32   // 数量
	Price     float64 // 单价
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Address 收货地址领域模型
type Address struct {
	ID            int64
	OrderID       string
	UserID        int64
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// IsExpired 检查订单是否过期
func (o *Order) IsExpired() bool {
	return time.Now().After(o.ExpireAt)
}

// IsPaid 检查订单是否已支付
func (o *Order) IsPaid() bool {
	return o.Status == OrderStatusPaid
}

// IsCanceled 检查订单是否已取消
func (o *Order) IsCanceled() bool {
	return o.Status == OrderStatusCanceled
}

// CanBePaid 检查订单是否可以支付
func (o *Order) CanBePaid() bool {
	return o.Status == OrderStatusPending && !o.IsExpired()
}

// CanBeCanceled 检查订单是否可以取消
func (o *Order) CanBeCanceled() bool {
	return o.Status == OrderStatusPending && !o.IsExpired()
}

// CalculateTotal 计算订单总金额
func (o *Order) CalculateTotal() float64 {
	var total float64
	for _, item := range o.OrderItems {
		total += float64(item.Quantity) * item.Price
	}
	return total
}
