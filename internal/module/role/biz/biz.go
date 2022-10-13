package rolebiz

import (
	"fmt"
	"strings"

	"github.com/bufbuild/connect-go"
	rolev1 "github.com/xdorro/proto-base-project/proto-gen-go/role/v1"

	"github.com/xdorro/golang-grpc-base-project/pkg/casbin"
)

var _ IRoleBiz = &Biz{}

// IRoleBiz role service interface.
type IRoleBiz interface {
	FindAllRoles(req *connect.Request[rolev1.FindAllRolesRequest]) (
		*connect.Response[rolev1.FindAllRolesResponse], error,
	)
	FindRoleByName(req *connect.Request[rolev1.CommonNameRequest]) (
		*connect.Response[rolev1.Role], error,
	)
	CreateRole(req *connect.Request[rolev1.CreateRoleRequest]) (
		*connect.Response[rolev1.CommonResponse], error,
	)
	UpdateRole(req *connect.Request[rolev1.UpdateRoleRequest]) (
		*connect.Response[rolev1.CommonResponse], error,
	)
	DeleteRole(req *connect.Request[rolev1.CommonNameRequest]) (
		*connect.Response[rolev1.CommonResponse], error,
	)
}

// Biz struct.
type Biz struct {
	// option
	casbin casbin.ICasbin
}

// Option service option.
type Option struct {
	Casbin casbin.ICasbin
}

// NewBiz new service.
func NewBiz(opt *Option) IRoleBiz {
	b := &Biz{
		casbin: opt.Casbin,
	}

	return b
}

// FindAllRoles find all roles
func (b *Biz) FindAllRoles(req *connect.Request[rolev1.FindAllRolesRequest]) (
	*connect.Response[rolev1.FindAllRolesResponse], error,
) {
	data := make([]*rolev1.Role, 0)

	roles := b.casbin.Enforcer().GetAllSubjects()
	for _, role := range roles {
		data = append(data, &rolev1.Role{
			Name: role,
		})
	}

	res := &rolev1.FindAllRolesResponse{
		Data: data,
	}

	return connect.NewResponse(res), nil
}

// FindRoleByName find role by name
func (b *Biz) FindRoleByName(req *connect.Request[rolev1.CommonNameRequest]) (
	*connect.Response[rolev1.Role], error,
) {
	name := strings.ToLower(req.Msg.GetName())

	policies := b.casbin.Enforcer().GetFilteredPolicy(0, name)
	if len(policies) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("role does not exists"))
	}

	permissions := make([]string, 0)
	for _, policy := range policies {
		permissions = append(permissions, policy[1])
	}

	res := &rolev1.Role{
		Name:        name,
		Permissions: permissions,
	}

	return connect.NewResponse(res), nil
}

// CreateRole create role
func (b *Biz) CreateRole(req *connect.Request[rolev1.CreateRoleRequest]) (
	*connect.Response[rolev1.CommonResponse], error,
) {
	name := strings.ToLower(req.Msg.GetName())

	policies := b.casbin.Enforcer().GetFilteredPolicy(0, name)
	if len(policies) > 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("role already exists"))
	}

	for _, per := range req.Msg.GetPermissions() {
		policies = append(policies, []string{name, per})
	}

	// add policies to casbin
	_, err := b.casbin.Enforcer().AddPolicies(policies)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &rolev1.CommonResponse{
		Data: "success",
	}

	return connect.NewResponse(res), nil
}

// UpdateRole update role
func (b *Biz) UpdateRole(req *connect.Request[rolev1.UpdateRoleRequest]) (
	*connect.Response[rolev1.CommonResponse], error,
) {
	name := strings.ToLower(req.Msg.GetName())

	oldPolicies := b.casbin.Enforcer().GetFilteredPolicy(0, name)
	if len(oldPolicies) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("role does not exists"))
	}

	_, err := b.casbin.Enforcer().RemovePolicies(oldPolicies)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	policies := make([][]string, 0)
	for _, per := range req.Msg.GetPermissions() {
		policies = append(policies, []string{name, per})
	}

	// update policies to casbin
	_, err = b.casbin.Enforcer().AddPolicies(policies)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &rolev1.CommonResponse{
		Data: "success",
	}

	return connect.NewResponse(res), nil
}

// DeleteRole delete role
func (b *Biz) DeleteRole(req *connect.Request[rolev1.CommonNameRequest]) (
	*connect.Response[rolev1.CommonResponse], error,
) {
	name := strings.ToLower(req.Msg.GetName())

	policies := b.casbin.Enforcer().GetFilteredPolicy(0, name)
	if len(policies) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("role does not exists"))
	}

	_, err := b.casbin.Enforcer().RemovePolicies(policies)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &rolev1.CommonResponse{
		Data: "success",
	}

	return connect.NewResponse(res), nil
}
