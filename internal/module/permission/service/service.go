package permissionservice

import (
	"context"

	"github.com/bufbuild/connect-go"
	permissionv1 "github.com/xdorro/proto-base-project/proto-gen-go/permission/v1"
	"github.com/xdorro/proto-base-project/proto-gen-go/permission/v1/permissionv1connect"

	permissionbiz "github.com/xdorro/golang-grpc-base-project/internal/module/permission/biz"
)

var _ IPermissionService = &Service{}

// IPermissionService permission service interface.
type IPermissionService interface {
	permissionv1connect.PermissionServiceHandler
}

// Service struct.
type Service struct {
	// option
	permissionBiz permissionbiz.IPermissionBiz

	permissionv1connect.UnimplementedPermissionServiceHandler
}

// Option service option.
type Option struct {
	PermissionBiz permissionbiz.IPermissionBiz
}

// NewService new service.
func NewService(opt *Option) IPermissionService {
	s := &Service{
		permissionBiz: opt.PermissionBiz,
	}

	return s
}

// FindAllPermissions find all permissions
func (s *Service) FindAllPermissions(_ context.Context, req *connect.Request[permissionv1.FindAllPermissionsRequest]) (
	*connect.Response[permissionv1.FindAllPermissionsResponse], error,
) {
	return s.permissionBiz.FindAllPermissions(req)
}

// FindPermissionByID find permission by id
func (s *Service) FindPermissionByID(_ context.Context, req *connect.Request[permissionv1.CommonUUIDRequest]) (
	*connect.Response[permissionv1.Permission], error,
) {
	return s.permissionBiz.FindPermissionByID(req)
}

// CreatePermission create permission
func (s *Service) CreatePermission(_ context.Context, req *connect.Request[permissionv1.CreatePermissionRequest]) (
	*connect.Response[permissionv1.CommonResponse], error,
) {
	return s.permissionBiz.CreatePermission(req)
}

// UpdatePermission update permission by id
func (s *Service) UpdatePermission(_ context.Context, req *connect.Request[permissionv1.UpdatePermissionRequest]) (
	*connect.Response[permissionv1.CommonResponse], error,
) {
	return s.permissionBiz.UpdatePermission(req)
}

// DeletePermission delete permission by id
func (s *Service) DeletePermission(_ context.Context, req *connect.Request[permissionv1.CommonUUIDRequest]) (
	*connect.Response[permissionv1.CommonResponse], error,
) {
	return s.permissionBiz.DeletePermission(req)
}
