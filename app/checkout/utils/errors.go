package utils

import (
	"github.com/cloudwego/kitex/pkg/kerrors"
)

const (
	// 结算相关错误码
	ErrCartEmpty         = 20001 // 购物车为空
	ErrInvalidQuantity   = 20002 // 无效的商品数量
	ErrInsufficientStock = 20003 // 商品库存不足
	ErrProductNotFound   = 20004 // 商品不存在
	ErrInvalidZipCode    = 20005 // 无效的邮编
	ErrCheckoutNotFound  = 20006 // 结算单不存在
)

var (
	// 错误信息映射
	errMsgMap = map[int32]string{
		ErrCartEmpty:         "cart is empty",
		ErrInvalidQuantity:   "invalid quantity",
		ErrInsufficientStock: "insufficient stock",
		ErrProductNotFound:   "product not found",
		ErrInvalidZipCode:    "invalid zip code",
		ErrCheckoutNotFound:  "checkout not found",
	}
)

// NewCheckoutError 创建结算错误
func NewCheckoutError(code int32) error {
	return kerrors.NewBizStatusError(code, errMsgMap[code])
}

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
