package config

import (
	"fmt"
	"os"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// NewConfig initializes the config
func NewConfig(env string) {
	viper.AutomaticEnv()

	// Replace env key
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	pwd, _ := os.Getwd()
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("%s/config", pwd))

	viper.SetConfigFile(fmt.Sprintf("%s/config/%s.toml", pwd, env))
	viper.SetConfigType("toml")
	viper.SetConfigName(env)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Error reading config file")
	}

	log.Info().
		Str("goarch", runtime.GOARCH).
		Str("goos", runtime.GOOS).
		Str("version", runtime.Version()).
		Msg("Runtime information")
}
