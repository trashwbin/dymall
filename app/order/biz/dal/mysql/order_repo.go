package mysql

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/order/biz/model"
	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{db: DB}
}

// CreateOrder 创建订单
func (r *OrderRepo) CreateOrder(order *model.Order) (*model.Order, error) {
	orderDO := &OrderDO{}
	orderDO.FromModel(order)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建订单
		if err := tx.Create(orderDO).Error; err != nil {
			return err
		}

		// 2. 创建订单地址
		addressDO := &AddressDO{}
		addressDO.FromModel(order.Address)
		addressDO.OrderID = orderDO.OrderID // 确保使用正确的订单ID
		if err := tx.Create(addressDO).Error; err != nil {
			return err
		}

		// 3. 创建订单商品
		for _, item := range order.OrderItems {
			itemDO := &OrderItemDO{}
			itemDO.FromModel(item)
			itemDO.OrderID = orderDO.OrderID // 确保使用正确的订单ID
			if err := tx.Create(itemDO).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetOrder(orderDO.OrderID)
}

// GetOrder 获取订单
func (r *OrderRepo) GetOrder(orderID string) (*model.Order, error) {
	var orderDO OrderDO
	if err := r.db.Where("order_id = ?", orderID).First(&orderDO).Error; err != nil {
		return nil, err
	}

	order := orderDO.ToModel()

	// 获取订单地址
	var addressDO AddressDO
	if err := r.db.Where("order_id = ?", orderID).First(&addressDO).Error; err != nil {
		return nil, err
	}
	order.Address = addressDO.ToModel()

	// 获取订单商品
	var itemDOs []OrderItemDO
	if err := r.db.Where("order_id = ?", orderID).Find(&itemDOs).Error; err != nil {
		return nil, err
	}

	items := make([]*model.OrderItem, len(itemDOs))
	for i, itemDO := range itemDOs {
		items[i] = itemDO.ToModel()
	}
	order.OrderItems = items

	return order, nil
}

// UpdateOrder 更新订单
func (r *OrderRepo) UpdateOrder(order *model.Order) error {
	orderDO := &OrderDO{}
	orderDO.FromModel(order)

	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 更新订单
		if err := tx.Model(&OrderDO{}).Where("order_id = ?", order.OrderID).Updates(orderDO).Error; err != nil {
			return err
		}

		// 2. 更新地址
		addressDO := &AddressDO{}
		addressDO.FromModel(order.Address)
		if err := tx.Model(&AddressDO{}).Where("order_id = ?", order.OrderID).Updates(addressDO).Error; err != nil {
			return err
		}

		return nil
	})
}

// ListOrders 获取订单列表
func (r *OrderRepo) ListOrders(userID int64) ([]*model.Order, error) {
	var orderDOs []OrderDO
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(2).Find(&orderDOs).Error; err != nil {
		return nil, err
	}

	orders := make([]*model.Order, 0, len(orderDOs))
	for _, orderDO := range orderDOs {
		order := orderDO.ToModel()

		// 获取订单地址
		var addressDO AddressDO
		if err := r.db.Where("order_id = ?", orderDO.OrderID).First(&addressDO).Error; err != nil {
			klog.Errorf("get order address failed: %v", err)
			continue
		}
		order.Address = addressDO.ToModel()

		// 获取订单商品
		var itemDOs []OrderItemDO
		if err := r.db.Where("order_id = ?", orderDO.OrderID).Find(&itemDOs).Error; err != nil {
			klog.Errorf("get order items failed: %v", err)
			continue
		}

		items := make([]*model.OrderItem, len(itemDOs))
		for i, itemDO := range itemDOs {
			items[i] = itemDO.ToModel()
		}
		order.OrderItems = items

		orders = append(orders, order)
	}

	return orders, nil
}

// Transaction 事务处理
func (r *OrderRepo) Transaction(fn func(txRepo *OrderRepo) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &OrderRepo{db: tx}
		return fn(txRepo)
	})
}
