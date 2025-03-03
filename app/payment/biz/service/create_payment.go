package service

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/uuid"
	"github.com/trashwbin/dymall/app/payment/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/payment/biz/dal/redis"
	"github.com/trashwbin/dymall/app/payment/biz/model"
	"github.com/trashwbin/dymall/app/payment/utils"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

type CreatePaymentService struct {
	ctx       context.Context
	mysqlRepo *mysql.PaymentRepo
	redisRepo *redis.PaymentRepo
}

func NewCreatePaymentService(ctx context.Context) *CreatePaymentService {
	return &CreatePaymentService{
		ctx:       ctx,
		mysqlRepo: mysql.NewPaymentRepo(),
		redisRepo: redis.NewPaymentRepo(),
	}
}

func (s *CreatePaymentService) Run(req *payment.CreatePaymentReq) (resp *payment.CreatePaymentResp, err error) {
	klog.CtxInfof(s.ctx, "CreatePaymentService - Run: req=%+v", req)

	// 1. 验证请求参数
	if req.Amount <= 0 {
		klog.CtxErrorf(s.ctx, "CreatePaymentService - invalid amount: %v", req.Amount)
		return nil, utils.NewPaymentError(utils.ErrPaymentAmountInvalid)
	}

	expireAt := time.Unix(req.ExpireAt, 0)
	if expireAt.Before(time.Now()) {
		klog.CtxErrorf(s.ctx, "CreatePaymentService - invalid expire time: %v", expireAt)
		return nil, utils.NewPaymentError(utils.ErrPaymentExpired)
	}

	// 2. 创建支付单模型
	now := time.Now()
	paymentModel := &model.Payment{
		PaymentID: uuid.New().String(),
		OrderID:   req.OrderId,
		UserID:    int64(req.UserId),
		Amount:    float64(req.Amount),
		Currency:  req.Currency,
		Status:    model.PayStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
		ExpireAt:  expireAt,
	}

	// 3. 保存到MySQL
	createdPayment, err := s.mysqlRepo.CreatePayment(paymentModel)
	if err != nil {
		klog.CtxErrorf(s.ctx, "CreatePaymentService - MySQL CreatePayment failed: %v", err)
		return nil, err
	}

	// 4. 保存到Redis缓存
	if err := s.redisRepo.SetPayment(s.ctx, createdPayment); err != nil {
		klog.CtxWarnf(s.ctx, "CreatePaymentService - Redis SetPayment failed: %v", err)
	}

	// 5. 构建响应
	resp = &payment.CreatePaymentResp{
		Payment: &payment.Payment{
			PaymentId: createdPayment.PaymentID,
			OrderId:   createdPayment.OrderID,
			UserId:    uint32(createdPayment.UserID),
			Amount:    float32(createdPayment.Amount),
			Currency:  createdPayment.Currency,
			Status:    payment.PaymentStatus(createdPayment.Status),
			CreatedAt: createdPayment.CreatedAt.Unix(),
			UpdatedAt: createdPayment.UpdatedAt.Unix(),
			ExpireAt:  createdPayment.ExpireAt.Unix(),
		},
	}

	klog.CtxInfof(s.ctx, "CreatePaymentService - Run success: paymentID=%s", createdPayment.PaymentID)
	return resp, nil
}
