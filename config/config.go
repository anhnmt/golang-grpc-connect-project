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

	// set default config
	defaultConfig()

	log.Info().
		Str("goarch", runtime.GOARCH).
		Str("goos", runtime.GOOS).
		Str("version", runtime.Version()).
		Msg("Runtime information")
}

// defaultConfig is the default configuration for the application.
func defaultConfig() {
	// APP
	viper.SetDefault("app.name", "golang-grpc-base-project")
	viper.SetDefault("app.port", 8088)
	viper.SetDefault("app.debug", true)
	viper.SetDefault("APP_SECRET", "your-256-bit-secret")

	// LOG
	viper.SetDefault("LOG_PAYLOAD", true)
	viper.SetDefault("LOG_FILE_URL", "logs/data.log")

	// DATABASE
	viper.SetDefault("database.url", "mongodb://localhost:27017")
	viper.SetDefault("database.name", "base")

	// REDIS
	viper.SetDefault("redis.url", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
}
