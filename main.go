package main

import (
	"authenticationService/internal/config"
	"authenticationService/internal/logger"
	"authenticationService/storage/postgres"
	"log/slog"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoad()

	log := logger.NewLogger(cfg.Env)

	log.Info("Starting the application", slog.String("env", cfg.Env))

	storage, err := postgres.NewStorage(cfg.Storage)
	if err != nil {
		log.Error("failed to create storage", "error", err)
		return
	}

	log.Info("Storage connected successfully")

	_ = storage
}
