package mysql

import (
	"github.com/trashwbin/dymall/app/cart/biz/model"
	"gorm.io/gorm"
)

type CartRepo struct {
	db *gorm.DB
}

func NewCartRepo() *CartRepo {
	return &CartRepo{db: DB}
}

// GetCartByUserID 根据用户ID获取购物车
func (r *CartRepo) GetCartByUserID(userID int64, status model.CartStatus) (*model.Cart, error) {
	var cartDO CartDO
	err := r.db.Where("user_id = ? AND status = ?", userID, status).First(&cartDO).Error
	if err != nil {
		return nil, err
	}
	return cartDO.ToModel(), nil
}

// CreateCart 创建购物车
func (r *CartRepo) CreateCart(cart *model.Cart) (*model.Cart, error) {
	cartDO := &CartDO{}
	cartDO.FromModel(cart)
	if err := r.db.Create(cartDO).Error; err != nil {
		return nil, err
	}
	return cartDO.ToModel(), nil
}

// GetCartItem 获取购物车商品
func (r *CartRepo) GetCartItem(cartID int64, productID int64) (*model.CartItem, error) {
	var itemDO CartItemDO
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&itemDO).Error
	if err != nil {
		return nil, err
	}
	return itemDO.ToModel(), nil
}

// UpdateCartItem 更新购物车商品
func (r *CartRepo) UpdateCartItem(item *model.CartItem) error {
	itemDO := &CartItemDO{}
	itemDO.FromModel(item)
	return r.db.Save(itemDO).Error
}

// CreateCartItem 创建购物车商品
func (r *CartRepo) CreateCartItem(item *model.CartItem) error {
	itemDO := &CartItemDO{}
	itemDO.FromModel(item)
	return r.db.Create(itemDO).Error
}

// EmptyCart 清空购物车
func (r *CartRepo) EmptyCart(userID int64) error {
	// 使用软删除清空购物车商品
	if err := r.db.Where("user_id = ?", userID).Delete(&CartItemDO{}).Error; err != nil {
		return err
	}

	// 更新购物车状态为已清空
	return r.db.Model(&CartDO{}).
		Where("user_id = ? AND status = ?", userID, model.CartStatusNormal).
		Update("status", model.CartStatusEmpty).Error
}

// GetCartItems 获取用户购物车所有商品
func (r *CartRepo) GetCartItems(userID int64) ([]*model.CartItem, error) {
	var itemDOs []CartItemDO
	err := r.db.Where("user_id = ?", userID).Find(&itemDOs).Error
	if err != nil {
		return nil, err
	}

	items := make([]*model.CartItem, len(itemDOs))
	for i, itemDO := range itemDOs {
		items[i] = itemDO.ToModel()
	}
	return items, nil
}

// Transaction 事务处理
func (r *CartRepo) Transaction(fn func(txRepo *CartRepo) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &CartRepo{db: tx}
		return fn(txRepo)
	})
}
