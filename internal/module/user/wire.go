package usermodule

import (
	"github.com/google/wire"

	userbiz "github.com/xdorro/golang-grpc-base-project/internal/module/user/biz"
	userrepo "github.com/xdorro/golang-grpc-base-project/internal/module/user/repo"
	userservice "github.com/xdorro/golang-grpc-base-project/internal/module/user/service"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	userrepo.ProviderRepoSet,
	userbiz.ProviderBizSet,
	userservice.ProviderServiceSet,
)
