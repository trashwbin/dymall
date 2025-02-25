package mysql

import (
	"time"

	"github.com/trashwbin/dymall/app/order/biz/model"
	"gorm.io/gorm"
)

// OrderDO 订单数据对象
type OrderDO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	OrderID     string         `gorm:"uniqueIndex;type:varchar(32);not null;comment:订单号"`
	UserID      int64          `gorm:"index:idx_user_id;not null;comment:用户ID"`
	TotalAmount float64        `gorm:"type:decimal(10,2);not null;comment:总金额"`
	Currency    string         `gorm:"type:varchar(10);not null;default:CNY;comment:货币类型"`
	Status      int            `gorm:"type:tinyint;not null;default:1;comment:订单状态"`
	PaymentID   string         `gorm:"type:varchar(32);comment:支付单ID"`
	Email       string         `gorm:"type:varchar(128);not null;comment:用户邮箱"`
	ExpireAt    time.Time      `gorm:"not null;comment:过期时间"`
	PaidAt      *time.Time     `gorm:"comment:支付时间"`
	CreatedAt   time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt   time.Time      `gorm:"not null;comment:更新时间"`
	DeletedAt   gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// OrderItemDO 订单商品数据对象
type OrderItemDO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	OrderID   string         `gorm:"index:idx_order_id;type:varchar(32);not null;comment:订单号"`
	ProductID int64          `gorm:"index:idx_product_id;not null;comment:商品ID"`
	Quantity  int32          `gorm:"not null;comment:数量"`
	Price     float64        `gorm:"type:decimal(10,2);not null;comment:单价"`
	CreatedAt time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"not null;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// AddressDO 收货地址数据对象
type AddressDO struct {
	ID            int64          `gorm:"primaryKey;autoIncrement"`
	OrderID       string         `gorm:"uniqueIndex;type:varchar(32);not null;comment:订单号"`
	UserID        int64          `gorm:"index:idx_user_id;not null;comment:用户ID"`
	StreetAddress string         `gorm:"type:varchar(256);not null;comment:街道地址"`
	City          string         `gorm:"type:varchar(64);not null;comment:城市"`
	State         string         `gorm:"type:varchar(64);not null;comment:州/省"`
	Country       string         `gorm:"type:varchar(64);not null;comment:国家"`
	ZipCode       string         `gorm:"type:varchar(20);not null;comment:邮编"`
	CreatedAt     time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt     time.Time      `gorm:"not null;comment:更新时间"`
	DeletedAt     gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// TableName 设置表名
func (OrderDO) TableName() string {
	return "orders"
}

// TableName 设置表名
func (OrderItemDO) TableName() string {
	return "order_items"
}

// TableName 设置表名
func (AddressDO) TableName() string {
	return "order_addresses"
}

// ToModel 转换为领域模型
func (o *OrderDO) ToModel() *model.Order {
	return &model.Order{
		ID:          o.ID,
		OrderID:     o.OrderID,
		UserID:      o.UserID,
		TotalAmount: o.TotalAmount,
		Currency:    o.Currency,
		Status:      model.OrderStatus(o.Status),
		PaymentID:   o.PaymentID,
		Email:       o.Email,
		ExpireAt:    o.ExpireAt,
		PaidAt:      o.PaidAt,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (o *OrderDO) FromModel(m *model.Order) {
	o.ID = m.ID
	o.OrderID = m.OrderID
	o.UserID = m.UserID
	o.TotalAmount = m.TotalAmount
	o.Currency = m.Currency
	o.Status = int(m.Status)
	o.PaymentID = m.PaymentID
	o.Email = m.Email
	o.ExpireAt = m.ExpireAt
	o.PaidAt = m.PaidAt
	o.CreatedAt = m.CreatedAt
	o.UpdatedAt = m.UpdatedAt
}

// ToModel 转换为领域模型
func (oi *OrderItemDO) ToModel() *model.OrderItem {
	return &model.OrderItem{
		ID:        oi.ID,
		OrderID:   oi.OrderID,
		ProductID: oi.ProductID,
		Quantity:  oi.Quantity,
		Price:     oi.Price,
		CreatedAt: oi.CreatedAt,
		UpdatedAt: oi.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (oi *OrderItemDO) FromModel(m *model.OrderItem) {
	oi.ID = m.ID
	oi.OrderID = m.OrderID
	oi.ProductID = m.ProductID
	oi.Quantity = m.Quantity
	oi.Price = m.Price
	oi.CreatedAt = m.CreatedAt
	oi.UpdatedAt = m.UpdatedAt
}

// ToModel 转换为领域模型
func (a *AddressDO) ToModel() *model.Address {
	return &model.Address{
		ID:            a.ID,
		OrderID:       a.OrderID,
		UserID:        a.UserID,
		StreetAddress: a.StreetAddress,
		City:          a.City,
		State:         a.State,
		Country:       a.Country,
		ZipCode:       a.ZipCode,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (a *AddressDO) FromModel(m *model.Address) {
	a.ID = m.ID
	a.OrderID = m.OrderID
	a.UserID = m.UserID
	a.StreetAddress = m.StreetAddress
	a.City = m.City
	a.State = m.State
	a.Country = m.Country
	a.ZipCode = m.ZipCode
	a.CreatedAt = m.CreatedAt
	a.UpdatedAt = m.UpdatedAt
}
