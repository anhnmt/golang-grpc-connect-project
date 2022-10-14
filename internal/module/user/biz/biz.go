package userbiz

import (
	"fmt"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog/log"
	userv1 "github.com/xdorro/proto-base-project/proto-gen-go/user/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	usermodel "github.com/xdorro/golang-grpc-base-project/internal/module/user/model"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
	"github.com/xdorro/golang-grpc-base-project/pkg/utils"
)

var _ IUserBiz = &Biz{}

// IUserBiz user service interface.
type IUserBiz interface {
	FindAllUsers(req *connect.Request[userv1.FindAllUsersRequest]) (
		*connect.Response[userv1.FindAllUsersResponse], error,
	)
	FindUserByID(req *connect.Request[userv1.CommonUUIDRequest]) (*connect.Response[userv1.User], error)
	CreateUser(req *connect.Request[userv1.CreateUserRequest]) (*connect.Response[userv1.CommonResponse], error)
	UpdateUser(req *connect.Request[userv1.UpdateUserRequest]) (*connect.Response[userv1.CommonResponse], error)
	DeleteUser(req *connect.Request[userv1.CommonUUIDRequest]) (*connect.Response[userv1.CommonResponse], error)
}

// Biz struct.
type Biz struct {
	// option
	userCollection *mongo.Collection
}

// Option service option.
type Option struct {
	Repo repo.IRepo
}

// NewBiz new service.
func NewBiz(opt *Option) IUserBiz {
	s := &Biz{
		userCollection: opt.Repo.CollectionModel(&usermodel.User{}),
	}

	return s
}

// FindAllUsers is the user.v1.UserBiz.FindAllUsers method.
func (s *Biz) FindAllUsers(req *connect.Request[userv1.FindAllUsersRequest]) (
	*connect.Response[userv1.FindAllUsersResponse], error,
) {
	// count all users with filter
	filter := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	count, _ := repo.CountDocuments(s.userCollection, filter)
	limit := int64(10)
	totalPages := utils.TotalPage(count, limit)
	page := utils.CurrentPage(req.Msg.GetPage(), totalPages)

	// find all genres with filter and option
	opt := options.
		Find().
		SetSort(bson.M{"created_at": -1}).
		SetProjection(bson.M{"password": 0}).
		SetLimit(limit).
		SetSkip((page - 1) * limit)

	data, err := repo.Find[usermodel.User](s.userCollection, filter, opt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &userv1.FindAllUsersResponse{
		TotalPage:   totalPages,
		CurrentPage: page,
		Data:        usermodel.UsersToProto(data),
	}

	return connect.NewResponse(res), nil
}

// FindUserByID is the user.v1.UserBiz.FindUserByID method.
func (s *Biz) FindUserByID(req *connect.Request[userv1.CommonUUIDRequest]) (
	*connect.Response[userv1.User], error,
) {
	id, err := primitive.ObjectIDFromHex(req.Msg.GetId())
	if err != nil {
		log.Err(err).Msg("Failed find user by id")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	opt := options.
		FindOne().
		SetSort(bson.M{"created_at": -1}).
		SetProjection(bson.M{"password": 0})

	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	data, err := repo.FindOne[usermodel.User](s.userCollection, filter, opt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := usermodel.UserToProto(data)
	return connect.NewResponse(res), nil
}

// CreateUser is the user.v1.UserBiz.CreateUser method.
func (s *Biz) CreateUser(req *connect.Request[userv1.CreateUserRequest]) (
	*connect.Response[userv1.CommonResponse], error,
) {
	// count all users with filter
	count, _ := repo.CountDocuments(s.userCollection, bson.M{
		"email": req.Msg.GetEmail(),
	})
	if count > 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("email already exists"))
	}

	role := req.Msg.GetRole()
	if role == "" {
		role = "user"
	}

	data := &usermodel.User{
		Name:     req.Msg.GetName(),
		Email:    req.Msg.GetEmail(),
		Password: req.Msg.GetPassword(),
		Role:     strings.ToLower(role),
	}
	data.PreCreate()

	// hash password
	err := data.HashPassword()
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	result, err := repo.InsertOne(s.userCollection, data)
	if err != nil {
		log.Err(err).Msg("Error create user")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &userv1.CommonResponse{
		Data: "success",
	}

	switch v := result.InsertedID.(type) {
	case primitive.ObjectID:
		res.Data = v.Hex()
	case string:
		res.Data = v
	}

	return connect.NewResponse(res), nil
}

// UpdateUser is the user.v1.UserBiz.UpdateUser method.
func (s *Biz) UpdateUser(req *connect.Request[userv1.UpdateUserRequest]) (
	*connect.Response[userv1.CommonResponse], error,
) {
	id, err := primitive.ObjectIDFromHex(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	data, err := repo.FindOne[usermodel.User](s.userCollection, filter)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// count all users with filter
	count, _ := repo.CountDocuments(s.userCollection, bson.M{
		"_id":   bson.M{"$ne": id},
		"email": req.Msg.GetEmail(),
	})
	if count > 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("email already exists"))
	}

	data.Name = utils.StringCompareOrPassValue(data.Name, req.Msg.GetName())
	data.Email = utils.StringCompareOrPassValue(data.Email, req.Msg.GetEmail())
	data.Role = utils.StringCompareOrPassValue(data.Role, strings.ToLower(req.Msg.GetRole()))
	data.PreUpdate()

	obj := bson.M{"$set": data}
	if _, err = repo.UpdateOne(s.userCollection, filter, obj); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &userv1.CommonResponse{
		Data: req.Msg.GetId(),
	}
	return connect.NewResponse(res), nil
}

// DeleteUser is the user.v1.UserBiz.DeleteUser method.
func (s *Biz) DeleteUser(req *connect.Request[userv1.CommonUUIDRequest]) (
	*connect.Response[userv1.CommonResponse], error,
) {
	id, err := primitive.ObjectIDFromHex(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	// count all users with filter
	count, _ := repo.CountDocuments(s.userCollection, filter)
	if count <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("user does not exists"))
	}

	if _, err = repo.SoftDeleteOne(s.userCollection, filter); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &userv1.CommonResponse{
		Data: req.Msg.GetId(),
	}
	return connect.NewResponse(res), nil
}
