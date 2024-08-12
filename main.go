package main

import (
	"authenticationService/internal/app"
	"authenticationService/internal/config"
	"authenticationService/internal/logger"
	"authenticationService/internal/server"
	"authenticationService/internal/storage/postgres"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"
)

// @title authenticationService API
// @version 1.0
// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.MustLoad()

	log := logger.NewLogger(cfg.Env)

	log.Info("Starting the application", slog.String("env", cfg.Env))

	storage, err := postgres.NewStorage(cfg.Storage)
	if err != nil {
		log.Error("failed to create storage", "error", err)
		return
	}

	log.Info("Connected PostgreSQL successfully",
		slog.String("host", cfg.Storage.Host),
		slog.Int("port", cfg.Storage.Port),
		slog.String("database", cfg.Storage.Database),
	)

	a := app.New(cfg, storage, log)

	router := server.New(*a)

	log.Info("Server started", slog.String("address", cfg.Address), slog.Int("port", cfg.Port))
	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		log.Error("failed to start server", "err", err)
		return
	}

}
