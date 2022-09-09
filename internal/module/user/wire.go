package usermodule

import (
	"github.com/google/wire"

	userbiz "github.com/xdorro/golang-grpc-base-project/internal/module/user/biz"
	userrepo "github.com/xdorro/golang-grpc-base-project/internal/module/user/repo"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	userrepo.ProviderRepoSet,
	userbiz.ProviderServiceSet,
)
