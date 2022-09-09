package pingservice

import (
	"github.com/google/wire"
)

// ProviderServiceSet is Service providers.
var ProviderServiceSet = wire.NewSet(
	NewService,
	wire.Struct(new(Option), "*"),
)
