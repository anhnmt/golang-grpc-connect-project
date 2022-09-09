package userservice

import (
	"context"

	"github.com/bufbuild/connect-go"
	userv1 "github.com/xdorro/proto-base-project/proto-gen-go/user/v1"
	"github.com/xdorro/proto-base-project/proto-gen-go/user/v1/userv1connect"

	userbiz "github.com/xdorro/golang-grpc-base-project/internal/module/user/biz"
)

var _ IUserService = &Service{}

// IUserService user service interface.
type IUserService interface {
	userv1connect.UserServiceHandler
}

// Service struct.
type Service struct {
	// option
	userBiz userbiz.IUserBiz

	userv1connect.UnimplementedUserServiceHandler
}

// Option service option.
type Option struct {
	UserBiz userbiz.IUserBiz
}

// NewService new service.
func NewService(opt *Option) IUserService {
	s := &Service{
		userBiz: opt.UserBiz,
	}

	return s
}

// FindAllUsers is the user.v1.UserService.FindAllUsers method.
func (s *Service) FindAllUsers(_ context.Context, req *connect.Request[userv1.FindAllUsersRequest]) (
	*connect.Response[userv1.FindAllUsersResponse], error,
) {
	return s.userBiz.FindAllUsers(req)
}

// FindUserByID is the user.v1.UserService.FindUserByID method.
func (s *Service) FindUserByID(_ context.Context, req *connect.Request[userv1.CommonUUIDRequest]) (
	*connect.Response[userv1.User], error,
) {
	return s.userBiz.FindUserByID(req)
}

// CreateUser is the user.v1.UserService.CreateUser method.
func (s *Service) CreateUser(_ context.Context, req *connect.Request[userv1.CreateUserRequest]) (
	*connect.Response[userv1.CommonResponse], error,
) {
	return s.userBiz.CreateUser(req)
}

// UpdateUser is the user.v1.UserService.UpdateUser method.
func (s *Service) UpdateUser(_ context.Context, req *connect.Request[userv1.UpdateUserRequest]) (
	*connect.Response[userv1.CommonResponse], error,
) {
	return s.userBiz.UpdateUser(req)
}

// DeleteUser is the user.v1.UserService.DeleteUser method.
func (s *Service) DeleteUser(_ context.Context, req *connect.Request[userv1.CommonUUIDRequest]) (
	*connect.Response[userv1.CommonResponse], error,
) {
	return s.userBiz.DeleteUser(req)
}
