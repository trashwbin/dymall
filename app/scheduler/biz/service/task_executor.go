package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/trashwbin/dymall/app/scheduler/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/scheduler/biz/dal/redis"
	"github.com/trashwbin/dymall/app/scheduler/biz/model"
	"github.com/trashwbin/dymall/app/scheduler/infra/rpc"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment/paymentservice"
)

const (
	// 每次获取待执行任务的数量
	batchSize = 100
	// 执行间隔
	executionInterval = 1 * time.Second
)

// TaskExecutor 任务执行器
type TaskExecutor struct {
	mysqlRepo     *mysql.TaskRepo
	redisRepo     *redis.TaskRepo
	orderClient   orderservice.Client
	paymentClient paymentservice.Client
}

// NewTaskExecutor 创建任务执行器
func NewTaskExecutor() *TaskExecutor {
	return &TaskExecutor{
		mysqlRepo:     mysql.NewTaskRepo(),
		redisRepo:     redis.NewTaskRepo(),
		orderClient:   rpc.OrderClient,
		paymentClient: rpc.PaymentClient,
	}
}

// Start 启动任务执行器
func (e *TaskExecutor) Start(ctx context.Context) {
	ticker := time.NewTicker(executionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := e.executeBatch(ctx); err != nil {
				log.Printf("execute batch failed: %v", err)
			}
		}
	}
}

// executeBatch 批量执行任务
func (e *TaskExecutor) executeBatch(ctx context.Context) error {
	// 1. 获取待执行的任务
	tasks, err := e.mysqlRepo.GetPendingTasks(time.Now(), batchSize)
	if err != nil {
		return fmt.Errorf("get pending tasks failed: %w", err)
	}

	// 2. 遍历执行任务
	for _, task := range tasks {
		// 2.1 尝试获取任务锁
		locked, err := e.redisRepo.AcquireLock(ctx, task.ID)
		if err != nil || !locked {
			continue // 获取锁失败，跳过该任务
		}

		// 2.2 执行任务并自动释放锁
		go e.executeTask(ctx, task)
	}

	return nil
}

// executeTask 执行单个任务
func (e *TaskExecutor) executeTask(ctx context.Context, task *model.Task) {
	defer e.redisRepo.ReleaseLock(ctx, task.ID)

	// 1. 标记任务开始执行
	task.MarkRunning()
	if err := e.mysqlRepo.UpdateTask(task); err != nil {
		log.Printf("update task status to running failed: %v", err)
		return
	}

	// 2. 根据任务类型执行不同的逻辑
	var err error
	switch task.Type {
	case model.TaskTypeOrderExpiration:
		err = e.handleOrderExpiration(ctx, task)
	default:
		err = fmt.Errorf("unsupported task type: %v", task.Type)
	}

	// 3. 更新任务状态
	if err != nil {
		task.MarkFailed(err.Error())
	} else {
		task.MarkCompleted()
	}

	if err := e.mysqlRepo.UpdateTask(task); err != nil {
		log.Printf("update task status failed: %v", err)
	}
}

// handleOrderExpiration 处理订单过期任务
func (e *TaskExecutor) handleOrderExpiration(ctx context.Context, task *model.Task) error {
	// 1. 先查询订单状态
	getOrderReq := &order.GetOrderReq{
		OrderId: task.TargetID,
	}
	orderResp, err := e.orderClient.GetOrder(ctx, getOrderReq)
	if err != nil {
		return fmt.Errorf("get order failed: %w", err)
	}

	// 2. 检查订单状态，只有待支付状态的订单才能被取消
	if orderResp.Order.Status != order.OrderStatus_ORDER_STATUS_PENDING {
		// 订单已经不是待支付状态，说明可能已经支付或已被取消
		// 将任务标记为完成并记录原因
		if task.Metadata == nil {
			task.Metadata = make(map[string]string)
		}
		task.Metadata["skip_reason"] = fmt.Sprintf("order status is %v", orderResp.Order.Status)
		return nil
	}

	// 3. 检查支付状态
	getPaymentReq := &payment.GetPaymentReq{
		PaymentId: orderResp.Order.PaymentId,
	}
	paymentResp, err := e.paymentClient.GetPayment(ctx, getPaymentReq)
	if err != nil {
		return fmt.Errorf("get payment failed: %w", err)
	}

	// 4. 如果支付单已支付，则不取消订单
	if paymentResp.Payment.Status == payment.PaymentStatus_PAYMENT_STATUS_SUCCESS {
		if task.Metadata == nil {
			task.Metadata = make(map[string]string)
		}
		task.Metadata["skip_reason"] = "payment already succeeded"
		return nil
	}

	// 5. 调用 OrderService 取消订单（会级联取消支付单）
	cancelOrderReq := &order.CancelOrderReq{
		OrderId: task.TargetID,
		Cascade: true, // 级联取消支付单
	}

	_, err = e.orderClient.CancelOrder(ctx, cancelOrderReq)
	if err != nil {
		return fmt.Errorf("cancel order failed: %w", err)
	}

	return nil
}
