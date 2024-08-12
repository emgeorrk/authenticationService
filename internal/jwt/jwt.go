package jwtlib

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const refreshTokenLength = 32

type JWTClaims struct {
	ClientIp string `json:"client_ip"`
	jwt.RegisteredClaims
}

func NewJWT(ip string, accessTokenDeadline time.Time) *jwt.Token {
	claims := JWTClaims{
		ClientIp: ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    "authenticationService",
			ExpiresAt: jwt.NewNumericDate(accessTokenDeadline),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token
}

func GenerateRefreshToken(accessToken string) (string, error) {
	const op = "jwtlib.GenerateRefreshToken"

	token := make([]byte, refreshTokenLength-7)
	if _, err := rand.Read(token); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Добавляем к refresh-токену последние 7 символов access-токена для связи
	return hex.EncodeToString(token) + accessToken[len(accessToken)-7:], nil
}

func ValidateToken(token string) (string, error) {
	return "user_id", nil
}
