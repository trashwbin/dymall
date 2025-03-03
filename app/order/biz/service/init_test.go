package service

import (
	"context"
	"os"
	"testing"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/joho/godotenv"
	"github.com/trashwbin/dymall/app/order/biz/dal"
	"github.com/trashwbin/dymall/app/order/infra/rpc"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

// MockCartClient 模拟购物车服务客户端
type MockCartClient struct{}

func (m *MockCartClient) GetCart(ctx context.Context, req *cart.GetCartReq, opts ...callopt.Option) (r *cart.GetCartResp, err error) {
	return &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: req.UserId,
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

func (m *MockCartClient) AddItem(ctx context.Context, req *cart.AddItemReq, opts ...callopt.Option) (r *cart.AddItemResp, err error) {
	return nil, nil
}

func (m *MockCartClient) EmptyCart(ctx context.Context, req *cart.EmptyCartReq, opts ...callopt.Option) (r *cart.EmptyCartResp, err error) {
	return &cart.EmptyCartResp{}, nil
}

// MockProductClient 模拟商品服务客户端
type MockProductClient struct{}

func (m *MockProductClient) BatchGetProducts(ctx context.Context, req *product.BatchGetProductsReq, opts ...callopt.Option) (r *product.BatchGetProductsResp, err error) {
	products := make([]*product.Product, 0, len(req.Ids))
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

func (m *MockProductClient) GetProduct(ctx context.Context, req *product.GetProductReq, opts ...callopt.Option) (r *product.GetProductResp, err error) {
	return nil, nil
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

// MockSchedulerClient 模拟调度器服务客户端
type MockSchedulerClient struct{}

func (m *MockSchedulerClient) CreateTask(ctx context.Context, req *scheduler.CreateTaskReq, opts ...callopt.Option) (r *scheduler.CreateTaskResp, err error) {
	return &scheduler.CreateTaskResp{
		Task: &scheduler.Task{
			TaskId:    req.TargetId,
			Type:      req.Type,
			Status:    scheduler.TaskStatus_TASK_STATUS_PENDING,
			TargetId:  req.TargetId,
			Metadata:  req.Metadata,
			ExecuteAt: req.ExecuteAt,
		},
	}, nil
}

func (m *MockSchedulerClient) CancelTask(ctx context.Context, req *scheduler.CancelTaskReq, opts ...callopt.Option) (r *scheduler.CancelTaskResp, err error) {
	return &scheduler.CancelTaskResp{}, nil
}

func (m *MockSchedulerClient) GetTask(ctx context.Context, req *scheduler.GetTaskReq, opts ...callopt.Option) (r *scheduler.GetTaskResp, err error) {
	return nil, nil
}

func (m *MockSchedulerClient) ExecuteTask(ctx context.Context, req *scheduler.ExecuteTaskReq, opts ...callopt.Option) (r *scheduler.ExecuteTaskResp, err error) {
	return nil, nil
}

// MockPaymentClient 模拟支付服务客户端
type MockPaymentClient struct{}

func (m *MockPaymentClient) CreatePayment(ctx context.Context, req *payment.CreatePaymentReq, opts ...callopt.Option) (r *payment.CreatePaymentResp, err error) {
	return &payment.CreatePaymentResp{
		Payment: &payment.Payment{
			PaymentId: "test_payment_id",
			OrderId:   req.OrderId,
			UserId:    req.UserId,
			Amount:    req.Amount,
			Currency:  req.Currency,
			Status:    payment.PaymentStatus_PAYMENT_STATUS_PENDING,
			ExpireAt:  req.ExpireAt,
		},
	}, nil
}

func (m *MockPaymentClient) ProcessPayment(ctx context.Context, req *payment.ProcessPaymentReq, opts ...callopt.Option) (r *payment.ProcessPaymentResp, err error) {
	return nil, nil
}

func (m *MockPaymentClient) CancelPayment(ctx context.Context, req *payment.CancelPaymentReq, opts ...callopt.Option) (r *payment.CancelPaymentResp, err error) {
	return &payment.CancelPaymentResp{}, nil
}

func (m *MockPaymentClient) GetPayment(ctx context.Context, req *payment.GetPaymentReq, opts ...callopt.Option) (r *payment.GetPaymentResp, err error) {
	return nil, nil
}

func TestMain(m *testing.M) {
	// 替换为mock客户端
	rpc.CartClient = &MockCartClient{}
	rpc.ProductClient = &MockProductClient{}
	rpc.SchedulerClient = &MockSchedulerClient{}
	rpc.PaymentClient = &MockPaymentClient{}

	// 加载环境变量
	_ = godotenv.Load()

	// 初始化数据库
	dal.Init()

	// 运行测试
	code := m.Run()
	os.Exit(code)
}
