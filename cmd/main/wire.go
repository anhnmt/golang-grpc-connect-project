//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"

	"github.com/xdorro/golang-grpc-base-project/internal/interceptor"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/internal/usecase/auth"
	"github.com/xdorro/golang-grpc-base-project/internal/usecase/ping"
	"github.com/xdorro/golang-grpc-base-project/internal/usecase/user"
)

func initServer() server.IServer {
	wire.Build(
		ping.ProviderServiceSet,
		user.ProviderServiceSet,
		auth.ProviderServiceSet,
		interceptor.ProviderInterceptorSet,
		service.ProviderServiceSet,
		server.ProviderServerSet,
	)

	return &server.Server{}
}
