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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token
}

func GenerateRefreshToken() (string, error) {
	const op = "jwtlib.GenerateRefreshToken"

	token := make([]byte, refreshTokenLength)
	if _, err := rand.Read(token); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return hex.EncodeToString(token), nil
}

func ValidateToken(token string) (string, error) {
	return "user_id", nil
}
