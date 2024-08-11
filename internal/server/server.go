package server

import (
	"authenticationService/internal/app"
	"authenticationService/internal/logger"
	"authenticationService/internal/server/handlers/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(a app.App) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.URLFormat,
		logger.MiddlewareLogger(a.Logger),
	)

	router.Get("/auth", auth.New(a))

	return router
}
