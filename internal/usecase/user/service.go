package user

import (
	"github.com/xdorro/proto-base-project/proto-gen-go/user/v1/userv1connect"
)

var _ IUserService = &Service{}

// IUserService user service interface.
type IUserService interface {
	userv1connect.UserServiceHandler
}

// Service struct.
type Service struct {
	// option

	userv1connect.UnimplementedUserServiceHandler
}

// Option service option.
type Option struct {
}

// NewService new service.
func NewService(opt *Option) IUserService {
	s := &Service{}

	return s
}
