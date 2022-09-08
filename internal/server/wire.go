package server

import (
	"net/http"

	"github.com/google/wire"
)

// ProviderServerSet is Server providers.
var ProviderServerSet = wire.NewSet(
	http.NewServeMux,
	NewServer,
	wire.Struct(new(Option), "*"),
)
