package server

import (
	"authenticationService/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func NewServer(log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.URLFormat,
		logger.MiddlewareLogger(log),
	)

	return router
}
