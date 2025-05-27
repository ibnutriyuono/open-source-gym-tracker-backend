package token

import (
	"caloria-backend/internal/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, duration time.Duration) (string, error) {
	now := time.Now()
	
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Unix() + int64(duration.Seconds()), 
	}
	
	jwtSecret := env.GetString("JWT_SECRET", "")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}