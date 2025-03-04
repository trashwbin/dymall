package service

import (
	"context"
	"os"
	"testing"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/joho/godotenv"
	"github.com/trashwbin/dymall/app/cart/biz/dal"
	"github.com/trashwbin/dymall/app/cart/infra/rpc"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

// MockProductClient 模拟商品服务客户端
type MockProductClient struct{}

func (m *MockProductClient) GetProduct(ctx context.Context, req *product.GetProductReq, opts ...callopt.Option) (r *product.GetProductResp, err error) {
	// 模拟商品服务响应
	if req.Id == 1001 {
		return &product.GetProductResp{
			Product: &product.Product{
				Id:          1001,
				Name:        "测试商品",
				Description: "这是一个测试商品",
				Price:       99.9,
			},
		}, nil
	}
	if req.Id == 1002 {
		return &product.GetProductResp{
			Product: &product.Product{
				Id:          1002,
				Name:        "测试商品2",
				Description: "这是一个测试商品2",
				Price:       199.9,
			},
		}, nil
	}
	return &product.GetProductResp{}, nil
}

func (m *MockProductClient) ListProducts(ctx context.Context, req *product.ListProductsReq, opts ...callopt.Option) (r *product.ListProductsResp, err error) {
	return nil, nil
}

func (m *MockProductClient) SearchProducts(ctx context.Context, req *product.SearchProductsReq, opts ...callopt.Option) (r *product.SearchProductsResp, err error) {
	return nil, nil
}

func (m *MockProductClient) CreateProduct(ctx context.Context, req *product.CreateProductReq, opts ...callopt.Option) (r *product.CreateProductResp, err error) {
	return nil, nil
}

func (m *MockProductClient) UpdateProduct(ctx context.Context, req *product.UpdateProductReq, opts ...callopt.Option) (r *product.UpdateProductResp, err error) {
	return nil, nil
}

func (m *MockProductClient) DeleteProduct(ctx context.Context, req *product.DeleteProductReq, opts ...callopt.Option) (r *product.DeleteProductResp, err error) {
	return nil, nil
}

func (m *MockProductClient) BatchGetProducts(ctx context.Context, req *product.BatchGetProductsReq, opts ...callopt.Option) (r *product.BatchGetProductsResp, err error) {
	return &product.BatchGetProductsResp{}, nil
}

func TestMain(m *testing.M) {
	// 替换为mock客户端
	rpc.ProductClient = &MockProductClient{}

	// 加载环境变量
	_ = godotenv.Load()

	// 初始化数据库
	dal.Init()

	// 运行测试
	code := m.Run()
	os.Exit(code)
}
