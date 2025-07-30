package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/imdinnesh/openfinstack/services/auth/config"
)

func GenerateJWT(userID uint, ttl time.Duration, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(ttl).Unix(),
		"role":    role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Load().JWTSecret))
}	