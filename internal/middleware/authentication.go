package middleware

import (
	"caloria-backend/internal/env"
	"caloria-backend/internal/helper/response"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Authentication(DB *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")

			message := "Invalid token"
			if authorizationHeader == "" || !strings.HasPrefix(authorizationHeader, "Bearer") {
				response.SendJSON(w, http.StatusUnauthorized, nil, message)
				return
			}

			tokenString := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "Bearer"))
			secretKey := env.GetString("JWT_SECRET", "")

			token, err := jwt.Parse(string(tokenString), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(secretKey), nil
			})

			if errors.Is(err, jwt.ErrTokenExpired) {
				path := r.URL.Path
				switch path {
				case "/v1/auth/refresh":
					next.ServeHTTP(w, r)
				default:
					message := "Token has expired"
					response.SendJSON(w, http.StatusUnauthorized, nil, message)
					return
				}

			}

			if err != nil || !token.Valid {
				response.SendJSON(w, http.StatusUnauthorized, nil, message)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
