package userbiz

import (
	"context"
	"fmt"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog/log"
	userv1 "github.com/xdorro/proto-base-project/proto-gen-go/user/v1"
	"github.com/xdorro/proto-base-project/proto-gen-go/user/v1/userv1connect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	usermodel "github.com/xdorro/golang-grpc-base-project/internal/module/user/model"
	userrepo "github.com/xdorro/golang-grpc-base-project/internal/module/user/repo"
	"github.com/xdorro/golang-grpc-base-project/utils"
)

var _ IUserService = &Service{}

// IUserService user service interface.
type IUserService interface {
	userv1connect.UserServiceHandler
}

// Service struct.
type Service struct {
	// option
	userRepo userrepo.IRepo

	userv1connect.UnimplementedUserServiceHandler
}

// Option service option.
type Option struct {
	UserRepo userrepo.IRepo
}

// NewService new service.
func NewService(opt *Option) IUserService {
	s := &Service{
		userRepo: opt.UserRepo,
	}

	return s
}

// FindAllUsers is the user.v1.UserService.FindAllUsers method.
func (s *Service) FindAllUsers(_ context.Context, req *connect.Request[userv1.FindAllUsersRequest]) (
	*connect.Response[userv1.FindAllUsersResponse], error,
) {

	// count all users with filter
	filter := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	count, _ := s.userRepo.CountDocuments(filter)
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
	data, err := s.userRepo.Find(filter, opt)
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

// FindUserByID is the user.v1.UserService.FindUserByID method.
func (s *Service) FindUserByID(_ context.Context, req *connect.Request[userv1.CommonUUIDRequest]) (
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

	data, err := s.userRepo.FindOne(filter, opt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := usermodel.UserToProto(data)
	return connect.NewResponse(res), nil
}

// CreateUser is the user.v1.UserService.CreateUser method.
func (s *Service) CreateUser(_ context.Context, req *connect.Request[userv1.CreateUserRequest]) (
	*connect.Response[userv1.CommonResponse], error,
) {
	// count all users with filter
	count, _ := s.userRepo.CountDocuments(bson.M{
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

	oid, err := s.userRepo.InsertOne(data)
	if err != nil {
		log.Err(err).Msg("Error create user")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	resID := oid.InsertedID.(primitive.ObjectID)
	res := &userv1.CommonResponse{
		Data: resID.Hex(),
	}
	return connect.NewResponse(res), nil
}

// UpdateUser is the user.v1.UserService.UpdateUser method.
func (s *Service) UpdateUser(_ context.Context, req *connect.Request[userv1.UpdateUserRequest]) (
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
	data, err := s.userRepo.FindOne(filter)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// count all users with filter
	count, _ := s.userRepo.CountDocuments(bson.M{
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
	if _, err = s.userRepo.UpdateOne(filter, obj); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &userv1.CommonResponse{
		Data: req.Msg.GetId(),
	}
	return connect.NewResponse(res), nil
}

// DeleteUser is the user.v1.UserService.DeleteUser method.
func (s *Service) DeleteUser(_ context.Context, req *connect.Request[userv1.CommonUUIDRequest]) (
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
	count, _ := s.userRepo.CountDocuments(filter)
	if count <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("user does not exists"))
	}

	if _, err = s.userRepo.SoftDeleteOne(filter); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &userv1.CommonResponse{
		Data: req.Msg.GetId(),
	}
	return connect.NewResponse(res), nil
}
