package service

import (
	"context"
	"fmt"
	"github.com/trashwbin/dymall/app/user/biz/dal/mysql" // 引入 mysql 包来操作数据库
	user "github.com/trashwbin/dymall/rpc_gen/kitex_gen/user"
	"gorm.io/gorm"
)

type DeleteUserService struct {
	ctx context.Context
}

// NewDeleteUserService 创建新的 DeleteUserService 实例
func NewDeleteUserService(ctx context.Context) *DeleteUserService {
	return &DeleteUserService{ctx: ctx}
}

// Run 执行删除用户操作
func (s *DeleteUserService) Run(req *user.UserInfoRequest) (resp *user.DeleteUserResponse, err error) {
	// 检查 user_id 是否提供
	if req.UserId == 0 {
		return &user.DeleteUserResponse{
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
			return &user.DeleteUserResponse{
				Code:    user.ErrorCode_UserNotFound,
				Message: "用户未找到",
			}, nil
		}
		// 其他数据库错误
		return &user.DeleteUserResponse{
			Code:    user.ErrorCode_InternalError,
			Message: fmt.Sprintf("数据库查询失败: %v", result.Error),
		}, result.Error
	}

	// 执行软删除
	result = mysql.DB.Delete(&userDO)
	if result.Error != nil {
		return &user.DeleteUserResponse{
			Code:    user.ErrorCode_InternalError,
			Message: fmt.Sprintf("删除用户失败: %v", result.Error),
		}, result.Error
	}

	// 成功删除
	return &user.DeleteUserResponse{
		Code:    user.ErrorCode_Success,
		Message: "用户删除成功",
	}, nil
}
