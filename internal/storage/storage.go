package storage

import (
	"authenticationService/internal/models"
	"errors"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type TokenKeeper interface {
	CreateUser(user *models.User) (string, error)
	CreateRefreshToken(refreshToken *models.RefreshToken) (string, error)
	GetUserByID(id string) (*models.User, error)
	GetRefreshTokenByID(id string) (*models.RefreshToken, error)
	UpdateRefreshTokenStatus(id, newStatus string) error
}
