package storage

import (
	"authenticationService/internal/models"
)

const (
	UserNotFound         = "user not found"
	RefreshTokenNotFound = "refresh token not found"
)

type TokenKeeper interface {
	CreateUser(user *models.User) (string, error)
	CreateRefreshToken(refreshToken *models.RefreshToken) (string, error)
	GetUserByID(id string) (*models.User, error)
	GetRefreshTokenByID(id string) (*models.RefreshToken, error)
	UpdateRefreshTokenStatus(id, newStatus string) error
}
