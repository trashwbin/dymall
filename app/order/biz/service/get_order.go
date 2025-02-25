package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/order/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/order/biz/dal/redis"
	"github.com/trashwbin/dymall/app/order/utils"
	pb "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
)

type GetOrderService struct {
	ctx       context.Context
	mysqlRepo *mysql.OrderRepo
	redisRepo *redis.OrderRepo
}

func NewGetOrderService(ctx context.Context) *GetOrderService {
	return &GetOrderService{
		ctx:       ctx,
		mysqlRepo: mysql.NewOrderRepo(),
		redisRepo: redis.NewOrderRepo(),
	}
}

// Run 获取订单详情
func (s *GetOrderService) Run(req *pb.GetOrderReq) (resp *pb.GetOrderResp, err error) {
	// 1. 先从缓存获取
	order, err := s.redisRepo.GetOrder(s.ctx, req.OrderId)
	if err == nil && order.UserID == int64(req.UserId) {
		return &pb.GetOrderResp{
			Order: convertModelToProto(order),
		}, nil
	}

	// 2. 缓存未命中，从数据库获取
	order, err = s.mysqlRepo.GetOrder(req.OrderId)
	if err != nil {
		klog.Errorf("get order from db failed: %v", err)
		return nil, utils.NewBizError(500, "获取订单失败")
	}

	// 3. 验证用户权限
	if order.UserID != int64(req.UserId) {
		return nil, utils.NewBizError(403, "无权访问该订单")
	}

	// 4. 设置缓存
	err = s.redisRepo.SetOrder(s.ctx, order)
	if err != nil {
		klog.Errorf("set order cache failed: %v", err)
		// 这里不需要返回错误，因为订单已经获取成功
	}

	return &pb.GetOrderResp{
		Order: convertModelToProto(order),
	}, nil
}
