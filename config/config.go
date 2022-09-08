package config

import (
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// NewConfig initializes the config
func NewConfig() {
	viper.AutomaticEnv()
	// Replace env key
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	defaultConfig()

	// Set default logger
	defaultLogger()

	log.Info().
		Str("goarch", runtime.GOARCH).
		Str("goos", runtime.GOOS).
		Str("version", runtime.Version()).
		Msg("Runtime information")
}
