package userbiz

import (
	"github.com/xdorro/proto-base-project/proto-gen-go/user/v1/userv1connect"

	userrepo "github.com/xdorro/golang-grpc-base-project/internal/module/user/repo"
)

var _ IUserService = &Service{}

// IUserService user service interface.
type IUserService interface {
	userv1connect.UserServiceHandler
}

// Service struct.
type Service struct {
	// option
	userRepo userrepo.IRepo

	userv1connect.UnimplementedUserServiceHandler
}

// Option service option.
type Option struct {
	UserRepo userrepo.IRepo
}

// NewService new service.
func NewService(opt *Option) IUserService {
	s := &Service{
		userRepo: opt.UserRepo,
	}

	return s
}
