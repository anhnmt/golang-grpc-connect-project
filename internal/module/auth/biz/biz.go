package authbiz

import (
	"fmt"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	authv1 "github.com/xdorro/proto-base-project/proto-gen-go/auth/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"

	usermodel "github.com/xdorro/golang-grpc-base-project/internal/module/user/model"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
	"github.com/xdorro/golang-grpc-base-project/pkg/utils"
)

var _ IAuthBiz = &Biz{}

// IAuthBiz auth service interface.
type IAuthBiz interface {
	Login(req *connect.Request[authv1.LoginRequest]) (*connect.Response[authv1.TokenResponse], error)
	RevokeToken(req *connect.Request[authv1.TokenRequest]) (*connect.Response[authv1.CommonResponse], error)
	RefreshToken(req *connect.Request[authv1.TokenRequest]) (*connect.Response[authv1.TokenResponse], error)
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
func NewBiz(opt *Option) IAuthBiz {
	s := &Biz{
		userCollection: opt.Repo.CollectionModel(&usermodel.User{}),
	}

	return s
}

func (s *Biz) Login(req *connect.Request[authv1.LoginRequest]) (
	*connect.Response[authv1.TokenResponse], error,
) {
	filter := bson.M{
		"email": req.Msg.GetEmail(),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	data, err := repo.FindOne[usermodel.User](s.userCollection, filter)
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

// RevokeToken is the auth.v1.AuthBiz.RevokeToken method.
func (s *Biz) RevokeToken(req *connect.Request[authv1.TokenRequest]) (
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

// RefreshToken is the auth.v1.AuthBiz.RefreshToken method.
func (s *Biz) RefreshToken(req *connect.Request[authv1.TokenRequest]) (
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
	data, err := repo.FindOne[usermodel.User](s.userCollection, filter)
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
func (s *Biz) generateAuthToken(data *usermodel.User) (
	*authv1.TokenResponse, error,
) {
	uid := data.Id
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
func (s *Biz) removeAuthToken(token string) (*jwt.RegisteredClaims, error) {
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
