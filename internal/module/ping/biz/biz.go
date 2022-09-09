package pingbiz

import (
	"github.com/bufbuild/connect-go"
	pingv1 "github.com/xdorro/proto-base-project/proto-gen-go/ping/v1"
)

var _ IPingBiz = &Biz{}

// IPingBiz ping service interface.
type IPingBiz interface {
	Ping(req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error)
}

// Biz struct.
type Biz struct {
	// option
}

// Option service option.
type Option struct {
}

// NewBiz new service.
func NewBiz() IPingBiz {
	return &Biz{}
}

// Ping is the ping.v1.PingBiz.Ping method.
func (s *Biz) Ping(req *connect.Request[pingv1.PingRequest]) (
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
