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
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type MarkOrderPaidService struct {
	ctx       context.Context
	mysqlRepo *mysql.OrderRepo
	redisRepo *redis.OrderRepo
}

func NewMarkOrderPaidService(ctx context.Context) *MarkOrderPaidService {
	return &MarkOrderPaidService{
		ctx:       ctx,
		mysqlRepo: mysql.NewOrderRepo(),
		redisRepo: redis.NewOrderRepo(),
	}
}

// Run 标记订单支付完成
func (s *MarkOrderPaidService) Run(req *pb.MarkOrderPaidReq) (resp *pb.MarkOrderPaidResp, err error) {
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
	if !order.CanBePaid() {
		return nil, utils.NewBizError(400, "订单状态不允许支付")
	}

	// 4. 更新订单状态
	order.Status = model.OrderStatusPaid
	now := time.Now()
	order.PaidAt = &now
	order.UpdatedAt = now

	// 5. 保存更新
	err = s.mysqlRepo.UpdateOrder(order)
	if err != nil {
		klog.Errorf("update order failed: %v", err)
		return nil, utils.NewBizError(500, "更新订单失败")
	}

	// 6. 取消订单过期任务
	_, err = rpc.SchedulerClient.CancelTask(s.ctx, &scheduler.CancelTaskReq{
		TaskId: order.OrderID, // 使用订单ID作为任务ID
	})
	if err != nil {
		klog.Errorf("cancel expire task failed: %v", err)
		// 这里不需要返回错误，因为订单已经更新成功
	}

	// 7. 更新缓存
	err = s.redisRepo.SetOrder(s.ctx, order)
	if err != nil {
		klog.Errorf("update order cache failed: %v", err)
		// 这里不需要返回错误，因为订单已经更新成功
	}

	return &pb.MarkOrderPaidResp{}, nil
}
