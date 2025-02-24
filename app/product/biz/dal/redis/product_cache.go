package redis

import (
	"fmt"
	"time"

	"github.com/trashwbin/dymall/app/product/biz/model"
)

const (
	// 商品缓存key前缀
	ProductKeyPrefix = "product:"
	// 商品列表缓存key前缀
	ProductListKeyPrefix = "product:list:"
	// 商品分类缓存key前缀
	CategoryKeyPrefix = "category:"
	// 商品缓存过期时间
	ProductExpiration = 24 * time.Hour
)

// ProductCache Redis商品缓存数据对象
type ProductCache struct {
	ID          uint32    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Picture     string    `json:"picture"`
	Price       float32   `json:"price"`
	Categories  []string  `json:"categories"`
	Status      int       `json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryCache Redis商品分类缓存数据对象
type CategoryCache struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToModel 转换为领域模型
func (p *ProductCache) ToModel() *model.Product {
	return &model.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Picture:     p.Picture,
		Price:       p.Price,
		Categories:  p.Categories,
		UpdatedAt:   p.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (p *ProductCache) FromModel(m *model.Product) {
	p.ID = m.ID
	p.Name = m.Name
	p.Description = m.Description
	p.Picture = m.Picture
	p.Price = m.Price
	p.Categories = m.Categories
	p.UpdatedAt = m.UpdatedAt
}

// GetKey 获取商品缓存key
func (p *ProductCache) GetKey() string {
	return fmt.Sprintf("%s%d", ProductKeyPrefix, p.ID)
}

// ToModel 转换为领域模型
func (c *CategoryCache) ToModel() *model.Category {
	return &model.Category{
		ID:        c.ID,
		Name:      c.Name,
		UpdatedAt: c.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (c *CategoryCache) FromModel(m *model.Category) {
	c.ID = m.ID
	c.Name = m.Name
	c.UpdatedAt = m.UpdatedAt
}

// GetKey 获取分类缓存key
func (c *CategoryCache) GetKey() string {
	return fmt.Sprintf("%s%d", CategoryKeyPrefix, c.ID)
}

// GetProductListKey 获取商品列表缓存key
func GetProductListKey(page int32, pageSize int64, categoryName string) string {
	if categoryName != "" {
		return fmt.Sprintf("%s%d:%d:%s", ProductListKeyPrefix, page, pageSize, categoryName)
	}
	return fmt.Sprintf("%s%d:%d", ProductListKeyPrefix, page, pageSize)
}
