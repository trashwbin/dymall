package service

import (
	"context"
	"fmt"
	"github.com/trashwbin/dymall/app/cart/biz/dal/redis"
	"github.com/trashwbin/dymall/app/user/biz/dal/mysql" // 引入 mysql 包来操作数据库
	"github.com/trashwbin/dymall/app/user/infra/rpc"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
	user "github.com/trashwbin/dymall/rpc_gen/kitex_gen/user"
	"gorm.io/gorm"
)

type LoginService struct {
	ctx context.Context
}

// NewLoginService 创建新的 LoginService 实例
func NewLoginService(ctx context.Context) *LoginService {
	return &LoginService{ctx: ctx}
}

// Run 执行登录操作
func (s *LoginService) Run(req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	// 检查用户名和密码是否提供
	if req.Username == "" || req.Password == "" {
		return &user.LoginResponse{
			Code:    user.ErrorCode_InvalidRequest,
			Message: "用户名或密码不能为空",
		}, nil
	}

	// 查找用户
	var userDO mysql.UserDO
	result := mysql.DB.Where("username = ?", req.Username).First(&userDO)
	if result.Error != nil {
		// 用户不存在
		if result.Error == gorm.ErrRecordNotFound {
			return &user.LoginResponse{
				Code:    user.ErrorCode_UserNotFound,
				Message: "用户名或密码错误",
			}, nil
		}
		// 其他数据库错误
		return &user.LoginResponse{
			Code:    user.ErrorCode_InternalError,
			Message: fmt.Sprintf("数据库查询失败: %v", result.Error),
		}, result.Error
	}

	// 验证密码是否正确
	if userDO.Password != req.Password {
		return &user.LoginResponse{
			Code:    user.ErrorCode_InvalidCredentials,
			Message: "用户名或密码错误",
		}, nil
	}

	// 生成JWT Token
	deliverTokenRPC, err := rpc.AuthClient.DeliverTokenByRPC(s.ctx, &auth.DeliverTokenReq{})
	token := deliverTokenRPC.Token
	//这里使用redis存储token
	//TODO测试
	redis.RedisClient.Set(s.ctx, token, userDO.ID, 0)

	if err != nil {
		return &user.LoginResponse{
			Code:    user.ErrorCode_InternalError,
			Message: "生成Token失败",
		}, err
	}

	// 返回登录成功信息
	return &user.LoginResponse{
		Code:    user.ErrorCode_Success,
		Message: "登录成功",
		Token:   string(token),
	}, nil
}
