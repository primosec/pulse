package main

import (
	"os"
	"time"

	"github.com/primosec/pulse/internal/api"
	"github.com/primosec/pulse/internal/config"
	"github.com/primosec/pulse/internal/repository"
	"github.com/primosec/pulse/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Msg("Starting pulse...")

	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	jwtExpiry, err := time.ParseDuration(cfg.JWTExpiration)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse JWT expiration")
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, jwtExpiry)

	router := api.NewRouter(authService)

	log.Info().Str("port", cfg.Port).Msg("server listening")
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
