package app

import (
	"authenticationService/internal/config"
	"authenticationService/internal/storage"
	"log/slog"
)

type App struct {
	Config  *config.Config
	Logger  *slog.Logger
	Storage storage.TokenKeeper
}

func New(config *config.Config, storage storage.TokenKeeper, logger *slog.Logger) *App {
	return &App{
		Config:  config,
		Logger:  logger,
		Storage: storage,
	}
}
