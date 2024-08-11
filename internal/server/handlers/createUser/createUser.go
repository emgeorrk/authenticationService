package createUser

import (
	"authenticationService/internal/app"
	"authenticationService/internal/logger"
	"authenticationService/internal/models"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type Request struct {
	Name                        string `json:"name" validate:"required"`
	Email                       string `json:"email" validate:"required,email"`
	MaxActiveTokenPairs         int    `json:"max_active_token_pairs"`
	AccessTokenLifetimeMinutes  int    `json:"access_token_lifetime_minutes"`
	RefreshTokenLifetimeMinutes int    `json:"refresh_token_lifetime_minutes"`
}

type Response struct {
	GUID  string `json:"GUID,omitempty"`
	Error string `json:"error,omitempty"`
}

func New(a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.handlers.auth.New"

		log := a.Logger.With(
			slog.String("handler", "auth"),
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request", logger.Err(err))

			w.WriteHeader(http.StatusBadRequest)

			render.JSON(w, r, Response{
				Error: err.Error(),
			})

			return
		}

		log.Info("request decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to validate request", logger.Err(err))

			w.WriteHeader(http.StatusBadRequest)

			render.JSON(w, r, Response{
				Error: err.Error(),
			})

			return
		}

		log.Info("request validated", slog.Any("request", req))

		newUser := &models.User{
			ID:                          uuid.New().String(),
			Name:                        req.Name,
			Email:                       req.Email,
			MaxActiveTokenPairs:         req.MaxActiveTokenPairs,
			AccessTokenLifetimeMinutes:  req.AccessTokenLifetimeMinutes,
			RefreshTokenLifetimeMinutes: req.RefreshTokenLifetimeMinutes,
		}

		// Создаем пользователя в базе данных
		id, err := a.Storage.CreateUser(newUser)
		if err != nil {
			log.Error("failed to create user", logger.Err(err))

			w.WriteHeader(http.StatusInternalServerError)

			render.JSON(w, r, Response{
				Error: "failed to create user",
			})

			return
		}

		// Возвращаем GUID созданного пользователя
		log.Info("user created", slog.String("GUID", id))

		w.WriteHeader(http.StatusCreated)

		render.JSON(w, r, Response{
			GUID: id,
		})

		return
	}
}
