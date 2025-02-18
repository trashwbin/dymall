package mysql

import (
	"time"

	"github.com/trashwbin/dymall/app/cart/biz/model"
	"gorm.io/gorm"
)

// CartDO 购物车数据对象
type CartDO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"index:idx_user_id;not null;comment:用户ID"`
	Status    int            `gorm:"type:tinyint;default:1;comment:购物车状态 1:正常 2:已清空"`
	CreatedAt time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"not null;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// CartItemDO 购物车商品数据对象
type CartItemDO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	CartID    int64          `gorm:"index:idx_cart_id;not null;comment:购物车ID"`
	UserID    int64          `gorm:"index:idx_user_id;not null;comment:用户ID"`
	ProductID int64          `gorm:"index:idx_product_id;not null;comment:商品ID"`
	Quantity  int32          `gorm:"not null;default:1;comment:商品数量"`
	CreatedAt time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"not null;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// TableName 设置表名
func (CartDO) TableName() string {
	return "carts"
}

// TableName 设置表名
func (CartItemDO) TableName() string {
	return "cart_items"
}

// ToModel 转换为领域模型
func (c *CartDO) ToModel() *model.Cart {
	return &model.Cart{
		ID:        c.ID,
		UserID:    c.UserID,
		Status:    model.CartStatus(c.Status),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (c *CartDO) FromModel(m *model.Cart) {
	c.ID = m.ID
	c.UserID = m.UserID
	c.Status = int(m.Status)
	c.CreatedAt = m.CreatedAt
	c.UpdatedAt = m.UpdatedAt
}

// ToModel 转换为领域模型
func (ci *CartItemDO) ToModel() *model.CartItem {
	return &model.CartItem{
		ID:        ci.ID,
		CartID:    ci.CartID,
		UserID:    ci.UserID,
		ProductID: ci.ProductID,
		Quantity:  ci.Quantity,
		CreatedAt: ci.CreatedAt,
		UpdatedAt: ci.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (ci *CartItemDO) FromModel(m *model.CartItem) {
	ci.ID = m.ID
	ci.CartID = m.CartID
	ci.UserID = m.UserID
	ci.ProductID = m.ProductID
	ci.Quantity = m.Quantity
	ci.CreatedAt = m.CreatedAt
	ci.UpdatedAt = m.UpdatedAt
}
