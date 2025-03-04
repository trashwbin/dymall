package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trashwbin/dymall/app/user/biz/dal/mysql"
	user "github.com/trashwbin/dymall/rpc_gen/kitex_gen/user"
)

// MockUserRepo 是 UserRepo 的 mock 实现，用于模拟数据库操作。
type MockUserRepo struct {
	GetUserByUsernameFunc func(string) (*user.User, error)
	CreateUserFunc        func(*user.User) error
}

// GetUserByUsername 实现了 UserRepo 接口的 GetUserByUsername 方法。
func (m *MockUserRepo) GetUserByUsername(username string) (*user.User, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(username)
	}
	return nil, errors.New("GetUserByUsername not implemented")
}

// CreateUser 实现了 UserRepo 接口的 CreateUser 方法。
func (m *MockUserRepo) CreateUser(user *user.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return errors.New("CreateUser not implemented")
}

func TestCreateUserService_Run(t *testing.T) {
	ctx := context.Background()

	// 定义测试用例集合
	tests := []struct {
		name            string                   // 测试用例名称
		req             *user.CreateUserRequest  // 创建用户的请求参数
		mockSetup       func(repo *MockUserRepo) // 设置 Mock 行为的函数
		expectedCode    user.ErrorCode           // 预期的错误码
		expectedMessage string                   // 预期的错误消息
	}{
		{
			name: "成功创建用户",
			req: &user.CreateUserRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "testuser@example.com",
				Gender:   "male",
				Age:      30,
				Address:  "1234 Elm Street",
			},
			mockSetup: func(repo *MockUserRepo) {
				// 模拟用户名不存在的情况
				repo.GetUserByUsernameFunc = func(username string) (*user.User, error) {
					return nil, nil
				}
				// 模拟用户创建成功，并返回一个模拟的用户 ID
				repo.CreateUserFunc = func(user *user.User) error {
					user.Id = 1 // 模拟生成的用户 ID
					return nil
				}
			},
			expectedCode:    user.ErrorCode_Success,
			expectedMessage: "用户创建成功",
		},
		{
			name: "用户名已存在",
			req: &user.CreateUserRequest{
				Username: "existinguser",
				Password: "password123",
				Email:    "existinguser@example.com",
				Gender:   "male",
				Age:      30,
				Address:  "1234 Elm Street",
			},
			mockSetup: func(repo *MockUserRepo) {
				// 模拟用户名已存在的场景
				repo.GetUserByUsernameFunc = func(username string) (*user.User, error) {
					return &user.User{Username: "existinguser"}, nil
				}
			},
			expectedCode:    user.ErrorCode_InvalidRequest,
			expectedMessage: "用户名已存在",
		},
		{
			name: "数据库查询失败",
			req: &user.CreateUserRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "testuser@example.com",
				Gender:   "male",
				Age:      30,
				Address:  "1234 Elm Street",
			},
			mockSetup: func(repo *MockUserRepo) {
				// 模拟数据库查询失败的场景
				repo.GetUserByUsernameFunc = func(username string) (*user.User, error) {
					return nil, errors.New("数据库查询失败")
				}
			},
			expectedCode:    user.ErrorCode_InternalError,
			expectedMessage: "数据库查询失败",
		},
		{
			name: "创建用户失败",
			req: &user.CreateUserRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "testuser@example.com",
				Gender:   "male",
				Age:      30,
				Address:  "1234 Elm Street",
			},
			mockSetup: func(repo *MockUserRepo) {
				// 模拟用户名不存在的情况
				repo.GetUserByUsernameFunc = func(username string) (*user.User, error) {
					return nil, nil
				}
				// 模拟用户创建失败的场景
				repo.CreateUserFunc = func(user *user.User) error {
					return errors.New("创建用户失败")
				}
			},
			expectedCode:    user.ErrorCode_InternalError,
			expectedMessage: "创建用户失败",
		},
	}

	// 遍历每个测试用例并执行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepo{}
			tt.mockSetup(mockRepo)

			s := &CreateUserService{
				ctx:      ctx,
				userRepo: mockRepo, // 使用 mock 的 UserRepo
			}

			resp, err := s.Run(tt.req)

			// 记录日志以方便调试
			t.Logf("err: %v", err)
			t.Logf("resp: %+v", resp)

			// 使用断言验证响应是否符合预期
			assert.Equal(t, tt.expectedCode, resp.Code)
			assert.Equal(t, tt.expectedMessage, resp.Message)

			// 对于成功的测试用例，验证返回的用户 ID 是否不为空
			if tt.expectedCode == user.ErrorCode_Success {
				assert.NotEmpty(t, resp.Id)
			} else {
				// 对于失败的测试用例，验证返回的用户 ID 是否为空
				assert.Empty(t, resp.Id)
			}
		})
	}
