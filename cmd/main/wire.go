//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"

	"github.com/xdorro/golang-grpc-base-project/internal/interceptor"
	authmodule "github.com/xdorro/golang-grpc-base-project/internal/module/auth"
	pingmodule "github.com/xdorro/golang-grpc-base-project/internal/module/ping"
	usermodule "github.com/xdorro/golang-grpc-base-project/internal/module/user"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

func initServer() server.IServer {
	wire.Build(
		repo.ProviderRepoSet,
		pingmodule.ProviderModuleSet,
		usermodule.ProviderModuleSet,
		authmodule.ProviderModuleSet,
		interceptor.ProviderInterceptorSet,
		service.ProviderServiceSet,
		server.ProviderServerSet,
	)

	return &server.Server{}
}
