package main

import (
	"github.com/rs/zerolog/log"

	"github.com/xdorro/golang-grpc-base-project/config"
)

func init() {
	config.NewConfig()
}

func main() {
	log.Info().Msg("Hello world")
}
