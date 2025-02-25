package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/order/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/order/biz/dal/redis"
	"github.com/trashwbin/dymall/app/order/biz/model"
	"github.com/trashwbin/dymall/app/order/infra/mq"
	"github.com/trashwbin/dymall/app/order/infra/rpc"
	"github.com/trashwbin/dymall/app/order/utils"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
	pb "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

type CreateOrderService struct {
	ctx       context.Context
	mysqlRepo *mysql.OrderRepo
	redisRepo *redis.OrderRepo
}

func NewCreateOrderService(ctx context.Context) *CreateOrderService {
	return &CreateOrderService{
		ctx:       ctx,
		mysqlRepo: mysql.NewOrderRepo(),
		redisRepo: redis.NewOrderRepo(),
	}
}

// Run 创建订单
func (s *CreateOrderService) Run(req *pb.CreateOrderReq) (resp *pb.CreateOrderResp, err error) {
	// 1. 验证请求参数
	if req == nil {
		return nil, utils.NewBizError(400, "请求参数不能为空")
	}
	if req.UserId == 0 {
		return nil, utils.NewBizError(400, "用户ID不能为空")
	}
	if req.Address == nil {
		return nil, utils.NewBizError(400, "地址信息不能为空")
	}
	if req.Email == "" {
		return nil, utils.NewBizError(400, "邮箱不能为空")
	}

	// 2. 获取购物车商品
	cartResp, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		klog.Errorf("get cart failed: %v", err)
		return nil, utils.NewBizError(500, "获取购物车失败")
	}

	if len(cartResp.Cart.Items) == 0 {
		return nil, utils.NewBizError(400, "购物车为空")
	}

	// 3. 获取商品信息
	productIDs := make([]uint32, 0, len(cartResp.Cart.Items))
	for _, item := range cartResp.Cart.Items {
		productIDs = append(productIDs, item.ProductId)
	}

	productsResp, err := rpc.ProductClient.BatchGetProducts(s.ctx, &product.BatchGetProductsReq{
		Ids: productIDs,
	})
	if err != nil {
		klog.Errorf("get products failed: %v", err)
		return nil, utils.NewBizError(500, "获取商品信息失败")
	}

	// 4. 创建订单
	orderModel := &model.Order{
		OrderID:   generateOrderID(req.UserId),
		UserID:    int64(req.UserId),
		Status:    model.OrderStatusPending,
		Currency:  "CNY", // 默认使用人民币
		Email:     req.Email,
		ExpireAt:  time.Now().Add(30 * time.Minute), // 订单30分钟后过期
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Address: &model.Address{
			UserID:        int64(req.UserId),
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       fmt.Sprintf("%d", req.Address.ZipCode), // 转换为字符串
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	// 5. 添加订单商品
	orderItems := make([]*model.OrderItem, 0, len(cartResp.Cart.Items))
	for _, item := range cartResp.Cart.Items {
		// 查找商品价格
		var price float64
		for _, p := range productsResp.Products {
			if p.Id == item.ProductId {
				price = float64(p.Price)
				break
			}
		}

		orderItem := &model.OrderItem{
			OrderID:   orderModel.OrderID,
			ProductID: int64(item.ProductId),
			Quantity:  item.Quantity,
			Price:     price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		orderItems = append(orderItems, orderItem)
	}
	orderModel.OrderItems = orderItems
	orderModel.TotalAmount = orderModel.CalculateTotal()

	// 6. 保存订单
	createdOrder, err := s.mysqlRepo.CreateOrder(orderModel)
	if err != nil {
		klog.Errorf("create order failed: %v", err)
		return nil, utils.NewBizError(500, "创建订单失败")
	}

	// 7. 创建订单过期任务
	_, err = rpc.SchedulerClient.CreateTask(s.ctx, &scheduler.CreateTaskReq{
		Type:      scheduler.TaskType_TASK_TYPE_ORDER_EXPIRATION,
		TargetId:  createdOrder.OrderID,
		ExecuteAt: createdOrder.ExpireAt.Unix(),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", createdOrder.UserID),
		},
	})
	if err != nil {
		klog.Errorf("create expire task failed: %v", err)
		// 这里不需要返回错误，因为订单已经创建成功
	}

	// 8. 发送订单创建事件
	err = mq.Nc.Publish("order.created", []byte(createdOrder.OrderID))
	if err != nil {
		klog.Errorf("publish order created event failed: %v", err)
		// 这里不需要返回错误，因为订单已经创建成功
	}

	// 9. 清空购物车
	_, err = rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		klog.Errorf("empty cart failed: %v", err)
		// 这里不需要返回错误，因为订单已经创建成功
	}

	// 10. 设置订单缓存
	err = s.redisRepo.SetOrder(s.ctx, createdOrder)
	if err != nil {
		klog.Errorf("set order cache failed: %v", err)
		// 这里不需要返回错误，因为订单已经创建成功
	}

	// 11. 返回响应
	return &pb.CreateOrderResp{
		Order: convertModelToProto(createdOrder),
	}, nil
}

// generateOrderID 生成订单号
func generateOrderID(userID uint32) string {
	return fmt.Sprintf("%d%d", userID, time.Now().UnixNano())
}

// convertModelToProto 将领域模型转换为 proto 模型
func convertModelToProto(order *model.Order) *pb.Order {
	items := make([]*pb.OrderItem, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		items = append(items, &pb.OrderItem{
			Item: &cart.CartItem{
				ProductId: uint32(item.ProductID),
				Quantity:  item.Quantity,
			},
			Cost: float32(item.Price),
		})
	}

	return &pb.Order{
		OrderId:      order.OrderID,
		UserId:       uint32(order.UserID),
		UserCurrency: order.Currency,
		Address: &pb.Address{
			StreetAddress: order.Address.StreetAddress,
			City:          order.Address.City,
			State:         order.Address.State,
			Country:       order.Address.Country,
			ZipCode:       int32(utils.MustAtoi(order.Address.ZipCode)),
		},
		Email:       order.Email,
		OrderItems:  items,
		Status:      pb.OrderStatus(order.Status),
		TotalAmount: float32(order.TotalAmount),
		CreatedAt:   order.CreatedAt.Unix(),
		UpdatedAt:   order.UpdatedAt.Unix(),
		ExpireAt:    order.ExpireAt.Unix(),
		PaymentId:   order.PaymentID,
	}
}
