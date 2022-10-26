package service

import (
	"context"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	userv1 "github.com/xdorro/proto-base-project/proto-gen-go/user/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	userservice "github.com/xdorro/golang-grpc-base-project/internal/module/user/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

var _ IService = &Service{}

// IService service interface.
type IService interface {
	Close() error
	RegisterGrpcServerHandler(grpcServer *grpc.Server)
	RegisterHttpServerHandler(httpServer *runtime.ServeMux)
}

// Option service option.
type Option struct {
	// Interceptor interceptor.IInterceptor
	Repo repo.IRepo

	UserService userservice.IUserService
}

// Service struct.
type Service struct {
	// options
	userService userservice.IUserService

	// interceptor interceptor.IInterceptor
	repo repo.IRepo

	mu       sync.Mutex
	services []string
}

// NewService new service.
func NewService(opt *Option) IService {
	s := &Service{
		repo:        opt.Repo,
		userService: opt.UserService,
	}

	return s
}

// Close the Service.
func (s *Service) Close() error {
	group := new(errgroup.Group)

	group.Go(func() error {
		return s.repo.Close()
	})

	return group.Wait()
}

// RegisterGrpcServerHandler adds a serviceHandler.
func (s *Service) RegisterGrpcServerHandler(grpcServer *grpc.Server) {
	userv1.RegisterUserServiceServer(grpcServer, s.userService)
}

// RegisterHttpServerHandler adds a serviceHandler.
func (s *Service) RegisterHttpServerHandler(httpServer *runtime.ServeMux) {
	ctx := context.Background()
	if err := userv1.RegisterUserServiceHandlerServer(ctx, httpServer, s.userService); err != nil {
		log.Panic().Err(err).Msg("RegisterUserServiceHandler failed")
	}
}
