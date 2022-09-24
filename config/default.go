package config

import (
	"github.com/spf13/viper"
)

// defaultConfig is the default configuration for the application.
func defaultConfig() {
	// APP
	viper.SetDefault("APP_NAME", "golang-grpc-base-project")
	viper.SetDefault("APP_PORT", 8000)
	viper.SetDefault("PPROF_PORT", 6060)
	viper.SetDefault("APP_DEBUG", true)
	viper.SetDefault("APP_SECRET", "your-256-bit-secret")

	// LOG
	viper.SetDefault("LOG_PAYLOAD", true)
	viper.SetDefault("LOG_FILE_URL", "logs/data.log")

	// DATABASE
	viper.SetDefault("DB_URL", "mongodb://localhost:27017")
	viper.SetDefault("DB_NAME", "base")

	// REDIS
	viper.SetDefault("REDIS_URL", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
}
