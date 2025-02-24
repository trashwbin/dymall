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

type ProcessPaymentService struct {
	ctx       context.Context
	mysqlRepo *mysql.PaymentRepo
	redisRepo *redis.PaymentRepo
}

func NewProcessPaymentService(ctx context.Context) *ProcessPaymentService {
	return &ProcessPaymentService{
		ctx:       ctx,
		mysqlRepo: mysql.NewPaymentRepo(),
		redisRepo: redis.NewPaymentRepo(),
	}
}

func (s *ProcessPaymentService) Run(req *payment.ProcessPaymentReq) (resp *payment.ProcessPaymentResp, err error) {
	klog.CtxInfof(s.ctx, "ProcessPaymentService - Run: req=%+v", req)

	// 1. 获取支付单
	paymentModel, err := s.mysqlRepo.GetPaymentByID(req.PaymentId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "ProcessPaymentService - GetPaymentByID failed: %v", err)
		return nil, utils.NewPaymentError(utils.ErrPaymentNotFound)
	}

	// 2. 验证支付单状态
	if !paymentModel.CanPay() {
		klog.CtxErrorf(s.ctx, "ProcessPaymentService - payment cannot be paid: status=%d, expired=%v",
			paymentModel.Status, paymentModel.IsExpired())
		if paymentModel.IsExpired() {
			return nil, utils.NewPaymentError(utils.ErrPaymentExpired)
		}
		return nil, utils.NewPaymentError(utils.ErrPaymentStatusInvalid)
	}

	// 3. 验证信用卡信息
	creditCard := &model.CreditCard{
		Number:          req.CreditCard.CreditCardNumber,
		CVV:             req.CreditCard.CreditCardCvv,
		ExpirationYear:  req.CreditCard.CreditCardExpirationYear,
		ExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
	}
	if valid, errMsg := creditCard.ValidateCreditCard(); !valid {
		klog.CtxErrorf(s.ctx, "ProcessPaymentService - invalid credit card info: %s", errMsg)
		return nil, utils.NewBizError(utils.ErrCreditCardInvalid, errMsg)
	}

	// 4. 处理支付（这里应该调用实际的支付网关）
	now := time.Now()
	paymentModel.Status = model.PayStatusSuccess
	paymentModel.UpdatedAt = now
	paymentModel.PaidAt = now

	// 5. 更新MySQL
	if err := s.mysqlRepo.UpdatePayment(paymentModel); err != nil {
		klog.CtxErrorf(s.ctx, "ProcessPaymentService - UpdatePayment failed: %v", err)
		return nil, err
	}

	// 6. 更新Redis缓存
	if err := s.redisRepo.SetPayment(s.ctx, paymentModel); err != nil {
		klog.CtxWarnf(s.ctx, "ProcessPaymentService - Redis SetPayment failed: %v", err)
	}

	// 7. 构建响应
	resp = &payment.ProcessPaymentResp{
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

	klog.CtxInfof(s.ctx, "ProcessPaymentService - Run success: paymentID=%s", paymentModel.PaymentID)
	return resp, nil
}
