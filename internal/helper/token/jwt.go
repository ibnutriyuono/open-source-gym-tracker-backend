package token

import (
	"caloria-backend/internal/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}
	jwtSecret :=  env.GetString("JWT_SECRET", "")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
