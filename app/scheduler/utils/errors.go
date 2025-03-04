package utils

import (
	"github.com/cloudwego/kitex/pkg/kerrors"
)

// NewBizError 创建业务错误
func NewBizError(code int32, msg string) error {
	return kerrors.NewBizStatusError(code, msg)
}

// MustHandleError 必须处理的错误
func MustHandleError(err error) {
	if err != nil {
		panic(err)
	}
}
