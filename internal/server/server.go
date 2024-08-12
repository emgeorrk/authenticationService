package server

import (
	_ "authenticationService/docs"
	"authenticationService/internal/app"
	"authenticationService/internal/logger"
	"authenticationService/internal/server/handlers/auth"
	"authenticationService/internal/server/handlers/createUser"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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

	router.Post("/users", createUser.New(a))
	router.Post("/auth", auth.New(a))

	if a.Config.Env == "local" {
		router.Get("/swagger/*", httpSwagger.WrapHandler)
	}

	return router
}
