package interceptor

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/xdorro/golang-grpc-base-project/pkg/casbin"
	"github.com/xdorro/golang-grpc-base-project/pkg/redis"
)

var _ IInterceptor = (*Interceptor)(nil)

// IInterceptor is the interface that must be implemented by an interceptor.
type IInterceptor interface {
	UnaryInterceptor() connect.UnaryInterceptorFunc
}

// Option is an interceptor option struct.
type Option struct {
	Redis  redis.IRedis
	Casbin casbin.ICasbin
}

// Interceptor is an interceptor struct.
type Interceptor struct {
	logPayload bool

	// options
	redis  redis.IRedis
	casbin casbin.ICasbin
}

// NewInterceptor returns a new interceptor.
func NewInterceptor(opt *Option) IInterceptor {
	i := &Interceptor{
		logPayload: viper.GetBool("log.payload"),
		redis:      opt.Redis,
		casbin:     opt.Casbin,
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
			return i.logPayloadHandler(request, response, err)
		}
	}
}

// logPayloadHandler is a log payload handler.
func (i *Interceptor) logPayloadHandler(request connect.AnyRequest, response connect.AnyResponse, err error) (
	connect.AnyResponse, error,
) {
	// Log the payload
	if i.logPayload {
		defer func(response connect.AnyResponse, err error) {
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
