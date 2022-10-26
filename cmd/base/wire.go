//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"

	usermodule "github.com/xdorro/golang-grpc-base-project/internal/module/user"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

func initServer() server.IServer {
	wire.Build(
		repo.ProviderRepoSet,
		// redis.ProviderRedisSet,
		usermodule.ProviderModuleSet,
		// casbin.ProviderCasbinSet,
		service.ProviderServiceSet,
		server.ProviderServerSet,
	)

	return &server.Server{}
}
