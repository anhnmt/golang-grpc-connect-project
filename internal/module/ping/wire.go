package pingmodule

import (
	"github.com/google/wire"

	pingbiz "github.com/xdorro/golang-grpc-base-project/internal/module/ping/biz"
	pingservice "github.com/xdorro/golang-grpc-base-project/internal/module/ping/service"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	pingbiz.ProviderBizSet,
	pingservice.ProviderServiceSet,
)
