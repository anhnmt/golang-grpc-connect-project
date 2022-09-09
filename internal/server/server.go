package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/utils"
)

var _ IServer = (*Server)(nil)

// IServer Server interface.
type IServer interface {
	Run()
	Close() error
}

// Server struct.
type Server struct {
	// config
	appName   string
	appPort   int
	pprofPort int
	appDebug  bool

	// option
	mux     *http.ServeMux
	service service.IService

	mu   sync.Mutex
	http *http.Server
}

// Option server.
type Option struct {
	Mux     *http.ServeMux
	Service service.IService
}

// NewServer new server.
func NewServer(opt *Option) IServer {
	s := &Server{
		appName:   viper.GetString("APP_NAME"),
		appPort:   viper.GetInt("APP_PORT"),
		pprofPort: viper.GetInt("PPROF_PORT"),
		appDebug:  viper.GetBool("APP_DEBUG"),
		mux:       opt.Mux,
		service:   opt.Service,
	}

	log.Info().
		Str("app-name", s.appName).
		Int("app-port", s.appPort).
		Msg("Server information loaded")

	return s
}

// Run runs the server.
func (s *Server) Run() {
	appPort := fmt.Sprintf(":%d", s.appPort)
	log.Info().Msgf("Starting application http://localhost%s", appPort)

	// create new http server
	s.setServer(&http.Server{
		Addr: appPort,
		// Use h2c, so we can serve HTTP/2 without TLS.
		Handler: h2c.NewHandler(
			s.customHandler(),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      1 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	})

	// we need a webserver to get the pprof webserver
	if s.appDebug {
		go func(s *Server) {
			pprofPort := fmt.Sprintf(":%d", s.pprofPort)
			log.Info().Msgf("Starting pprof http://localhost:%s", pprofPort)

			err := http.ListenAndServe(pprofPort, nil)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal().Err(err).Msg("Failed to run pprof server")
			}
		}(s)
	}

	// Serve the http server on the http listener.
	go func(s *Server) {
		err := s.http.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to run http server")
		}
	}(s)
}

// Close closes the server.
func (s *Server) Close() error {
	g, gCtx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		ctx, cancel := context.WithTimeout(gCtx, 10*time.Second)
		defer cancel()

		if err := s.http.Shutdown(ctx); err != nil {
			log.Err(err).Msg("Failed to shutdown http server")
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

// Server adds a new server.
func (s *Server) setServer(http *http.Server) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.http = http
}

// customHandler adds custom handlers to the server.
func (s *Server) customHandler() http.Handler {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseWithJson(w, http.StatusOK, "Hello, World!")
	})

	return newCORS().Handler(s.mux)
}
