package userrepo

import (
	"github.com/google/wire"
)

// ProviderRepoSet is Repo providers.
var ProviderRepoSet = wire.NewSet(
	NewRepo,
	wire.Struct(new(Option), "*"),
)
