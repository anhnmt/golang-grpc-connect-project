package roleservice

import (
	"context"

	"github.com/bufbuild/connect-go"
	rolev1 "github.com/xdorro/proto-base-project/proto-gen-go/role/v1"
	"github.com/xdorro/proto-base-project/proto-gen-go/role/v1/rolev1connect"

	rolebiz "github.com/xdorro/golang-grpc-base-project/internal/module/role/biz"
)

var _ IRoleService = &Service{}

// IRoleService role service interface.
type IRoleService interface {
	rolev1connect.RoleServiceHandler
}

// Service struct.
type Service struct {
	// option
	roleBiz rolebiz.IRoleBiz

	rolev1connect.UnimplementedRoleServiceHandler
}

// Option service option.
type Option struct {
	RoleBiz rolebiz.IRoleBiz
}

// NewService new service.
func NewService(opt *Option) IRoleService {
	s := &Service{
		roleBiz: opt.RoleBiz,
	}

	return s
}

// FindAllRoles is the role.v1.RoleService.FindAllRoles method.
func (s *Service) FindAllRoles(_ context.Context, req *connect.Request[rolev1.FindAllRolesRequest]) (
	*connect.Response[rolev1.FindAllRolesResponse], error,
) {
	return s.roleBiz.FindAllRoles(req)
}

// FindRoleByName is the role.v1.RoleService.FindRoleByName method.
func (s *Service) FindRoleByName(_ context.Context, req *connect.Request[rolev1.CommonNameRequest]) (
	*connect.Response[rolev1.Role], error,
) {
	return s.roleBiz.FindRoleByName(req)
}

// CreateRole is the role.v1.RoleService.CreateRole method.
func (s *Service) CreateRole(_ context.Context, req *connect.Request[rolev1.CreateRoleRequest]) (
	*connect.Response[rolev1.CommonResponse], error,
) {
	return s.roleBiz.CreateRole(req)
}

// UpdateRole is the role.v1.RoleService.UpdateRole method.
func (s *Service) UpdateRole(_ context.Context, req *connect.Request[rolev1.UpdateRoleRequest]) (
	*connect.Response[rolev1.CommonResponse], error,
) {
	return s.roleBiz.UpdateRole(req)
}

// DeleteRole is the role.v1.RoleService.DeleteRole method.
func (s *Service) DeleteRole(_ context.Context, req *connect.Request[rolev1.CommonNameRequest]) (
	*connect.Response[rolev1.CommonResponse], error,
) {
	return s.roleBiz.DeleteRole(req)
}
