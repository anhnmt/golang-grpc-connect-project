package interceptor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	permissionmodel "github.com/xdorro/golang-grpc-base-project/internal/module/permission/model"
	permissionrepo "github.com/xdorro/golang-grpc-base-project/internal/module/permission/repo"
	"github.com/xdorro/golang-grpc-base-project/pkg/casbin"
	"github.com/xdorro/golang-grpc-base-project/pkg/redis"
	"github.com/xdorro/golang-grpc-base-project/pkg/utils"
	"github.com/xdorro/golang-grpc-base-project/pkg/utils/constants"
)

var _ IInterceptor = (*Interceptor)(nil)

// IInterceptor is the interface that must be implemented by an interceptor.
type IInterceptor interface {
	UnaryInterceptor() connect.UnaryInterceptorFunc
}

// Option is an interceptor option struct.
type Option struct {
	Casbin         casbin.ICasbin
	Redis          redis.IRedis
	PermissionRepo permissionrepo.IRepo
}

// Interceptor is an interceptor struct.
type Interceptor struct {
	logPayload bool

	// options
	casbin         casbin.ICasbin
	redis          redis.IRedis
	permissionRepo permissionrepo.IRepo
}

// NewInterceptor returns a new interceptor.
func NewInterceptor(opt *Option) IInterceptor {
	i := &Interceptor{
		logPayload:     viper.GetBool("log.payload"),
		casbin:         opt.Casbin,
		redis:          opt.Redis,
		permissionRepo: opt.PermissionRepo,
	}

	return i
}

// UnaryInterceptor is a unary interceptor.
func (i *Interceptor) UnaryInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, request connect.AnyRequest) (
			connect.AnyResponse, error,
		) {
			response, err := next(ctx, request)

			// get full method
			permissions := i.getListPermissions()
			if len(permissions) == 0 {
				return i.logPayloadHandler(request, response, err)
			}

			// check permission
			procedure := request.Spec().Procedure
			per, ok := permissions[procedure]
			if !ok || per == nil {
				return i.logPayloadHandler(request, response, err)
			}

			// check require auth
			if per.RequireAuth {
				var token string
				token, err = utils.AuthFromHeader(request.Header(), utils.TokenType)
				if err != nil {
					log.Err(err).Msg("Error get token from header")
					return i.logPayloadHandler(request, nil, connect.NewError(connect.CodeUnauthenticated, err))
				}

				var claims *jwt.RegisteredClaims
				claims, err = utils.DecryptToken(token)
				if err != nil {
					return i.logPayloadHandler(request, nil, connect.NewError(connect.CodeInvalidArgument, err))
				}

				// check role
				var role string
				if len(claims.Audience) > 0 {
					role = claims.Audience[0]
				}

				allowed, _ := i.casbin.Enforcer().Enforce(role, procedure)
				if !allowed {
					err = fmt.Errorf("Permission denied")
					return i.logPayloadHandler(request, nil, connect.NewError(connect.CodePermissionDenied, err))
				}
			}

			// check require hash
			if !per.RequireHash {
				return i.logPayloadHandler(request, response, err)
			}

			// handler hash payload here
			// // test custom response
			//  response = connect.NewResponse(&permissionv1.CommonResponse{
			// 	Data: "ahihi",
			// })
			return i.logPayloadHandler(request, response, err)
		}
	}
}

// getAllPermissions returns all permissions.
func (i *Interceptor) getListPermissions() map[string]*permissionmodel.Permission {
	// get all permissions
	permissions := make(map[string]*permissionmodel.Permission)

	if val := redis.Get(i.redis, constants.ListAuthPermissionsKey); val != "" {
		_ = json.Unmarshal([]byte(val), &permissions)
		return permissions
	}

	// count all permissions with filter
	filter := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	// find all permissions with filter and option
	opt := options.
		Find().
		SetSort(bson.M{"created_at": -1})

	data, err := i.permissionRepo.Find(filter, opt)
	if err != nil {
		return permissions
	}

	for _, per := range data {
		permissions[per.Slug] = per
	}

	log.Info().
		Interface("permissions", permissions).
		Msg("Log get all permissions")

	go func() {
		_ = redis.SetObject(i.redis, constants.ListAuthPermissionsKey, permissions, 7*24*time.Hour)
	}()

	return permissions
}

// logPayloadHandler is a log payload handler.
func (i *Interceptor) logPayloadHandler(request connect.AnyRequest, response connect.AnyResponse, err error) (
	connect.AnyResponse, error,
) {
	// Log the payload
	if i.logPayload {
		go func(response connect.AnyResponse, err error) {
			logger := log.Info()
			if err != nil {
				logger = log.Error().Err(err)
			} else {
				logger.Interface("response", response.Any())
			}

			logger.
				Str("procedure", request.Spec().Procedure).
				Interface("request", request.Any()).
				Interface("header", request.Header()).
				Msg("Log payload interceptor")
		}(response, err)
	}

	return response, err
}
