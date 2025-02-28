package service

import (
	"context"
	"fmt"
	"github.com/trashwbin/dymall/app/user/biz/dal/mysql" // 引入 mysql 包来操作数据库
	user "github.com/trashwbin/dymall/rpc_gen/kitex_gen/user"
	"gorm.io/gorm"
)

type GetUserService struct {
	ctx context.Context
}

// NewGetUserService 创建新的 GetUserService 实例
func NewGetUserInfoService(ctx context.Context) *GetUserService {
	return &GetUserService{ctx: ctx}
}

// Run 获取用户信息
func (s *GetUserService) Run(req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	// 检查 user_id 是否提供
	if req.UserId == 0 {
		return &user.UserInfoResponse{
			Code:    user.ErrorCode_InvalidRequest,
			Message: "用户ID不能为空",
		}, nil
	}

	// 查找用户
	var userDO mysql.UserDO
	result := mysql.DB.First(&userDO, req.UserId)
	if result.Error != nil {
		// 用户不存在
		if result.Error == gorm.ErrRecordNotFound {
			return &user.UserInfoResponse{
				Code:    user.ErrorCode_UserNotFound,
				Message: "用户未找到",
			}, nil
		}
		// 其他数据库错误
		return &user.UserInfoResponse{
			Code:    user.ErrorCode_InternalError,
			Message: fmt.Sprintf("数据库查询失败: %v", result.Error),
		}, result.Error
	}

	// 返回用户信息
	return &user.UserInfoResponse{
		Code:    user.ErrorCode_Success,
		Message: "用户信息获取成功",
		User: &user.User{
			Id:        int64(userDO.ID),
			Username:  userDO.Username,
			Email:     userDO.Email,
			Gender:    userDO.Gender,
			Age:       int32(userDO.Age),
			Address:   userDO.Address,
			CreatedAt: userDO.CreatedAt.String(),
			UpdatedAt: userDO.UpdatedAt.String(),
		},
	}, nil
}
