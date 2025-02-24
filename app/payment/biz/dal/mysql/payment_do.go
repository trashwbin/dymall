package mysql

import (
	"time"

	"github.com/trashwbin/dymall/app/payment/biz/model"
	"gorm.io/gorm"
)

// PaymentDO 支付数据对象
type PaymentDO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	PaymentID string         `gorm:"type:varchar(36);uniqueIndex;not null;comment:支付单号"`
	OrderID   string         `gorm:"type:varchar(32);index;not null;comment:订单号"`
	UserID    int64          `gorm:"index;not null;comment:用户ID"`
	Amount    float64        `gorm:"type:decimal(10,2);not null;comment:支付金额"`
	Currency  string         `gorm:"type:varchar(3);not null;comment:货币类型"`
	Status    int            `gorm:"type:tinyint;not null;default:1;comment:支付状态 1:待支付 2:支付成功 3:支付失败 4:已取消 5:已过期"`
	CreatedAt time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"not null;comment:更新时间"`
	ExpireAt  time.Time      `gorm:"not null;comment:过期时间"`
	PaidAt    *time.Time     `gorm:"comment:支付时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// TableName 设置表名
func (PaymentDO) TableName() string {
	return "payments"
}

// ToModel 转换为领域模型
func (p *PaymentDO) ToModel() *model.Payment {
	var paidAt time.Time
	if p.PaidAt != nil {
		paidAt = *p.PaidAt
	}
	return &model.Payment{
		ID:        p.ID,
		PaymentID: p.PaymentID,
		OrderID:   p.OrderID,
		UserID:    p.UserID,
		Amount:    p.Amount,
		Currency:  p.Currency,
		Status:    model.PayStatus(p.Status),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		ExpireAt:  p.ExpireAt,
		PaidAt:    paidAt,
	}
}

// FromModel 从领域模型转换
func (p *PaymentDO) FromModel(m *model.Payment) {
	p.ID = m.ID
	p.PaymentID = m.PaymentID
	p.OrderID = m.OrderID
	p.UserID = m.UserID
	p.Amount = m.Amount
	p.Currency = m.Currency
	p.Status = int(m.Status)
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
	p.ExpireAt = m.ExpireAt
	if !m.PaidAt.IsZero() {
		paidAt := m.PaidAt
		p.PaidAt = &paidAt
	} else {
		p.PaidAt = nil
	}
}
