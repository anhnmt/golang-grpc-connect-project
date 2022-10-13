package service

import (
	"net/http"
	"strings"
	"sync"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/rs/zerolog/log"
	"github.com/xdorro/proto-base-project/proto-gen-go/auth/v1/authv1connect"
	"github.com/xdorro/proto-base-project/proto-gen-go/permission/v1/permissionv1connect"
	"github.com/xdorro/proto-base-project/proto-gen-go/role/v1/rolev1connect"
	"github.com/xdorro/proto-base-project/proto-gen-go/user/v1/userv1connect"
	"golang.org/x/sync/errgroup"

	"github.com/xdorro/golang-grpc-base-project/internal/interceptor"
	authservice "github.com/xdorro/golang-grpc-base-project/internal/module/auth/service"
	permissionservice "github.com/xdorro/golang-grpc-base-project/internal/module/permission/service"
	roleservice "github.com/xdorro/golang-grpc-base-project/internal/module/role/service"
	userservice "github.com/xdorro/golang-grpc-base-project/internal/module/user/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

var _ IService = &Service{}

// IService service interface.
type IService interface {
	Close() error
}

// Option service option.
type Option struct {
	Mux         *http.ServeMux
	Interceptor interceptor.IInterceptor
	Repo        repo.IRepo

	UserService       userservice.IUserService
	AuthService       authservice.IAuthService
	PermissionService permissionservice.IPermissionService
	RoleService       roleservice.IRoleService
}

// Service struct.
type Service struct {
	// options
	mux         *http.ServeMux
	interceptor interceptor.IInterceptor
	repo        repo.IRepo

	mu       sync.Mutex
	services []string
}

// NewService new service.
func NewService(opt *Option) IService {
	s := &Service{
		mux:  opt.Mux,
		repo: opt.Repo,
	}

	// Add connect options
	connectOption := connect.WithOptions(
		connect.WithCompressMinBytes(1024),
		connect.WithInterceptors(opt.Interceptor.UnaryInterceptor()),
	)

	// Add your handlers here
	s.addHandler(userv1connect.NewUserServiceHandler(opt.UserService, connectOption))
	s.addHandler(authv1connect.NewAuthServiceHandler(opt.AuthService, connectOption))
	s.addHandler(permissionv1connect.NewPermissionServiceHandler(opt.PermissionService, connectOption))
	s.addHandler(rolev1connect.NewRoleServiceHandler(opt.RoleService, connectOption))

	// Add service handlers
	s.serviceHandler(connectOption)

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

// serviceHandler add the service handler.
func (s *Service) serviceHandler(opts connect.Option) {
	// Health check
	checker := grpchealth.NewStaticChecker(s.services...)
	s.addHandler(grpchealth.NewHandler(checker, opts))

	// Reflect serviceHandler
	reflector := grpcreflect.NewStaticReflector(s.services...)
	s.addHandler(grpcreflect.NewHandlerV1(reflector, opts))
	// Many tools still expect the older version of the server reflection API, so
	// most servers should mount both handlers.
	s.addHandler(grpcreflect.NewHandlerV1Alpha(reflector, opts))
}

// addHandler adds a serviceHandler.
func (s *Service) addHandler(str string, handler http.Handler) {
	s.mu.Lock()
	// add service name to list of services
	svcName := strings.TrimSpace(strings.ReplaceAll(str, "/", ""))
	s.services = append(s.services, svcName)
	s.mu.Unlock()

	// add serviceHandler
	s.mux.Handle(str, handler)

	log.Info().Msgf("Added serviceHandler for %s", svcName)
}
