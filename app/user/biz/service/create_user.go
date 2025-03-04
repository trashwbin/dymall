package service

import (
	"context"
	"fmt"
	"github.com/trashwbin/dymall/app/user/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/user/biz/util"
	user "github.com/trashwbin/dymall/rpc_gen/kitex_gen/user"
)

type CreateUserService struct {
	ctx      context.Context
	userRepo *mysql.UserRepo
} // NewCreateUserService new CreateUserService

func NewCreateUserService(ctx context.Context) *CreateUserService {
	return &CreateUserService{
		ctx:      ctx,
		userRepo: mysql.NewUserRepo(),
	}
}

// Run create note info
func (s *CreateUserService) Run(req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	// Finish your business logic.
	// 1. 检查用户名是否已存在
	existingUser, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return &user.CreateUserResponse{
			Code:    user.ErrorCode_InternalError,
			Message: "数据库查询失败",
		}, err
	}
	if existingUser != nil {
		return &user.CreateUserResponse{
			Code:    user.ErrorCode_InvalidRequest,
			Message: "用户名已存在",
		}, nil
	}

	// 2. 加密密码
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return &user.CreateUserResponse{
			Code:    user.ErrorCode_InternalError,
			Message: "密码加密失败",
		}, err
	}

	// 3. 创建用户
	newUser := &user.User{
		Username: req.Username,
		Password: hashedPassword, // 使用加密后的密码
		Email:    req.Email,
		Gender:   req.Gender,
		Age:      req.Age,
		Address:  req.Address,
	}

	// 4. 将新用户插入数据库
	err = s.userRepo.CreateUser(newUser)
	if err != nil {
		return &user.CreateUserResponse{
			Code:    user.ErrorCode_InternalError,
			Message: fmt.Sprintf("创建用户失败: %v", err),
		}, err
	}

	// 5. 返回成功响应
	return &user.CreateUserResponse{
		Code:    user.ErrorCode_Success,
		Message: "用户创建成功",
		Id:      newUser.Id, // 你可以返回数据库中分配的用户 ID
	}, nil
}
