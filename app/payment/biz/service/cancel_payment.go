package service

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/trashwbin/dymall/app/payment/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/payment/biz/dal/redis"
	"github.com/trashwbin/dymall/app/payment/biz/model"
	"github.com/trashwbin/dymall/app/payment/utils"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

type CancelPaymentService struct {
	ctx       context.Context
	mysqlRepo *mysql.PaymentRepo
	redisRepo *redis.PaymentRepo
}

func NewCancelPaymentService(ctx context.Context) *CancelPaymentService {
	return &CancelPaymentService{
		ctx:       ctx,
		mysqlRepo: mysql.NewPaymentRepo(),
		redisRepo: redis.NewPaymentRepo(),
	}
}

func (s *CancelPaymentService) Run(req *payment.CancelPaymentReq) (resp *payment.CancelPaymentResp, err error) {
	klog.CtxInfof(s.ctx, "CancelPaymentService - Run: req=%+v", req)

	// 1. 获取支付单
	paymentModel, err := s.mysqlRepo.GetPaymentByID(req.PaymentId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "CancelPaymentService - GetPaymentByID failed: %v", err)
		return nil, utils.NewPaymentError(utils.ErrPaymentNotFound)
	}

	// 2. 验证用户ID
	if paymentModel.UserID != int64(req.UserId) {
		klog.CtxErrorf(s.ctx, "CancelPaymentService - user_id not match: want=%d, got=%d",
			paymentModel.UserID, req.UserId)
		return nil, utils.NewPaymentError(utils.ErrPaymentNotFound)
	}

	// 3. 验证支付单状态
	if paymentModel.Status != model.PayStatusPending {
		klog.CtxErrorf(s.ctx, "CancelPaymentService - payment status invalid: %d", paymentModel.Status)
		return nil, utils.NewPaymentError(utils.ErrPaymentStatusInvalid)
	}

	// 4. 更新支付单状态
	now := time.Now()
	paymentModel.Status = model.PayStatusCancelled
	paymentModel.UpdatedAt = now

	// 5. 更新MySQL
	if err := s.mysqlRepo.UpdatePayment(paymentModel); err != nil {
		klog.CtxErrorf(s.ctx, "CancelPaymentService - UpdatePayment failed: %v", err)
		return nil, err
	}

	// 6. 更新Redis缓存
	if err := s.redisRepo.SetPayment(s.ctx, paymentModel); err != nil {
		klog.CtxWarnf(s.ctx, "CancelPaymentService - Redis SetPayment failed: %v", err)
	}

	resp = &payment.CancelPaymentResp{}
	klog.CtxInfof(s.ctx, "CancelPaymentService - Run success: paymentID=%s", paymentModel.PaymentID)
	return resp, nil
}
