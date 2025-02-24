package model

import "time"

// Product 商品领域模型
type Product struct {
	ID          uint32
	Name        string
	Description string
	Picture     string
	Price       float32
	Categories  []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Category 商品分类领域模型
type Category struct {
	ID        uint32
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ProductStatus 商品状态
type ProductStatus int

const (
	ProductStatusOnSale  ProductStatus = 1 // 在售
	ProductStatusOffSale ProductStatus = 2 // 下架
	ProductStatusDeleted ProductStatus = 3 // 已删除
)

// ValidatePrice 验证商品价格
func (p *Product) ValidatePrice() bool {
	return p.Price >= 0
}

// ValidateName 验证商品名称
func (p *Product) ValidateName() bool {
	return len(p.Name) > 0 && len(p.Name) <= 100
}

// ValidateDescription 验证商品描述
func (p *Product) ValidateDescription() bool {
	return len(p.Description) <= 1000
}

// IsValid 验证商品信息是否有效
func (p *Product) IsValid() bool {
	return p.ValidateName() && p.ValidatePrice() && p.ValidateDescription()
}
