package mysql

import (
	"time"

	"github.com/trashwbin/dymall/app/checkout/biz/model"
	"gorm.io/gorm"
)

// CheckoutDO 结算单数据对象
type CheckoutDO struct {
	ID          string         `gorm:"primaryKey;type:varchar(32)"`
	UserID      int64          `gorm:"index:idx_user_id;not null;comment:用户ID"`
	TotalAmount float64        `gorm:"type:decimal(10,2);not null;comment:总金额"`
	Currency    string         `gorm:"type:varchar(3);not null;default:CNY;comment:货币类型"`
	Status      int            `gorm:"type:tinyint;default:1;comment:结算单状态 1:待提交 2:已提交 3:已过期"`
	CreatedAt   time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt   time.Time      `gorm:"not null;comment:更新时间"`
	ExpireAt    time.Time      `gorm:"not null;comment:过期时间"`
	DeletedAt   gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// CheckoutItemDO 结算单商品数据对象
type CheckoutItemDO struct {
	ID         int64          `gorm:"primaryKey;autoIncrement"`
	CheckoutID string         `gorm:"index:idx_checkout_id;type:varchar(32);not null;comment:结算单ID"`
	ProductID  int64          `gorm:"index:idx_product_id;not null;comment:商品ID"`
	Quantity   int32          `gorm:"not null;comment:商品数量"`
	Price      float64        `gorm:"type:decimal(10,2);not null;comment:商品单价"`
	Currency   string         `gorm:"type:varchar(3);not null;default:CNY;comment:货币类型"`
	CreatedAt  time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt  time.Time      `gorm:"not null;comment:更新时间"`
	DeletedAt  gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// TableName 设置表名
func (CheckoutDO) TableName() string {
	return "checkouts"
}

// TableName 设置表名
func (CheckoutItemDO) TableName() string {
	return "checkout_items"
}

// ToModel 转换为领域模型
func (c *CheckoutDO) ToModel() *model.Checkout {
	return &model.Checkout{
		ID:          c.ID,
		UserID:      c.UserID,
		TotalAmount: c.TotalAmount,
		Currency:    c.Currency,
		Status:      model.CheckoutStatus(c.Status),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		ExpireAt:    c.ExpireAt,
	}
}

// FromModel 从领域模型转换
func (c *CheckoutDO) FromModel(m *model.Checkout) {
	c.ID = m.ID
	c.UserID = m.UserID
	c.TotalAmount = m.TotalAmount
	c.Currency = m.Currency
	c.Status = int(m.Status)
	c.CreatedAt = m.CreatedAt
	c.UpdatedAt = m.UpdatedAt
	c.ExpireAt = m.ExpireAt
}

// ToModel 转换为领域模型
func (ci *CheckoutItemDO) ToModel() *model.CheckoutItem {
	return &model.CheckoutItem{
		ID:        ci.ID,
		ProductID: ci.ProductID,
		Quantity:  ci.Quantity,
		Price:     ci.Price,
		Currency:  ci.Currency,
	}
}

// FromModel 从领域模型转换
func (ci *CheckoutItemDO) FromModel(m *model.CheckoutItem) {
	ci.ID = m.ID
	ci.ProductID = m.ProductID
	ci.Quantity = m.Quantity
	ci.Price = m.Price
	ci.Currency = m.Currency
}
