package service

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/trashwbin/dymall/app/product/biz/dal"
	"github.com/trashwbin/dymall/app/product/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/product/biz/model"
)

var (
	testProduct *model.Product
)

func TestMain(m *testing.M) {
	// 加载环境变量
	_ = godotenv.Load()

	// 初始化数据库和缓存
	dal.Init()

	// 清理可能存在的测试数据
	if err := cleanupTestData(); err != nil {
		panic(err)
	}

	// 创建测试数据
	if err := setupTestData(); err != nil {
		panic(err)
	}

	// 运行测试
	code := m.Run()

	// 清理测试数据
	if err := cleanupTestData(); err != nil {
		panic(err)
	}

	os.Exit(code)
}

// 创建测试数据
func setupTestData() error {
	repo := mysql.NewProductRepo()

	// 创建测试商品
	product := &model.Product{
		Name:        "测试商品",
		Description: "这是一个测试商品",
		Picture:     "http://example.com/test.jpg",
		Price:       99.9,
		Categories:  []string{"测试分类", "电子产品"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	var err error
	testProduct, err = repo.CreateProduct(product)
	if err != nil {
		return err
	}

	// 创建测试分类
	for _, categoryName := range testProduct.Categories {
		// 使用 GetOrCreateCategory 来处理分类
		category, err := repo.GetOrCreateCategory(categoryName)
		if err != nil {
			return err
		}

		// 添加商品分类关系
		if err := repo.AddProductCategory(testProduct.ID, category.ID); err != nil {
			return err
		}
	}

	return nil
}

// 清理测试数据
func cleanupTestData() error {
	repo := mysql.NewProductRepo()

	// 删除所有测试商品（使用模糊匹配查找所有测试相关的商品）
	products, err := repo.SearchProducts("测试商品")
	if err != nil {
		return err
	}

	for _, p := range products {
		// 删除商品分类关系
		categories, err := repo.GetCategories(p.ID)
		if err != nil {
			continue
		}
		for _, category := range categories {
			_ = repo.RemoveCategory(p.ID, category.ID)
		}
		// 删除商品
		_ = repo.DeleteProduct(p.ID)
	}

	return nil
}
