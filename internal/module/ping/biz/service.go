package pingbiz

import (
	"context"

	"github.com/bufbuild/connect-go"
	pingv1 "github.com/xdorro/proto-base-project/proto-gen-go/ping/v1"
	"github.com/xdorro/proto-base-project/proto-gen-go/ping/v1/pingv1connect"
)

var _ IPingService = &Service{}

// IPingService ping service interface.
type IPingService interface {
	pingv1connect.PingServiceHandler
}

// Service struct.
type Service struct {
	// option

	pingv1connect.UnimplementedPingServiceHandler
}

// Option service option.
type Option struct {
}

// NewService new service.
func NewService() IPingService {
	return &Service{}
}

// Ping is the ping.v1.PingService.Ping method.
func (s *Service) Ping(_ context.Context, req *connect.Request[pingv1.PingRequest]) (
	*connect.Response[pingv1.PingResponse], error,
) {
	text := req.Msg.GetText()
	if text == "" {
		text = "pong"
	}

	res := &pingv1.PingResponse{
		Text: text,
	}

	return connect.NewResponse(res), nil
}
