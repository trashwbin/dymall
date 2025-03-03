package service

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/order/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/order/biz/dal/redis"
	"github.com/trashwbin/dymall/app/order/biz/model"
	"github.com/trashwbin/dymall/app/order/infra/rpc"
	"github.com/trashwbin/dymall/app/order/utils"
	pb "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type CancelOrderService struct {
	ctx       context.Context
	mysqlRepo *mysql.OrderRepo
	redisRepo *redis.OrderRepo
}

func NewCancelOrderService(ctx context.Context) *CancelOrderService {
	return &CancelOrderService{
		ctx:       ctx,
		mysqlRepo: mysql.NewOrderRepo(),
		redisRepo: redis.NewOrderRepo(),
	}
}

// Run 取消订单
func (s *CancelOrderService) Run(req *pb.CancelOrderReq) (resp *pb.CancelOrderResp, err error) {
	// 1. 获取订单
	order, err := s.mysqlRepo.GetOrder(req.OrderId)
	if err != nil {
		klog.Errorf("get order failed: %v", err)
		return nil, utils.NewBizError(500, "获取订单失败")
	}

	// 2. 验证用户权限
	if order.UserID != int64(req.UserId) {
		return nil, utils.NewBizError(403, "无权操作该订单")
	}

	// 3. 验证订单状态
	if !order.CanBeCanceled() {
		return nil, utils.NewBizError(400, "订单状态不允许取消")
	}

	// 4. 更新订单状态
	order.Status = model.OrderStatusCanceled
	order.UpdatedAt = time.Now()

	// 5. 保存更新
	err = s.mysqlRepo.UpdateOrder(order)
	if err != nil {
		klog.Errorf("update order failed: %v", err)
		return nil, utils.NewBizError(500, "更新订单失败")
	}

	// 6. 如果需要级联取消支付单
	if req.Cascade && order.PaymentID != "" {
		_, err = rpc.PaymentClient.CancelPayment(s.ctx, &payment.CancelPaymentReq{
			PaymentId: order.PaymentID,
			UserId:    req.UserId,
		})
		if err != nil {
			klog.Errorf("cancel payment failed: %v", err)
			// 这里不需要返回错误，因为订单已经取消成功
		}
	}

	// 7. 取消订单过期任务
	_, err = rpc.SchedulerClient.CancelTask(s.ctx, &scheduler.CancelTaskReq{
		TaskId: order.OrderID, // 使用订单ID作为任务ID
	})
	if err != nil {
		klog.Errorf("cancel expire task failed: %v", err)
		// 这里不需要返回错误，因为订单已经取消成功
	}

	// 8. 更新缓存
	err = s.redisRepo.SetOrder(s.ctx, order)
	if err != nil {
		klog.Errorf("update order cache failed: %v", err)
		// 这里不需要返回错误，因为订单已经取消成功
	}

	return &pb.CancelOrderResp{}, nil
}
