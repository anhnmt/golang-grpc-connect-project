package authmodule

import (
	"github.com/google/wire"

	authbiz "github.com/xdorro/golang-grpc-base-project/internal/module/auth/biz"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	authbiz.ProviderServiceSet,
)
