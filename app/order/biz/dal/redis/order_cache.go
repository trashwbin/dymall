package redis

import (
	"fmt"
	"time"

	"github.com/trashwbin/dymall/app/order/biz/model"
)

const (
	// 订单缓存key前缀
	OrderKeyPrefix = "order:"
	// 用户订单列表缓存key前缀
	UserOrdersKeyPrefix = "user:orders:"
	// 订单缓存过期时间
	OrderExpiration = 24 * time.Hour
)

// OrderCache Redis订单缓存数据对象
type OrderCache struct {
	ID          int64             `json:"id"`
	OrderID     string            `json:"order_id"`
	UserID      int64             `json:"user_id"`
	TotalAmount float64           `json:"total_amount"`
	Currency    string            `json:"currency"`
	Status      int               `json:"status"`
	PaymentID   string            `json:"payment_id"`
	Email       string            `json:"email"`
	ExpireAt    time.Time         `json:"expire_at"`
	PaidAt      *time.Time        `json:"paid_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Address     *AddressCache     `json:"address"`
	OrderItems  []*OrderItemCache `json:"order_items"`
}

// OrderItemCache Redis订单商品缓存数据对象
type OrderItemCache struct {
	ID        int64     `json:"id"`
	OrderID   string    `json:"order_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int32     `json:"quantity"`
	Price     float64   `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddressCache Redis订单地址缓存数据对象
type AddressCache struct {
	ID            int64     `json:"id"`
	OrderID       string    `json:"order_id"`
	UserID        int64     `json:"user_id"`
	StreetAddress string    `json:"street_address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	ZipCode       string    `json:"zip_code"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ToModel 转换为领域模型
func (c *OrderCache) ToModel() *model.Order {
	order := &model.Order{
		ID:          c.ID,
		OrderID:     c.OrderID,
		UserID:      c.UserID,
		TotalAmount: c.TotalAmount,
		Currency:    c.Currency,
		Status:      model.OrderStatus(c.Status),
		PaymentID:   c.PaymentID,
		Email:       c.Email,
		ExpireAt:    c.ExpireAt,
		PaidAt:      c.PaidAt,
		UpdatedAt:   c.UpdatedAt,
	}

	if c.Address != nil {
		order.Address = c.Address.ToModel()
	}

	if len(c.OrderItems) > 0 {
		items := make([]*model.OrderItem, len(c.OrderItems))
		for i, item := range c.OrderItems {
			items[i] = item.ToModel()
		}
		order.OrderItems = items
	}

	return order
}

// FromModel 从领域模型转换
func (c *OrderCache) FromModel(m *model.Order) {
	c.ID = m.ID
	c.OrderID = m.OrderID
	c.UserID = m.UserID
	c.TotalAmount = m.TotalAmount
	c.Currency = m.Currency
	c.Status = int(m.Status)
	c.PaymentID = m.PaymentID
	c.Email = m.Email
	c.ExpireAt = m.ExpireAt
	c.PaidAt = m.PaidAt
	c.UpdatedAt = m.UpdatedAt

	if m.Address != nil {
		address := &AddressCache{}
		address.FromModel(m.Address)
		c.Address = address
	}

	if len(m.OrderItems) > 0 {
		items := make([]*OrderItemCache, len(m.OrderItems))
		for i, item := range m.OrderItems {
			itemCache := &OrderItemCache{}
			itemCache.FromModel(item)
			items[i] = itemCache
		}
		c.OrderItems = items
	}
}

// GetKey 获取订单缓存key
func (c *OrderCache) GetKey() string {
	return fmt.Sprintf("%s%s", OrderKeyPrefix, c.OrderID)
}

// ToModel 转换为领域模型
func (i *OrderItemCache) ToModel() *model.OrderItem {
	return &model.OrderItem{
		ID:        i.ID,
		OrderID:   i.OrderID,
		ProductID: i.ProductID,
		Quantity:  i.Quantity,
		Price:     i.Price,
		UpdatedAt: i.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (i *OrderItemCache) FromModel(m *model.OrderItem) {
	i.ID = m.ID
	i.OrderID = m.OrderID
	i.ProductID = m.ProductID
	i.Quantity = m.Quantity
	i.Price = m.Price
	i.UpdatedAt = m.UpdatedAt
}

// ToModel 转换为领域模型
func (a *AddressCache) ToModel() *model.Address {
	return &model.Address{
		ID:            a.ID,
		OrderID:       a.OrderID,
		UserID:        a.UserID,
		StreetAddress: a.StreetAddress,
		City:          a.City,
		State:         a.State,
		Country:       a.Country,
		ZipCode:       a.ZipCode,
		UpdatedAt:     a.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (a *AddressCache) FromModel(m *model.Address) {
	a.ID = m.ID
	a.OrderID = m.OrderID
	a.UserID = m.UserID
	a.StreetAddress = m.StreetAddress
	a.City = m.City
	a.State = m.State
	a.Country = m.Country
	a.ZipCode = m.ZipCode
	a.UpdatedAt = m.UpdatedAt
}

// GetUserOrdersKey 获取用户订单列表缓存key
func GetUserOrdersKey(userID int64) string {
	return fmt.Sprintf("%s%d", UserOrdersKeyPrefix, userID)
}
