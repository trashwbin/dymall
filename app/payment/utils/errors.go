package utils

import (
	"github.com/cloudwego/kitex/pkg/kerrors"
)

const (
	// 支付相关错误码
	ErrPaymentNotFound      = 10001 // 支付单不存在
	ErrPaymentExpired       = 10002 // 支付单已过期
	ErrPaymentStatusInvalid = 10003 // 支付单状态不正确
	ErrPaymentAmountInvalid = 10004 // 支付金额不正确
	ErrCreditCardInvalid    = 10005 // 信用卡信息无效
)

var (
	// 错误信息映射
	errMsgMap = map[int32]string{
		ErrPaymentNotFound:      "payment not found",
		ErrPaymentExpired:       "payment expired",
		ErrPaymentStatusInvalid: "payment status invalid",
		ErrPaymentAmountInvalid: "payment amount invalid",
		ErrCreditCardInvalid:    "credit card info invalid",
	}
)

// MustHandleError 必须处理的错误
func MustHandleError(err error) {
	if err != nil {
		panic(err)
	}
}

// NewBizError 创建业务错误
func NewBizError(code int32, msg string) error {
	return kerrors.NewBizStatusError(code, msg)
}

// NewPaymentError 创建支付错误
func NewPaymentError(code int32) error {
	return NewBizError(code, errMsgMap[code])
}
