package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/order/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/order/biz/dal/redis"
	"github.com/trashwbin/dymall/app/order/biz/model"
	"github.com/trashwbin/dymall/app/order/utils"
	pb "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
)

type UpdateOrderService struct {
	ctx       context.Context
	mysqlRepo *mysql.OrderRepo
	redisRepo *redis.OrderRepo
}

func NewUpdateOrderService(ctx context.Context) *UpdateOrderService {
	return &UpdateOrderService{
		ctx:       ctx,
		mysqlRepo: mysql.NewOrderRepo(),
		redisRepo: redis.NewOrderRepo(),
	}
}

// Run 更新订单信息（仅支持更新地址）
func (s *UpdateOrderService) Run(req *pb.UpdateOrderReq) (resp *pb.UpdateOrderResp, err error) {
	// 1. 获取订单
	order, err := s.mysqlRepo.GetOrder(req.OrderId)
	if err != nil {
		klog.Errorf("get order failed: %v", err)
		return nil, utils.NewBizError(500, "获取订单失败")
	}

	// 2. 验证用户权限
	if order.UserID != int64(req.UserId) {
		return nil, utils.NewBizError(403, "无权修改该订单")
	}

	// 3. 验证订单状态
	if !order.CanBeCanceled() {
		return nil, utils.NewBizError(400, "订单状态不允许修改")
	}

	// 4. 更新地址
	order.Address = &model.Address{
		OrderID:       order.OrderID,
		UserID:        order.UserID,
		StreetAddress: req.Address.StreetAddress,
		City:          req.Address.City,
		State:         req.Address.State,
		Country:       req.Address.Country,
		ZipCode:       fmt.Sprintf("%d", req.Address.ZipCode),
		CreatedAt:     order.Address.CreatedAt,
		UpdatedAt:     time.Now(),
	}
	order.UpdatedAt = time.Now()

	// 5. 保存更新
	err = s.mysqlRepo.UpdateOrder(order)
	if err != nil {
		klog.Errorf("update order failed: %v", err)
		return nil, utils.NewBizError(500, "更新订单失败")
	}

	// 6. 更新缓存
	err = s.redisRepo.SetOrder(s.ctx, order)
	if err != nil {
		klog.Errorf("update order cache failed: %v", err)
		// 这里不需要返回错误，因为订单已经更新成功
	}

	return &pb.UpdateOrderResp{
		Order: convertModelToProto(order),
	}, nil
}
