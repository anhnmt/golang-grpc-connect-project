package pingservice

import (
	"context"

	"github.com/bufbuild/connect-go"
	pingv1 "github.com/xdorro/proto-base-project/proto-gen-go/ping/v1"
	"github.com/xdorro/proto-base-project/proto-gen-go/ping/v1/pingv1connect"

	pingbiz "github.com/xdorro/golang-grpc-base-project/internal/module/ping/biz"
)

var _ IPingService = &Service{}

// IPingService ping service interface.
type IPingService interface {
	pingv1connect.PingServiceHandler
}

// Service struct.
type Service struct {
	// option
	pingBiz pingbiz.IPingBiz

	pingv1connect.UnimplementedPingServiceHandler
}

// Option service option.
type Option struct {
	PingBiz pingbiz.IPingBiz
}

// NewService new service.
func NewService(opt *Option) IPingService {
	return &Service{
		pingBiz: opt.PingBiz,
	}
}

// Ping is the ping.v1.PingService.Ping method.
func (s *Service) Ping(_ context.Context, req *connect.Request[pingv1.PingRequest]) (
	*connect.Response[pingv1.PingResponse], error,
) {
	return s.pingBiz.Ping(req)
}
