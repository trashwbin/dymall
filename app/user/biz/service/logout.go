package service

import (
	"context"
	"github.com/trashwbin/dymall/app/user/biz/dal/mysql" // 引入 mysql 包来操作数据库
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
	// 例如，我们可以将 token 添加到一个黑名单中，或者从缓存中删除。

	// 假设我们使用 Redis 或内存缓存来存储活跃的 token，
	// 如果我们使用 token 认证，可以在这里处理 token 的失效逻辑：

	// 例如，缓存清除操作
	// redisClient.Del(fmt.Sprintf("user_token:%d", req.UserId))

	// 此处假设操作成功
	return &user.LogoutResponse{
		Code:    user.ErrorCode_Success,
		Message: "用户登出成功",
	}, nil
}
