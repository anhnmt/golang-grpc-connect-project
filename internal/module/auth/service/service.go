package authservice

import (
	"context"

	"github.com/bufbuild/connect-go"
	authv1 "github.com/xdorro/proto-base-project/proto-gen-go/auth/v1"
	"github.com/xdorro/proto-base-project/proto-gen-go/auth/v1/authv1connect"

	authbiz "github.com/xdorro/golang-grpc-base-project/internal/module/auth/biz"
)

var _ IAuthService = &Service{}

// IAuthService auth service interface.
type IAuthService interface {
	authv1connect.AuthServiceHandler
}

// Service struct.
type Service struct {
	// option
	authBiz authbiz.IAuthBiz

	authv1connect.UnimplementedAuthServiceHandler
}

// Option service option.
type Option struct {
	AuthBiz authbiz.IAuthBiz
}

// NewService new service.
func NewService(opt *Option) IAuthService {
	s := &Service{
		authBiz: opt.AuthBiz,
	}

	return s
}

// Login is the auth.v1.AuthService.Login method.
func (s *Service) Login(_ context.Context, req *connect.Request[authv1.LoginRequest]) (
	*connect.Response[authv1.TokenResponse], error,
) {
	return s.authBiz.Login(req)
}

// RevokeToken is the auth.v1.AuthService.RevokeToken method.
func (s *Service) RevokeToken(_ context.Context, req *connect.Request[authv1.TokenRequest]) (
	*connect.Response[authv1.CommonResponse], error,
) {
	return s.authBiz.RevokeToken(req)
}

// RefreshToken is the auth.v1.AuthService.RefreshToken method.
func (s *Service) RefreshToken(_ context.Context, req *connect.Request[authv1.TokenRequest]) (
	*connect.Response[authv1.TokenResponse], error,
) {
	return s.authBiz.RefreshToken(req)
}
