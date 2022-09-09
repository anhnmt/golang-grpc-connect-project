package authbiz

import (
	"github.com/xdorro/proto-base-project/proto-gen-go/auth/v1/authv1connect"
)

var _ IAuthService = &Service{}

// IAuthService auth service interface.
type IAuthService interface {
	authv1connect.AuthServiceHandler
}

// Service struct.
type Service struct {
	// option

	authv1connect.UnimplementedAuthServiceHandler
}

// Option service option.
type Option struct {
}

// NewService new service.
func NewService(*Option) IAuthService {
	s := &Service{}

	return s
}
