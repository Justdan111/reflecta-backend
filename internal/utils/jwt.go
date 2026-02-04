package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"reflecta/internal/config"
)

func GenerateToken(userID string) (string, error) {
	secret := config.GetEnv("JWT_SECRET")

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(48 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
