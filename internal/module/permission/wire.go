package permissionmodule

import (
	"github.com/google/wire"

	permissionbiz "github.com/xdorro/golang-grpc-base-project/internal/module/permission/biz"
	permissionrepo "github.com/xdorro/golang-grpc-base-project/internal/module/permission/repo"
	permissionservice "github.com/xdorro/golang-grpc-base-project/internal/module/permission/service"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	permissionrepo.ProviderRepoSet,
	permissionbiz.ProviderBizSet,
	permissionservice.ProviderServiceSet,
)
