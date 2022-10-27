package service

import (
	"fmt"
	"net/http"
	"reflect"
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
	"github.com/xdorro/golang-grpc-base-project/pkg/redis"
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
	Redis       redis.IRedis

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
	redis       redis.IRedis

	mu       sync.Mutex
	services []string
	methods  []string
}

// NewService new service.
func NewService(opt *Option) IService {
	s := &Service{
		mux:   opt.Mux,
		repo:  opt.Repo,
		redis: opt.Redis,
	}

	// Add connect options
	connectOption := connect.WithOptions(
		connect.WithCompressMinBytes(1024),
		connect.WithInterceptors(opt.Interceptor.UnaryInterceptor()),
	)

	// Add your handlers here
	s.addServiceHandler(userv1connect.UnimplementedUserServiceHandler{},
		func() (string, http.Handler) {
			return userv1connect.NewUserServiceHandler(opt.UserService, connectOption)
		})

	s.addServiceHandler(authv1connect.UnimplementedAuthServiceHandler{},
		func() (string, http.Handler) {
			return authv1connect.NewAuthServiceHandler(opt.AuthService, connectOption)
		})

	s.addServiceHandler(permissionv1connect.UnimplementedPermissionServiceHandler{},
		func() (string, http.Handler) {
			return permissionv1connect.NewPermissionServiceHandler(opt.PermissionService, connectOption)
		})

	s.addServiceHandler(rolev1connect.UnimplementedRoleServiceHandler{},
		func() (string, http.Handler) {
			return rolev1connect.NewRoleServiceHandler(opt.RoleService, connectOption)
		})

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

	group.Go(func() error {
		return s.redis.Close()
	})

	return group.Wait()
}

// serviceHandler add the service handler.
func (s *Service) serviceHandler(opts connect.Option) {
	// Health check
	checker := grpchealth.NewStaticChecker(s.services...)
	s.mux.Handle(grpchealth.NewHandler(checker, opts))

	// Reflect serviceHandler
	reflector := grpcreflect.NewStaticReflector(s.services...)
	s.mux.Handle(grpcreflect.NewHandlerV1(reflector, opts))
	// Many tools still expect the older version of the server reflection API, so
	// most servers should mount both handlers.
	s.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector, opts))
}

// addServiceHandler adds a serviceHandler.
func (s *Service) addServiceHandler(svcMethod any, fn func() (string, http.Handler)) {
	str, handler := fn()

	logger := log.Info()

	s.mu.Lock()
	if svcMethod != nil {
		t := reflect.TypeOf(svcMethod)
		methods := make([]string, 0, t.NumMethod())

		for i := 0; i < t.NumMethod(); i++ {
			methods = append(methods, fmt.Sprintf("%s%s", str, t.Method(i).Name))
		}

		s.methods = append(s.methods, methods...)

		logger.Strs("methods", methods)
	}

	// add service name to list of services
	svcName := strings.TrimSpace(strings.ReplaceAll(str, "/", ""))
	s.services = append(s.services, svcName)
	s.mu.Unlock()

	// add serviceHandler
	s.mux.Handle(str, handler)

	logger.Msgf("Added service handler for %s", svcName)
}
