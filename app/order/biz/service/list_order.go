package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/order/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/order/biz/dal/redis"
	"github.com/trashwbin/dymall/app/order/biz/model"
	"github.com/trashwbin/dymall/app/order/utils"
	pb "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
)

type ListOrderService struct {
	ctx       context.Context
	mysqlRepo *mysql.OrderRepo
	redisRepo *redis.OrderRepo
}

func NewListOrderService(ctx context.Context) *ListOrderService {
	return &ListOrderService{
		ctx:       ctx,
		mysqlRepo: mysql.NewOrderRepo(),
		redisRepo: redis.NewOrderRepo(),
	}
}

// Run 获取订单列表
func (s *ListOrderService) Run(req *pb.ListOrderReq) (resp *pb.ListOrderResp, err error) {
	// 1. 先从缓存获取
	orders, err := s.redisRepo.GetUserOrders(s.ctx, int64(req.UserId))
	if err == nil && len(orders) > 0 {
		orderProtos := make([]*pb.Order, 0, len(orders))
		for _, order := range orders {
			orderProtos = append(orderProtos, convertModelToProto(order))
		}
		return &pb.ListOrderResp{
			Orders: orderProtos,
		}, nil
	}

	// 2. 缓存未命中，从数据库获取
	orders, err = s.mysqlRepo.ListOrders(int64(req.UserId))
	if err != nil {
		klog.Errorf("list orders from db failed: %v", err)
		return nil, utils.NewBizError(500, "获取订单列表失败")
	}

	// 如果用户没有订单，直接返回空列表
	if len(orders) == 0 {
		return &pb.ListOrderResp{
			Orders: make([]*pb.Order, 0),
		}, nil
	}

	// 3. 过滤掉无效的订单
	validOrders := make([]*model.Order, 0, len(orders))
	for _, order := range orders {
		if order != nil && order.Address != nil && len(order.OrderItems) > 0 {
			validOrders = append(validOrders, order)
		}
	}

	// 4. 设置缓存
	orderIDs := make([]string, 0, len(validOrders))
	for _, order := range validOrders {
		orderIDs = append(orderIDs, order.OrderID)
		err = s.redisRepo.SetOrder(s.ctx, order)
		if err != nil {
			klog.Errorf("set order cache failed: %v", err)
			continue
		}
	}
	err = s.redisRepo.SetUserOrders(s.ctx, int64(req.UserId), orderIDs)
	if err != nil {
		klog.Errorf("set user orders cache failed: %v", err)
		// 这里不需要返回错误，因为订单列表已经获取成功
	}

	// 5. 返回响应
	orderProtos := make([]*pb.Order, 0, len(validOrders))
	for _, order := range validOrders {
		orderProtos = append(orderProtos, convertModelToProto(order))
	}
	return &pb.ListOrderResp{
		Orders: orderProtos,
	}, nil
}
