package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/product/biz/model"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

// mockProduct 创建模拟商品数据
func mockProduct(id uint32) *model.Product {
	return &model.Product{
		ID:          id,
		Name:        fmt.Sprintf("商品%d", id),
		Description: fmt.Sprintf("这是商品%d的描述", id),
		Picture:     fmt.Sprintf("http://example.com/product/%d.jpg", id),
		Price:       float32(id) * 10.0,
		Categories:  []string{"测试分类1", "测试分类2"},
	}
}

func TestBatchGetProducts_Run(t *testing.T) {
	ctx := context.Background()

	// 确保测试商品已创建
	assert.NotNil(t, testProduct, "测试商品应该已经创建")

	tests := []struct {
		name     string
		req      *product.BatchGetProductsReq
		wantErr  bool
		validate func(*testing.T, *product.BatchGetProductsResp)
	}{
		{
			name: "空请求测试",
			req:  &product.BatchGetProductsReq{},
			validate: func(t *testing.T, resp *product.BatchGetProductsResp) {
				assert.NotNil(t, resp)
				assert.Empty(t, resp.Products)
				assert.Empty(t, resp.MissingIds)
			},
		},
		{
			name: "正常查询测试",
			req: &product.BatchGetProductsReq{
				Ids: []uint32{testProduct.ID},
			},
			validate: func(t *testing.T, resp *product.BatchGetProductsResp) {
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Products)
				assert.Len(t, resp.Products, 1)

				// 验证返回的商品数据
				p := resp.Products[0]
				assert.Equal(t, testProduct.ID, p.Id)
				assert.Equal(t, testProduct.Name, p.Name)
				assert.Equal(t, testProduct.Description, p.Description)
				assert.Equal(t, testProduct.Picture, p.Picture)
				assert.Equal(t, testProduct.Price, p.Price)
				assert.Equal(t, testProduct.Categories, p.Categories)
			},
		},
		{
			name: "部分商品不存在测试",
			req: &product.BatchGetProductsReq{
				Ids: []uint32{testProduct.ID, 999, 1000},
			},
			validate: func(t *testing.T, resp *product.BatchGetProductsResp) {
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Products)
				assert.NotEmpty(t, resp.MissingIds)
				assert.Contains(t, resp.MissingIds, uint32(999))
				assert.Contains(t, resp.MissingIds, uint32(1000))
				assert.Len(t, resp.Products, 1)
				assert.Equal(t, testProduct.ID, resp.Products[0].Id)
			},
		},
		{
			name: "重复ID测试",
			req: &product.BatchGetProductsReq{
				Ids: []uint32{testProduct.ID, testProduct.ID, testProduct.ID},
			},
			validate: func(t *testing.T, resp *product.BatchGetProductsResp) {
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Products)
				assert.Len(t, resp.Products, 1)
				assert.Equal(t, testProduct.ID, resp.Products[0].Id)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewBatchGetProductsService(ctx)
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			tt.validate(t, resp)
		})
	}
}

// TestConvertToProtoProducts 测试模型转换函数
func TestConvertToProtoProducts(t *testing.T) {
	products := []*model.Product{
		testProduct,
		mockProduct(999), // 使用mock数据作为第二个测试数据
	}

	protoProducts := convertToProtoProducts(products)
	assert.Equal(t, len(products), len(protoProducts))

	// 验证第一个商品（真实测试数据）
	assert.Equal(t, testProduct.ID, protoProducts[0].Id)
	assert.Equal(t, testProduct.Name, protoProducts[0].Name)
	assert.Equal(t, testProduct.Description, protoProducts[0].Description)
	assert.Equal(t, testProduct.Picture, protoProducts[0].Picture)
	assert.Equal(t, testProduct.Price, protoProducts[0].Price)
	assert.Equal(t, testProduct.Categories, protoProducts[0].Categories)

	// 验证第二个商品（mock数据）
	mockP := mockProduct(999)
	assert.Equal(t, mockP.ID, protoProducts[1].Id)
	assert.Equal(t, mockP.Name, protoProducts[1].Name)
	assert.Equal(t, mockP.Description, protoProducts[1].Description)
	assert.Equal(t, mockP.Picture, protoProducts[1].Picture)
	assert.Equal(t, mockP.Price, protoProducts[1].Price)
	assert.Equal(t, mockP.Categories, protoProducts[1].Categories)
}
