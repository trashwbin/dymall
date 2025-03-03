package service

import (
	"context"
	"os"
	"testing"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/joho/godotenv"
	"github.com/trashwbin/dymall/app/checkout/infra/rpc"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

// MockCartClient 模拟购物车服务客户端
type MockCartClient struct{}

func (m *MockCartClient) GetCart(ctx context.Context, req *cart.GetCartReq, opts ...callopt.Option) (r *cart.GetCartResp, err error) {
	if req.UserId == 1001 {
		return &cart.GetCartResp{
			Cart: &cart.Cart{
				Items: []*cart.CartItem{
					{
						ProductId: 1001,
						Quantity:  2,
					},
					{
						ProductId: 1002,
						Quantity:  1,
					},
				},
			},
		}, nil
	}
	return &cart.GetCartResp{
		Cart: &cart.Cart{
			Items: []*cart.CartItem{},
		},
	}, nil
}

func (m *MockCartClient) AddItem(ctx context.Context, req *cart.AddItemReq, opts ...callopt.Option) (r *cart.AddItemResp, err error) {
	return &cart.AddItemResp{}, nil
}

func (m *MockCartClient) EmptyCart(ctx context.Context, req *cart.EmptyCartReq, opts ...callopt.Option) (r *cart.EmptyCartResp, err error) {
	return &cart.EmptyCartResp{}, nil
}

// MockProductClient 模拟商品服务客户端
type MockProductClient struct{}

func (m *MockProductClient) BatchGetProducts(ctx context.Context, req *product.BatchGetProductsReq, opts ...callopt.Option) (r *product.BatchGetProductsResp, err error) {
	products := make([]*product.Product, 0)
	for _, id := range req.Ids {
		if id == 1001 {
			products = append(products, &product.Product{
				Id:    1001,
				Name:  "测试商品1",
				Price: 99.9,
			})
		} else if id == 1002 {
			products = append(products, &product.Product{
				Id:    1002,
				Name:  "测试商品2",
				Price: 199.9,
			})
		}
	}
	return &product.BatchGetProductsResp{
		Products: products,
	}, nil
}

func (m *MockProductClient) CreateProduct(ctx context.Context, req *product.CreateProductReq, opts ...callopt.Option) (r *product.CreateProductResp, err error) {
	return &product.CreateProductResp{}, nil
}

func (m *MockProductClient) DeleteProduct(ctx context.Context, req *product.DeleteProductReq, opts ...callopt.Option) (r *product.DeleteProductResp, err error) {
	return &product.DeleteProductResp{}, nil
}

// MockOrderClient 模拟订单服务客户端
type MockOrderClient struct{}

func (m *MockOrderClient) CreateOrder(ctx context.Context, req *order.CreateOrderReq, opts ...callopt.Option) (r *order.CreateOrderResp, err error) {
	return &order.CreateOrderResp{
		Order: &order.Order{
			OrderId:     "test_order_001",
			PaymentId:   "test_payment_001",
			TotalAmount: 399.7,
		},
	}, nil
}

func (m *MockOrderClient) CancelOrder(ctx context.Context, req *order.CancelOrderReq, opts ...callopt.Option) (r *order.CancelOrderResp, err error) {
	return &order.CancelOrderResp{}, nil
}

func (m *MockOrderClient) GetOrder(ctx context.Context, req *order.GetOrderReq, opts ...callopt.Option) (r *order.GetOrderResp, err error) {
	return &order.GetOrderResp{}, nil
}

func (m *MockOrderClient) UpdateOrder(ctx context.Context, req *order.UpdateOrderReq, opts ...callopt.Option) (r *order.UpdateOrderResp, err error) {
	return &order.UpdateOrderResp{}, nil
}

func (m *MockOrderClient) ListOrder(ctx context.Context, req *order.ListOrderReq, opts ...callopt.Option) (r *order.ListOrderResp, err error) {
	return &order.ListOrderResp{}, nil
}

func (m *MockOrderClient) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq, opts ...callopt.Option) (r *order.MarkOrderPaidResp, err error) {
	return &order.MarkOrderPaidResp{}, nil
}
func (m *MockProductClient) GetProduct(ctx context.Context, req *product.GetProductReq, opts ...callopt.Option) (r *product.GetProductResp, err error) {
	return &product.GetProductResp{}, nil
}
func (m *MockProductClient) ListProducts(ctx context.Context, req *product.ListProductsReq, opts ...callopt.Option) (r *product.ListProductsResp, err error) {
	return nil, nil
}

func (m *MockProductClient) SearchProducts(ctx context.Context, req *product.SearchProductsReq, opts ...callopt.Option) (r *product.SearchProductsResp, err error) {
	return nil, nil
}

func (m *MockProductClient) UpdateProduct(ctx context.Context, req *product.UpdateProductReq, opts ...callopt.Option) (r *product.UpdateProductResp, err error) {
	return nil, nil
}

func TestMain(m *testing.M) {
	// 替换为mock客户端
	rpc.CartClient = &MockCartClient{}
	rpc.ProductClient = &MockProductClient{}
	rpc.OrderClient = &MockOrderClient{}

	// 加载环境变量
	_ = godotenv.Load()

	// 运行测试
	code := m.Run()
	os.Exit(code)
}
