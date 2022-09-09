package interceptor

import (
	"github.com/google/wire"
)

// ProviderInterceptorSet is Interceptor providers.
var ProviderInterceptorSet = wire.NewSet(
	NewInterceptor,
	wire.Struct(new(Option), "*"),
)
