package main

import (
	"context"
	"github.com/trashwbin/dymall/app/user/biz/service"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// CreateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	resp, err = service.NewCreateUserService(ctx).Run(req)

	return resp, err
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	resp, err = service.NewLoginService(ctx).Run(req)

	return resp, err
}

// Logout implements the UserServiceImpl interface.
func (s *UserServiceImpl) Logout(ctx context.Context, req *user.LogoutRequest) (resp *user.LogoutResponse, err error) {
	resp, err = service.NewLogoutService(ctx).Run(req)

	return resp, err
}

// UpdateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (resp *user.UpdateUserResponse, err error) {
	resp, err = service.NewUpdateUserService(ctx).Run(req)

	return resp, err
}

// DeleteUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *user.UserInfoRequest) (resp *user.DeleteUserResponse, err error) {
	resp, err = service.NewDeleteUserService(ctx).Run(req)

	return resp, err
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	resp, err = service.NewGetUserInfoService(ctx).Run(req)

	return resp, err
}
