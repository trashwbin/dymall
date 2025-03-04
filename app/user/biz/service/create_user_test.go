package service

import (
	"context"
	"testing"

	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/user"
)

func TestCreateUser_Run(t *testing.T) {
	// ctx := context.Background()
	// s := NewCreateUserService(ctx) // 初始化 CreateUserService

	// // 创建一个新用户请求
	// req := &user.CreateUserRequest{
	// 	Username: "testuser",
	// 	Password: "password123", // 密码需要加密
	// 	Email:    "testuser@example.com",
	// 	Gender:   "male",
	// 	Age:      30,
	// 	Address:  "1234 Elm Street",
	// }

	// // 调用 Run 方法，执行创建用户操作
	// resp, err := s.Run(req)

	// // 记录日志
	// t.Logf("err: %v", err)
	// t.Logf("resp: %v", resp)

	// // 使用断言来验证响应
	// assert.NoError(t, err)                             // 确保没有错误
	// assert.Equal(t, user.ErrorCode_Success, resp.Code) // 验证是否成功
	// assert.NotEmpty(t, resp.Id)                        // 验证是否返回了用户ID

	// // 可选：进一步验证返回的响应消息
	// assert.Equal(t, "用户创建成功", resp.Message)
	ctx := context.Background()
	s := NewCreateUserService(ctx)
	// init req and assert value

	req := &user.CreateUserRequest{
		Username: "testuser",
		Password: "password123", // 密码需要加密
		Email:    "testuser@example.com",
		Gender:   "male",
		Age:      30,
		Address:  "1234 Elm Street",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
