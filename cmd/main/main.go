package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/xdorro/golang-grpc-base-project/config"
)

func init() {
	config.NewConfig()
}

func main() {
	exit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	// New server
	srv := initServer()

	// Run server
	srv.Run()

	<-exit
	if err := srv.Close(); err != nil {
		log.Err(err).Msg("Failed to close server")
		return
	}

	log.Info().Msg("Graceful shutdown complete")
}
