package pingmodule

import (
	"github.com/google/wire"

	pingbiz "github.com/xdorro/golang-grpc-base-project/internal/module/ping/biz"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	pingbiz.ProviderServiceSet,
)
