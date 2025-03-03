package mysql

import (
	"github.com/trashwbin/dymall/app/payment/biz/model"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo() *PaymentRepo {
	return &PaymentRepo{db: DB}
}

// CreatePayment 创建支付单
func (r *PaymentRepo) CreatePayment(payment *model.Payment) (*model.Payment, error) {
	paymentDO := &PaymentDO{}
	paymentDO.FromModel(payment)
	if err := r.db.Create(paymentDO).Error; err != nil {
		return nil, err
	}
	return paymentDO.ToModel(), nil
}

// GetPaymentByID 根据支付单号获取支付单
func (r *PaymentRepo) GetPaymentByID(paymentID string) (*model.Payment, error) {
	var paymentDO PaymentDO
	err := r.db.Where("payment_id = ?", paymentID).First(&paymentDO).Error
	if err != nil {
		return nil, err
	}
	return paymentDO.ToModel(), nil
}

// GetPaymentByOrderID 根据订单号获取支付单
func (r *PaymentRepo) GetPaymentByOrderID(orderID string) (*model.Payment, error) {
	var paymentDO PaymentDO
	err := r.db.Where("order_id = ?", orderID).First(&paymentDO).Error
	if err != nil {
		return nil, err
	}
	return paymentDO.ToModel(), nil
}

// UpdatePayment 更新支付单
func (r *PaymentRepo) UpdatePayment(payment *model.Payment) error {
	paymentDO := &PaymentDO{}
	paymentDO.FromModel(payment)
	return r.db.Save(paymentDO).Error
}

// UpdatePaymentStatus 更新支付单状态
func (r *PaymentRepo) UpdatePaymentStatus(paymentID string, status model.PayStatus) error {
	return r.db.Model(&PaymentDO{}).
		Where("payment_id = ?", paymentID).
		Updates(map[string]interface{}{
			"status": int(status),
		}).Error
}

// ListPaymentsByUserID 获取用户的支付单列表
func (r *PaymentRepo) ListPaymentsByUserID(userID int64) ([]*model.Payment, error) {
	var paymentDOs []PaymentDO
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&paymentDOs).Error
	if err != nil {
		return nil, err
	}

	payments := make([]*model.Payment, len(paymentDOs))
	for i, paymentDO := range paymentDOs {
		payments[i] = paymentDO.ToModel()
	}
	return payments, nil
}

// Transaction 事务处理
func (r *PaymentRepo) Transaction(fn func(txRepo *PaymentRepo) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &PaymentRepo{db: tx}
		return fn(txRepo)
	})
}

// DeletePaymentForTest 删除支付单（仅用于测试）
func (r *PaymentRepo) DeletePaymentForTest(paymentID string) error {
	return r.db.Unscoped().Where("payment_id = ?", paymentID).Delete(&PaymentDO{}).Error
}
