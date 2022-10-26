package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

var _ IServer = (*Server)(nil)

// IServer Server interface.
type IServer interface {
	Run() error
	Close() error
}

// Server struct.
type Server struct {
	mu sync.Mutex

	// config
	appName    string
	appPort    int
	pprofPort  int
	appDebug   bool
	logPayload bool

	// option
	grpcServer *grpc.Server
	httpServer *runtime.ServeMux
	service    service.IService
}

// Option server.
type Option struct {
	Service service.IService
}

// NewServer new server.
func NewServer(opt *Option) IServer {
	s := &Server{
		appName:    viper.GetString("app.name"),
		appDebug:   viper.GetBool("app.debug"),
		appPort:    viper.GetInt("app.port"),
		pprofPort:  viper.GetInt("pprof.port"),
		logPayload: viper.GetBool("log.payload"),
		service:    opt.Service,
	}

	s.NewGrpcServer()
	s.NewHttpServer()

	log.Info().
		Str("app-name", s.appName).
		Int("app-port", s.appPort).
		Msg("Server information loaded")

	return s
}

// Run runs the server.
func (s *Server) Run() error {
	// we're going to run the different protocol servers in parallel, so
	// make an errgroup
	group := new(errgroup.Group)

	// we need a webserver to get the pprof webserver
	if s.appDebug {
		group.Go(func() error {
			pprofPort := fmt.Sprintf(":%d", s.pprofPort)
			log.Info().Msgf("Starting pprof http://localhost%s", pprofPort)

			return http.ListenAndServe(pprofPort, nil)
		})
	}

	// Serve the http server on the http listener.
	group.Go(func() error {
		appPort := fmt.Sprintf(":%d", s.appPort)
		log.Info().Msgf("Starting application http://localhost%s", appPort)

		// create new http server
		srv := &http.Server{
			Addr: appPort,
			// Use h2c, so we can serve HTTP/2 without TLS.
			Handler:           s.grpcHandlerFunc(),
			ReadHeaderTimeout: time.Second,
			ReadTimeout:       1 * time.Minute,
			WriteTimeout:      1 * time.Minute,
			MaxHeaderBytes:    8 * 1024, // 8KiB
		}

		// run the server
		return srv.ListenAndServe()
	})

	return group.Wait()
}

// Close closes the server.
func (s *Server) Close() error {
	s.grpcServer.GracefulStop()

	return nil
}

func (s *Server) grpcHandlerFunc() http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.grpcServer.ServeHTTP(w, r)
			return
		}

		if !s.logPayload {
			s.httpServer.ServeHTTP(w, r)
			return
		}

		s.logPayloadHandler(w, r)
	}), &http2.Server{})
}

// logPayloadHandler is a log payload handler.
func (s *Server) logPayloadHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	// Work / inspect body. You may even modify it!

	// And now set a new body, which will simulate the same data we read:
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Create a response wrapper:
	mrw := &MyResponseWriter{
		ResponseWriter: w,
		buf:            &bytes.Buffer{},
	}

	logger := log.Info().
		Interface("header", r.Header.Clone())

	if len(body) > 0 {
		logger.RawJSON("body", body)
	}

	s.httpServer.ServeHTTP(mrw, r)

	logger.
		RawJSON("response", mrw.buf.Bytes())

	// Now inspect response, and finally send it out:
	// (You can also modify it before sending it out!)
	if _, err = io.Copy(w, mrw.buf); err != nil {
		log.Printf("Failed to send out response: %v", err)
	}

	logger.
		Msg("Log payload interceptor")
	return
}

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (mrw *MyResponseWriter) Write(p []byte) (int, error) {
	return mrw.buf.Write(p)
}
