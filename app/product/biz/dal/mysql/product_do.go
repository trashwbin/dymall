package mysql

import (
	"time"

	"github.com/trashwbin/dymall/app/product/biz/model"
	"gorm.io/gorm"
)

// ProductDO 商品数据对象
type ProductDO struct {
	ID          uint32         `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"size:100;not null;comment:商品名称"`
	Description string         `gorm:"size:1000;comment:商品描述"`
	Picture     string         `gorm:"size:255;comment:商品图片"`
	Price       float32        `gorm:"not null;comment:商品价格"`
	Status      int            `gorm:"type:tinyint;default:1;comment:商品状态 1:在售 2:下架 3:已删除"`
	CreatedAt   time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt   time.Time      `gorm:"not null;comment:更新时间"`
	DeletedAt   gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// CategoryDO 商品分类数据对象
type CategoryDO struct {
	ID        uint32         `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"size:50;not null;uniqueIndex;comment:分类名称"`
	CreatedAt time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"not null;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// ProductCategoryDO 商品分类关系数据对象
type ProductCategoryDO struct {
	ID         uint32         `gorm:"primaryKey;autoIncrement"`
	ProductID  uint32         `gorm:"index:idx_product_id;not null;comment:商品ID"`
	CategoryID uint32         `gorm:"index:idx_category_id;not null;comment:分类ID"`
	CreatedAt  time.Time      `gorm:"not null;comment:创建时间"`
	DeletedAt  gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// TableName 设置表名
func (ProductDO) TableName() string {
	return "products"
}

// TableName 设置表名
func (CategoryDO) TableName() string {
	return "categories"
}

// TableName 设置表名
func (ProductCategoryDO) TableName() string {
	return "product_categories"
}

// ToModel 转换为领域模型
func (p *ProductDO) ToModel() *model.Product {
	return &model.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Picture:     p.Picture,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (p *ProductDO) FromModel(m *model.Product) {
	p.ID = m.ID
	p.Name = m.Name
	p.Description = m.Description
	p.Picture = m.Picture
	p.Price = m.Price
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
}

// ToModel 转换为领域模型
func (c *CategoryDO) ToModel() *model.Category {
	return &model.Category{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (c *CategoryDO) FromModel(m *model.Category) {
	c.ID = m.ID
	c.Name = m.Name
	c.CreatedAt = m.CreatedAt
	c.UpdatedAt = m.UpdatedAt
}
