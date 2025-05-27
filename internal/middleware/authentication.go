package middleware

import (
	"caloria-backend/internal/env"
	"caloria-backend/internal/helper/response"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Authentication(DB *gorm.DB) func(http.Handler) http.Handler {
	type contextKey string

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
				message := "Token has expired"
				response.SendJSON(w, http.StatusUnauthorized, nil, message)
				return

			}

			if err != nil || !token.Valid {
				fmt.Println(err)
				response.SendJSON(w, http.StatusUnauthorized, nil, message)
				return
			}

			accessToken := struct {
				AccessToken string `json:"access_token"`
				UserId      string `json:"id"`
			}{
				AccessToken: tokenString,
			}
			query := "SELECT user_id, access_token FROM user_tokens WHERE access_token = ?"
			result := DB.Raw(query, tokenString).Scan(&accessToken)

			if result.Error != nil {
				fmt.Println(result.Error.Error())
				message = "Database error"
				response.SendJSON(w, http.StatusInternalServerError, nil, message)
				return
			}

			if result.RowsAffected == 0 {
				message = "Invalid access token"
				response.SendJSON(w, http.StatusUnauthorized, nil, message)
				return
			}

			const contextKey contextKey = "userID"
			ctx := context.WithValue(r.Context(), contextKey, string(accessToken.UserId))
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
