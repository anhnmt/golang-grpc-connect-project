package authbiz

import (
	"context"
	"fmt"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	authv1 "github.com/xdorro/proto-base-project/proto-gen-go/auth/v1"
	"github.com/xdorro/proto-base-project/proto-gen-go/auth/v1/authv1connect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"

	usermodel "github.com/xdorro/golang-grpc-base-project/internal/module/user/model"
	userrepo "github.com/xdorro/golang-grpc-base-project/internal/module/user/repo"
	"github.com/xdorro/golang-grpc-base-project/utils"
)

var _ IAuthService = &Service{}

// IAuthService auth service interface.
type IAuthService interface {
	authv1connect.AuthServiceHandler
}

// Service struct.
type Service struct {
	// option
	userRepo userrepo.IRepo

	authv1connect.UnimplementedAuthServiceHandler
}

// Option service option.
type Option struct {
	UserRepo userrepo.IRepo
}

// NewService new service.
func NewService(opt *Option) IAuthService {
	s := &Service{
		userRepo: opt.UserRepo,
	}

	return s
}

// Login is the auth.v1.AuthService.Login method.
func (s *Service) Login(_ context.Context, req *connect.Request[authv1.LoginRequest]) (
	*connect.Response[authv1.TokenResponse], error,
) {
	filter := bson.M{
		"email": req.Msg.GetEmail(),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	data, err := s.userRepo.FindOne(filter)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// verify password
	if !data.ComparePassword(req.Msg.GetPassword()) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("password is incorrect"))
	}

	// generate a new auth token
	res, err := s.generateAuthToken(data)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}

// RevokeToken is the auth.v1.AuthService.RevokeToken method.
func (s *Service) RevokeToken(_ context.Context, req *connect.Request[authv1.TokenRequest]) (
	*connect.Response[authv1.CommonResponse], error,
) {
	token := req.Msg.GetToken()

	// verify & remove old token
	_, err := s.removeAuthToken(token)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := &authv1.CommonResponse{
		Token: token,
	}

	return connect.NewResponse(res), nil
}

// RefreshToken is the auth.v1.AuthService.RefreshToken method.
func (s *Service) RefreshToken(_ context.Context, req *connect.Request[authv1.TokenRequest]) (
	*connect.Response[authv1.TokenResponse], error,
) {
	// verify & remove old token
	claims, err := s.removeAuthToken(req.Msg.GetToken())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	id, err := primitive.ObjectIDFromHex(claims.Subject)
	if err != nil {
		log.Err(err).Msg("Failed find user by id")
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

	// generate a new auth token
	res, err := s.generateAuthToken(data)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}

// generateAuthToken generates a new auth token for the user.
func (s *Service) generateAuthToken(data *usermodel.User) (
	*authv1.TokenResponse, error,
) {
	uid := data.ID.Hex()
	sessionID := uuid.NewString()
	now := time.Now()
	refreshExpire := now.Add(utils.RefreshExpire)
	accessExpire := now.Add(utils.AccessExpire)

	result := &authv1.TokenResponse{
		TokenType:     utils.TokenType,
		RefreshExpire: refreshExpire.Unix(),
		AccessExpire:  accessExpire.Unix(),
	}

	var eg errgroup.Group

	// Create a new refreshToken
	eg.Go(func() error {
		var err error
		result.RefreshToken, err = utils.EncryptToken(&jwt.RegisteredClaims{
			Subject:   uid,
			ExpiresAt: jwt.NewNumericDate(refreshExpire),
			ID:        sessionID,
		})
		if err != nil {
			return err
		}

		// key := fmt.Sprintf(utils.AuthSessionKey, uid, sessionID)
		// err = redis.Set(s.redis, key, result.RefreshToken, utils.RefreshExpire)
		// if err != nil {
		// 	log.Err(err).Msg("Failed to set auth session")
		// 	return err
		// }

		return nil
	})

	// Create a new accessToken
	eg.Go(func() error {
		var err error
		result.AccessToken, err = utils.EncryptToken(&jwt.RegisteredClaims{
			Subject:   uid,
			ExpiresAt: jwt.NewNumericDate(accessExpire),
			ID:        sessionID,
			Audience:  []string{data.Role},
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err := eg.Wait(); err != nil {

		return nil, err
	}

	return result, nil
}

// removeAuthToken removes the auth token from the redis.
func (s *Service) removeAuthToken(token string) (*jwt.RegisteredClaims, error) {
	// verify refresh token
	claims, err := utils.DecryptToken(token)
	if err != nil {
		return nil, err
	}

	log.Info().
		Interface("claims", claims).
		Msg("Token decrypted")

	// // check if the refresh token is existed
	// key := fmt.Sprintf(utils.AuthSessionKey, claims.Subject, claims.ID)
	// if check := redis.Exists(s.redis, key); !check {
	// 	return nil, fmt.Errorf("token is not found")
	// }
	//
	// if err = redis.Del(s.redis, key); err != nil {
	// 	return nil, err
	// }

	return claims, nil
}
