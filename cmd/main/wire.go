//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"

	"github.com/xdorro/golang-grpc-base-project/internal/interceptor"
	authbiz "github.com/xdorro/golang-grpc-base-project/internal/module/auth/biz"
	pingbiz "github.com/xdorro/golang-grpc-base-project/internal/module/ping/biz"
	userbiz "github.com/xdorro/golang-grpc-base-project/internal/module/user/biz"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

func initServer() server.IServer {
	wire.Build(
		pingbiz.ProviderServiceSet,
		userbiz.ProviderServiceSet,
		authbiz.ProviderServiceSet,
		interceptor.ProviderInterceptorSet,
		service.ProviderServiceSet,
		server.ProviderServerSet,
	)

	return &server.Server{}
}
