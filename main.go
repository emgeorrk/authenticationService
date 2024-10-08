package main

import (
	"authenticationService/internal/app"
	"authenticationService/internal/config"
	"authenticationService/internal/logger"
	"authenticationService/internal/server"
	smtplib "authenticationService/internal/smtp"
	"authenticationService/internal/storage/postgres"
	"log/slog"
	"net/http"
	smtp2 "net/smtp"

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

	var smtp smtp2.Auth
	if cfg.SMTP.IsEnabled {
		smtp = smtplib.New(cfg.SMTP)
		log.Info("Connected SMTP successfully")
	} else {
		log.Info("SMTP is disabled")
	}

	a := app.New(cfg, storage, log, smtp)

	router := server.New(*a)

	log.Info("Server started", slog.String("address", cfg.Address))
	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		log.Error("failed to start server", "err", err)
		return
	}

}
