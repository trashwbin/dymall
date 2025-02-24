package mysql

import (
	"github.com/trashwbin/dymall/app/checkout/biz/model"
	"gorm.io/gorm"
)

type CheckoutRepo struct {
	db *gorm.DB
}

func NewCheckoutRepo() *CheckoutRepo {
	return &CheckoutRepo{db: DB}
}

// GetCheckout 获取结算单
func (r *CheckoutRepo) GetCheckout(id string) (*model.Checkout, error) {
	var checkoutDO CheckoutDO
	err := r.db.First(&checkoutDO, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// 获取结算单商品
	var itemDOs []CheckoutItemDO
	err = r.db.Where("checkout_id = ?", id).Find(&itemDOs).Error
	if err != nil {
		return nil, err
	}

	checkout := checkoutDO.ToModel()
	checkout.Items = make([]*model.CheckoutItem, len(itemDOs))
	for i, itemDO := range itemDOs {
		checkout.Items[i] = itemDO.ToModel()
	}

	return checkout, nil
}

// CreateCheckout 创建结算单
func (r *CheckoutRepo) CreateCheckout(checkout *model.Checkout) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建结算单
		checkoutDO := &CheckoutDO{}
		checkoutDO.FromModel(checkout)
		if err := tx.Create(checkoutDO).Error; err != nil {
			return err
		}

		// 创建结算单商品
		for _, item := range checkout.Items {
			itemDO := &CheckoutItemDO{}
			itemDO.FromModel(item)
			itemDO.CheckoutID = checkout.ID
			if err := tx.Create(itemDO).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateCheckout 更新结算单
func (r *CheckoutRepo) UpdateCheckout(checkout *model.Checkout) error {
	checkoutDO := &CheckoutDO{}
	checkoutDO.FromModel(checkout)
	return r.db.Save(checkoutDO).Error
}

// DeleteCheckout 删除结算单
func (r *CheckoutRepo) DeleteCheckout(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除结算单商品
		if err := tx.Where("checkout_id = ?", id).Delete(&CheckoutItemDO{}).Error; err != nil {
			return err
		}

		// 删除结算单
		return tx.Delete(&CheckoutDO{}, "id = ?", id).Error
	})
}

// GetCheckoutByUserID 根据用户ID获取结算单
func (r *CheckoutRepo) GetCheckoutByUserID(userID int64) ([]*model.Checkout, error) {
	var checkoutDOs []CheckoutDO
	err := r.db.Where("user_id = ?", userID).Find(&checkoutDOs).Error
	if err != nil {
		return nil, err
	}

	checkouts := make([]*model.Checkout, len(checkoutDOs))
	for i, checkoutDO := range checkoutDOs {
		checkout := checkoutDO.ToModel()

		// 获取结算单商品
		var itemDOs []CheckoutItemDO
		err = r.db.Where("checkout_id = ?", checkout.ID).Find(&itemDOs).Error
		if err != nil {
			return nil, err
		}

		checkout.Items = make([]*model.CheckoutItem, len(itemDOs))
		for j, itemDO := range itemDOs {
			checkout.Items[j] = itemDO.ToModel()
		}

		checkouts[i] = checkout
	}

	return checkouts, nil
}

// Transaction 事务处理
func (r *CheckoutRepo) Transaction(fn func(txRepo *CheckoutRepo) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &CheckoutRepo{db: tx}
		return fn(txRepo)
	})
}
