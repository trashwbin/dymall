package service

import (
	"context"
	"fmt"

	"github.com/trashwbin/dymall/app/user/biz/dal/mysql" // 引入 mysql 包来操作数据库
	"github.com/trashwbin/dymall/app/user/biz/dal/redis"
	user "github.com/trashwbin/dymall/rpc_gen/kitex_gen/user"
)

type LogoutService struct {
	ctx context.Context
}

// NewLogoutService 创建新的 LogoutService 实例
func NewLogoutService(ctx context.Context) *LogoutService {
	return &LogoutService{ctx: ctx}
}

// Run 执行登出操作
func (s *LogoutService) Run(req *user.LogoutRequest) (resp *user.LogoutResponse, err error) {
	// 检查 user_id 是否提供
	if req.UserId == 0 {
		return &user.LogoutResponse{
			Code:    user.ErrorCode_InvalidRequest,
			Message: "用户ID不能为空",
		}, nil
	}

	// 查找用户是否存在
	var userDO mysql.UserDO
	result := mysql.DB.First(&userDO, req.UserId)
	if result.Error != nil {
		// 用户不存在
		return &user.LogoutResponse{
			Code:    user.ErrorCode_UserNotFound,
			Message: "用户未找到",
		}, nil
	}

	// 执行登出操作
	// 对于 Token 登出，通常是将 token 标记为无效，或者从缓存中移除 token
	// 如果使用的是 session，可以在这里删除与用户 session 相关的所有数据

	// 我们可以将 token从缓存中删除。
	// 生成 Redis Key
	redisKey := fmt.Sprintf("user:token:%d", req.UserId)

	// 从 Redis 删除 Token
	err = redis.RedisClient.Del(s.ctx, redisKey).Err()
	if err != nil {
		return &user.LogoutResponse{
			Code:    user.ErrorCode_InternalError,
			Message: "登出失败，无法清除 Token",
		}, err
	}

	// 此处假设操作成功
	return &user.LogoutResponse{
		Code:    user.ErrorCode_Success,
		Message: "用户登出成功",
	}, nil
}
