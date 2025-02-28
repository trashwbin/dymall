package service

import (
	"context"
	"fmt"
	"github.com/trashwbin/dymall/app/user/biz/dal/mysql" // 引入 mysql 包来操作数据库
	user "github.com/trashwbin/dymall/rpc_gen/kitex_gen/user"
	"gorm.io/gorm"
)

type UpdateUserService struct {
	ctx context.Context
}

// NewUpdateUserService 创建新的 UpdateUserService 实例
func NewUpdateUserService(ctx context.Context) *UpdateUserService {
	return &UpdateUserService{ctx: ctx}
}

// Run 更新用户信息
func (s *UpdateUserService) Run(req *user.UpdateUserRequest) (resp *user.UpdateUserResponse, err error) {
	// 检查 user_id 是否提供
	if req.Id == 0 {
		return &user.UpdateUserResponse{
			Code:    user.ErrorCode_InvalidRequest,
			Message: "用户ID不能为空",
		}, nil
	}

	// 查找用户
	var userDO mysql.UserDO
	result := mysql.DB.First(&userDO, req.Id)
	if result.Error != nil {
		// 用户不存在
		if result.Error == gorm.ErrRecordNotFound {
			return &user.UpdateUserResponse{
				Code:    user.ErrorCode_UserNotFound,
				Message: "用户未找到",
			}, nil
		}
		// 其他数据库错误
		return &user.UpdateUserResponse{
			Code:    user.ErrorCode_InternalError,
			Message: fmt.Sprintf("数据库查询失败: %v", result.Error),
		}, result.Error
	}

	// 更新用户信息
	if req.Username != "" {
		userDO.Username = req.Username
	}
	if req.Password != "" {
		userDO.Password = req.Password
	}
	if req.Email != "" {
		userDO.Email = req.Email
	}
	if req.Gender != "" {
		userDO.Gender = req.Gender
	}
	if req.Age != 0 {
		userDO.Age = int(req.Age)
	}
	if req.Address != "" {
		userDO.Address = req.Address
	}

	// 保存更新后的用户信息
	updateResult := mysql.DB.Save(&userDO)
	if updateResult.Error != nil {
		return &user.UpdateUserResponse{
			Code:    user.ErrorCode_InternalError,
			Message: fmt.Sprintf("更新用户失败: %v", updateResult.Error),
		}, updateResult.Error
	}

	// 返回更新后的用户信息
	return &user.UpdateUserResponse{
		Code:    user.ErrorCode_Success,
		Message: "用户信息更新成功",
	}, nil
}
