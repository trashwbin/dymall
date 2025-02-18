package redis

import (
	"fmt"
	"time"

	"github.com/trashwbin/dymall/app/cart/biz/model"
)

const (
	// 购物车缓存key前缀
	CartKeyPrefix = "cart:user:"
	// 购物车商品缓存key前缀
	CartItemKeyPrefix = "cart:item:"
	// 购物车缓存过期时间
	CartExpiration = 72 * time.Hour
)

// CartCache Redis购物车缓存数据对象
type CartCache struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Status    int       `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CartItemCache Redis购物车商品缓存数据对象
type CartItemCache struct {
	ID        int64     `json:"id"`
	CartID    int64     `json:"cart_id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int32     `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToModel 转换为领域模型
func (c *CartCache) ToModel() *model.Cart {
	return &model.Cart{
		ID:        c.ID,
		UserID:    c.UserID,
		Status:    model.CartStatus(c.Status),
		UpdatedAt: c.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (c *CartCache) FromModel(m *model.Cart) {
	c.ID = m.ID
	c.UserID = m.UserID
	c.Status = int(m.Status)
	c.UpdatedAt = m.UpdatedAt
}

// GetKey 获取购物车缓存key
func (c *CartCache) GetKey() string {
	return fmt.Sprintf("%s%d", CartKeyPrefix, c.UserID)
}

// ToModel 转换为领域模型
func (ci *CartItemCache) ToModel() (*model.CartItem, error) {
	return &model.CartItem{
		ID:        ci.ID,
		CartID:    ci.CartID,
		UserID:    ci.UserID,
		ProductID: ci.ProductID,
		Quantity:  ci.Quantity,
		UpdatedAt: ci.UpdatedAt,
	}, nil
}

// FromModel 从领域模型转换
func (ci *CartItemCache) FromModel(m *model.CartItem) error {
	ci.ID = m.ID
	ci.CartID = m.CartID
	ci.UserID = m.UserID
	ci.ProductID = m.ProductID
	ci.Quantity = m.Quantity
	ci.UpdatedAt = m.UpdatedAt
	return nil
}

// GetKey 获取购物车商品缓存key
func (ci *CartItemCache) GetKey() string {
	return fmt.Sprintf("%s%d:%d", CartItemKeyPrefix, ci.CartID, ci.ProductID)
}

// GetCartItemPattern 获取购物车所有商品的key模式
func GetCartItemPattern(cartID int64) string {
	return fmt.Sprintf("%s%d:*", CartItemKeyPrefix, cartID)
}
