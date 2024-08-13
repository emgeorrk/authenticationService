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

func ValidateToken(secretKey, tokenString string) (*jwt.Token, error) {
	const op = "jwtlib.ValidateToken"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: unexpected signing method: %v", op, token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, fmt.Errorf("%s: token is invalid", op)
	}

	return token, nil
}
