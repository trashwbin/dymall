package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/payment/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/payment/biz/dal/redis"
	"github.com/trashwbin/dymall/app/payment/utils"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

type GetPaymentService struct {
	ctx       context.Context
	mysqlRepo *mysql.PaymentRepo
	redisRepo *redis.PaymentRepo
}

func NewGetPaymentService(ctx context.Context) *GetPaymentService {
	return &GetPaymentService{
		ctx:       ctx,
		mysqlRepo: mysql.NewPaymentRepo(),
		redisRepo: redis.NewPaymentRepo(),
	}
}

func (s *GetPaymentService) Run(req *payment.GetPaymentReq) (resp *payment.GetPaymentResp, err error) {
	klog.CtxInfof(s.ctx, "GetPaymentService - Run: req=%+v", req)

	// 1. 尝试从Redis获取
	paymentModel, err := s.redisRepo.GetPayment(s.ctx, req.PaymentId)
	if err != nil {
		klog.CtxDebugf(s.ctx, "GetPaymentService - Redis GetPayment failed: %v", err)
		// 2. Redis不存在，从MySQL获取
		paymentModel, err = s.mysqlRepo.GetPaymentByID(req.PaymentId)
		if err != nil {
			klog.CtxErrorf(s.ctx, "GetPaymentService - MySQL GetPaymentByID failed: %v", err)
			return nil, utils.NewPaymentError(utils.ErrPaymentNotFound)
		}
		// 3. 写入Redis缓存
		if err := s.redisRepo.SetPayment(s.ctx, paymentModel); err != nil {
			klog.CtxWarnf(s.ctx, "GetPaymentService - Redis SetPayment failed: %v", err)
		}
	}

	// 4. 验证用户ID
	if paymentModel.UserID != int64(req.UserId) {
		klog.CtxErrorf(s.ctx, "GetPaymentService - user_id not match: want=%d, got=%d",
			paymentModel.UserID, req.UserId)
		return nil, utils.NewPaymentError(utils.ErrPaymentNotFound)
	}

	// 5. 构建响应
	resp = &payment.GetPaymentResp{
		Payment: &payment.Payment{
			PaymentId: paymentModel.PaymentID,
			OrderId:   paymentModel.OrderID,
			UserId:    uint32(paymentModel.UserID),
			Amount:    float32(paymentModel.Amount),
			Currency:  paymentModel.Currency,
			Status:    payment.PaymentStatus(paymentModel.Status),
			CreatedAt: paymentModel.CreatedAt.Unix(),
			UpdatedAt: paymentModel.UpdatedAt.Unix(),
			ExpireAt:  paymentModel.ExpireAt.Unix(),
			PaidAt:    paymentModel.PaidAt.Unix(),
		},
	}

	klog.CtxInfof(s.ctx, "GetPaymentService - Run success: paymentID=%s", paymentModel.PaymentID)
	return resp, nil
}
