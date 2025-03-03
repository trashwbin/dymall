package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/joho/godotenv"
	"github.com/trashwbin/dymall/app/scheduler/biz/dal"
	"github.com/trashwbin/dymall/app/scheduler/infra/rpc"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

// MockOrderClient 模拟订单服务客户端
type MockOrderClient struct{}

func (m *MockOrderClient) GetOrder(ctx context.Context, req *order.GetOrderReq, opts ...callopt.Option) (r *order.GetOrderResp, err error) {
	// 模拟订单服务响应
	if req.OrderId == "test_order_1" {
		return &order.GetOrderResp{
			Order: &order.Order{
				OrderId:     "test_order_1",
				UserId:      1001,
				Status:      order.OrderStatus_ORDER_STATUS_PENDING,
				PaymentId:   "test_payment_1",
				TotalAmount: 99.9,
			},
		}, nil
	}
	if req.OrderId == "test_order_2" {
		return &order.GetOrderResp{
			Order: &order.Order{
				OrderId:     "test_order_2",
				UserId:      1002,
				Status:      order.OrderStatus_ORDER_STATUS_PAID,
				PaymentId:   "test_payment_2",
				TotalAmount: 199.9,
			},
		}, nil
	}
	return &order.GetOrderResp{}, nil
}

func (m *MockOrderClient) CancelOrder(ctx context.Context, req *order.CancelOrderReq, opts ...callopt.Option) (r *order.CancelOrderResp, err error) {
	return &order.CancelOrderResp{}, nil
}

func (m *MockOrderClient) CreateOrder(ctx context.Context, req *order.CreateOrderReq, opts ...callopt.Option) (r *order.CreateOrderResp, err error) {
	return nil, nil
}

func (m *MockOrderClient) UpdateOrder(ctx context.Context, req *order.UpdateOrderReq, opts ...callopt.Option) (r *order.UpdateOrderResp, err error) {
	return nil, nil
}

func (m *MockOrderClient) ListOrder(ctx context.Context, req *order.ListOrderReq, opts ...callopt.Option) (r *order.ListOrderResp, err error) {
	return nil, nil
}

func (m *MockOrderClient) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq, opts ...callopt.Option) (r *order.MarkOrderPaidResp, err error) {
	return nil, nil
}

// MockPaymentClient 模拟支付服务客户端
type MockPaymentClient struct{}

func (m *MockPaymentClient) GetPayment(ctx context.Context, req *payment.GetPaymentReq, opts ...callopt.Option) (r *payment.GetPaymentResp, err error) {
	// 模拟支付服务响应
	if req.PaymentId == "test_payment_1" {
		return &payment.GetPaymentResp{
			Payment: &payment.Payment{
				PaymentId: "test_payment_1",
				OrderId:   "test_order_1",
				UserId:    1001,
				Status:    payment.PaymentStatus_PAYMENT_STATUS_PENDING,
				Amount:    99.9,
			},
		}, nil
	}
	if req.PaymentId == "test_payment_2" {
		return &payment.GetPaymentResp{
			Payment: &payment.Payment{
				PaymentId: "test_payment_2",
				OrderId:   "test_order_2",
				UserId:    1002,
				Status:    payment.PaymentStatus_PAYMENT_STATUS_SUCCESS,
				Amount:    199.9,
			},
		}, nil
	}
	return &payment.GetPaymentResp{}, nil
}

func (m *MockPaymentClient) CreatePayment(ctx context.Context, req *payment.CreatePaymentReq, opts ...callopt.Option) (r *payment.CreatePaymentResp, err error) {
	return nil, nil
}

func (m *MockPaymentClient) ProcessPayment(ctx context.Context, req *payment.ProcessPaymentReq, opts ...callopt.Option) (r *payment.ProcessPaymentResp, err error) {
	return nil, nil
}

func (m *MockPaymentClient) CancelPayment(ctx context.Context, req *payment.CancelPaymentReq, opts ...callopt.Option) (r *payment.CancelPaymentResp, err error) {
	return &payment.CancelPaymentResp{}, nil
}

func TestMain(m *testing.M) {
	// 替换为mock客户端
	rpc.OrderClient = &MockOrderClient{}
	rpc.PaymentClient = &MockPaymentClient{}

	// 加载环境变量
	_ = godotenv.Load()

	// 初始化数据库
	dal.Init()

	// 运行测试
	code := m.Run()
	os.Exit(code)
}

// 辅助函数：创建测试任务
func createTestTask(t *testing.T) string {
	ctx := context.Background()
	s := NewCreateTaskService(ctx)
	resp, err := s.Run(&scheduler.CreateTaskReq{
		Type:      scheduler.TaskType_TASK_TYPE_ORDER_EXPIRATION,
		TargetId:  "test_order_1",
		ExecuteAt: time.Now().Add(time.Hour).Unix(),
	})
	if err != nil {
		t.Fatalf("create test task failed: %v", err)
	}
	return resp.Task.TaskId
}
