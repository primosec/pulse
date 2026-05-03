package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Port              string
	Env               string
	DatabaseURL       string
	RedisURL          string
	JWTSecret         string
	JWTExpiration     string
	WorkerConcurrency int
}

func Load() *Config {
	godotenv.Load()
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("JWT_EXPIRATION", "24h")
	viper.SetDefault("WORKER_CONCURRENCY", 10)

	cfg := &Config{
		Port:              viper.GetString("PORT"),
		Env:               viper.GetString("ENV"),
		DatabaseURL:       viper.GetString("DATABASE_URL"),
		RedisURL:          viper.GetString("REDIS_URL"),
		JWTSecret:         viper.GetString("JWT_SECRET"),
		JWTExpiration:     viper.GetString("JWT_EXPIRATION"),
		WorkerConcurrency: viper.GetInt("WORKER_CONCURRENCY"),
	}

	log.Info().Str("env", cfg.Env).Str("port", cfg.Port).Msg("Configuration loaded")
	return cfg
}
