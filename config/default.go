package config

import (
	"github.com/spf13/viper"
)

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
