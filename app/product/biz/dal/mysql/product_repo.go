package mysql

import (
	"github.com/trashwbin/dymall/app/product/biz/model"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{db: DB}
}

// CreateProduct 创建商品
func (r *ProductRepo) CreateProduct(product *model.Product) (*model.Product, error) {
	productDO := &ProductDO{}
	productDO.FromModel(product)
	if err := r.db.Create(productDO).Error; err != nil {
		return nil, err
	}
	return productDO.ToModel(), nil
}

// UpdateProduct 更新商品
func (r *ProductRepo) UpdateProduct(product *model.Product) error {
	productDO := &ProductDO{}
	productDO.FromModel(product)
	return r.db.Save(productDO).Error
}

// DeleteProduct 删除商品
func (r *ProductRepo) DeleteProduct(id uint32) error {
	return r.db.Model(&ProductDO{}).Where("id = ?", id).Update("status", model.ProductStatusDeleted).Error
}

// GetProduct 获取商品
func (r *ProductRepo) GetProduct(id uint32) (*model.Product, error) {
	var productDO ProductDO
	if err := r.db.First(&productDO, id).Error; err != nil {
		return nil, err
	}
	return productDO.ToModel(), nil
}

// ListProducts 获取商品列表
func (r *ProductRepo) ListProducts(page int32, pageSize int64, categoryName string) ([]*model.Product, error) {
	var productDOs []ProductDO
	query := r.db.Model(&ProductDO{})

	if categoryName != "" {
		query = query.Joins("JOIN product_categories pc ON products.id = pc.product_id").
			Joins("JOIN categories c ON pc.category_id = c.id").
			Where("c.name = ?", categoryName)
	}

	err := query.Offset(int((page - 1) * int32(pageSize))).
		Limit(int(pageSize)).
		Find(&productDOs).Error
	if err != nil {
		return nil, err
	}

	products := make([]*model.Product, len(productDOs))
	for i, productDO := range productDOs {
		products[i] = productDO.ToModel()
	}
	return products, nil
}

// SearchProducts 搜索商品
func (r *ProductRepo) SearchProducts(query string) ([]*model.Product, error) {
	var productDOs []ProductDO
	err := r.db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&productDOs).Error
	if err != nil {
		return nil, err
	}

	products := make([]*model.Product, len(productDOs))
	for i, productDO := range productDOs {
		products[i] = productDO.ToModel()
	}
	return products, nil
}

// AddCategory 添加商品分类
func (r *ProductRepo) AddCategory(productID uint32, categoryID uint32) error {
	return r.db.Create(&ProductCategoryDO{
		ProductID:  productID,
		CategoryID: categoryID,
	}).Error
}

// RemoveCategory 移除商品分类
func (r *ProductRepo) RemoveCategory(productID uint32, categoryID uint32) error {
	return r.db.Where("product_id = ? AND category_id = ?", productID, categoryID).
		Delete(&ProductCategoryDO{}).Error
}

// GetCategories 获取商品分类列表
func (r *ProductRepo) GetCategories(productID uint32) ([]*model.Category, error) {
	var categories []CategoryDO
	err := r.db.Model(&CategoryDO{}).
		Joins("JOIN product_categories pc ON categories.id = pc.category_id").
		Where("pc.product_id = ?", productID).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	result := make([]*model.Category, len(categories))
	for i, category := range categories {
		result[i] = category.ToModel()
	}
	return result, nil
}

// Transaction 事务处理
func (r *ProductRepo) Transaction(fn func(txRepo *ProductRepo) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &ProductRepo{db: tx}
		return fn(txRepo)
	})
}
