package permissionbiz

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog/log"
	permissionv1 "github.com/xdorro/proto-base-project/proto-gen-go/permission/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	permissionmodel "github.com/xdorro/golang-grpc-base-project/internal/module/permission/model"
	permissionrepo "github.com/xdorro/golang-grpc-base-project/internal/module/permission/repo"
	"github.com/xdorro/golang-grpc-base-project/pkg/utils"
)

var _ IPermissionBiz = &Biz{}

// IPermissionBiz permission service interface.
type IPermissionBiz interface {
	FindAllPermissions(req *connect.Request[permissionv1.FindAllPermissionsRequest]) (
		*connect.Response[permissionv1.FindAllPermissionsResponse], error,
	)
	FindPermissionByID(req *connect.Request[permissionv1.CommonUUIDRequest]) (
		*connect.Response[permissionv1.Permission], error,
	)
	CreatePermission(req *connect.Request[permissionv1.CreatePermissionRequest]) (
		*connect.Response[permissionv1.CommonResponse], error,
	)
	UpdatePermission(req *connect.Request[permissionv1.UpdatePermissionRequest]) (
		*connect.Response[permissionv1.CommonResponse], error,
	)
	DeletePermission(req *connect.Request[permissionv1.CommonUUIDRequest]) (
		*connect.Response[permissionv1.CommonResponse], error,
	)
}

// Biz struct.
type Biz struct {
	// option
	permissionRepo permissionrepo.IRepo
}

// Option service option.
type Option struct {
	PermissionRepo permissionrepo.IRepo
}

// NewBiz new service.
func NewBiz(opt *Option) IPermissionBiz {
	s := &Biz{
		permissionRepo: opt.PermissionRepo,
	}

	return s
}

// FindAllPermissions is the permission.v1.PermissionBiz.FindAllPermissions method.
func (s *Biz) FindAllPermissions(req *connect.Request[permissionv1.FindAllPermissionsRequest]) (
	*connect.Response[permissionv1.FindAllPermissionsResponse], error,
) {
	// count all permissions with filter
	filter := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	count, _ := s.permissionRepo.CountDocuments(filter)
	limit := int64(10)
	totalPages := utils.TotalPage(count, limit)
	page := utils.CurrentPage(req.Msg.GetPage(), totalPages)

	// find all permissions with filter and option
	opt := options.
		Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(limit).
		SetSkip((page - 1) * limit)
	data, err := s.permissionRepo.Find(filter, opt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &permissionv1.FindAllPermissionsResponse{
		TotalPage:   totalPages,
		CurrentPage: page,
		Data:        permissionmodel.PermissionsToProto(data),
	}

	return connect.NewResponse(res), nil
}

// FindPermissionByID is the permission.v1.PermissionBiz.FindPermissionByID method.
func (s *Biz) FindPermissionByID(req *connect.Request[permissionv1.CommonUUIDRequest]) (
	*connect.Response[permissionv1.Permission], error,
) {
	id := req.Msg.GetId()
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Err(err).Msg("Failed find permission by id")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	opt := options.
		FindOne().
		SetSort(bson.M{"created_at": -1})
	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	data, err := s.permissionRepo.FindOne(filter, opt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := permissionmodel.PermissionToProto(data)
	return connect.NewResponse(res), nil
}

// CreatePermission is the permission.v1.PermissionBiz.CreatePermission method.
func (s *Biz) CreatePermission(req *connect.Request[permissionv1.CreatePermissionRequest]) (
	*connect.Response[permissionv1.CommonResponse], error,
) {
	// count all permissions with filter
	countFilter := bson.M{
		"slug": req.Msg.GetSlug(),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	count, _ := s.permissionRepo.CountDocuments(countFilter)
	if count > 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("slug already exists"))
	}

	data := &permissionmodel.Permission{
		Name:        req.Msg.GetName(),
		Slug:        req.Msg.GetSlug(),
		RequireAuth: req.Msg.GetRequireAuth(),
		RequireHash: req.Msg.GetRequireHash(),
	}
	data.PreCreate()

	oid, err := s.permissionRepo.InsertOne(data)
	if err != nil {
		log.Err(err).Msg("Error create permission")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	resID := oid.InsertedID.(string)
	res := &permissionv1.CommonResponse{
		Data: resID,
	}

	// if err = redis.Del(s.redis, utils.ListAuthPermissionsKey); err != nil {
	// 	return nil, err
	// }

	return connect.NewResponse(res), nil
}

// UpdatePermission is the permission.v1.PermissionBiz.UpdatePermission method.
func (s *Biz) UpdatePermission(req *connect.Request[permissionv1.UpdatePermissionRequest]) (
	*connect.Response[permissionv1.CommonResponse], error,
) {
	id := req.Msg.GetId()
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Err(err).Msg("Failed find permission by id")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	data, err := s.permissionRepo.FindOne(filter)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// count all permissions with filter
	countFilter := bson.M{
		"_id":  bson.M{"$ne": id},
		"slug": req.Msg.GetSlug(),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	count, _ := s.permissionRepo.CountDocuments(countFilter)
	if count > 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("slug already exists"))
	}

	data.Name = utils.StringCompareOrPassValue(data.Name, req.Msg.GetName())
	data.Slug = utils.StringCompareOrPassValue(data.Slug, req.Msg.GetSlug())

	if req.Msg.RequireAuth != nil {
		data.RequireAuth = req.Msg.GetRequireAuth()
	}

	if req.Msg.RequireHash != nil {
		data.RequireHash = req.Msg.GetRequireHash()
	}

	data.PreUpdate()

	opt := bson.M{"$set": data}
	if _, err = s.permissionRepo.UpdateOne(filter, opt); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &permissionv1.CommonResponse{
		Data: req.Msg.GetId(),
	}

	// if err = redis.Del(s.redis, utils.ListAuthPermissionsKey); err != nil {
	// 	return nil, err
	// }

	return connect.NewResponse(res), nil
}

// DeletePermission is the permission.v1.PermissionBiz.DeletePermission method.
func (s *Biz) DeletePermission(req *connect.Request[permissionv1.CommonUUIDRequest]) (
	*connect.Response[permissionv1.CommonResponse], error,
) {
	id := req.Msg.GetId()
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Err(err).Msg("Failed find permission by id")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	// count all permissions with filter
	count, _ := s.permissionRepo.CountDocuments(filter)
	if count <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("permission does not exists"))
	}

	if _, err = s.permissionRepo.SoftDeleteOne(filter); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &permissionv1.CommonResponse{
		Data: req.Msg.GetId(),
	}

	// if err = redis.Del(s.redis, utils.ListAuthPermissionsKey); err != nil {
	// 	return nil, err
	// }

	return connect.NewResponse(res), nil
}
