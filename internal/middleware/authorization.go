package middleware

import (
	"caloria-backend/internal/helper/response"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func CheckPermission(DB *gorm.DB, permissionName string,) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			contextKey := "userID"
			userID, ok := r.Context().Value(contextKey).(string)
            if !ok || userID == "" {
				message := "Unauthorized"
				response.SendJSON(w, http.StatusUnauthorized, nil, message)
				return
            }

			fmt.Println(userID, "blues in closet")
			var count int64
			query := `
				SELECT COUNT(*) 
				FROM users u
				JOIN user_roles ur ON ur.user_id = u.id
				JOIN roles r ON r.id = ur.role_id
				JOIN role_permissions rp ON rp.role_id = r.id
				JOIN permissions p ON p.id = rp.permission_id
				WHERE u.id = ? AND p.name = ?
			`

			fmt.Println(query)

			if err := DB.Raw(query, userID, permissionName).Count(&count).Error; err != nil || count == 0 {
				message := "Forbidden"
				response.SendJSON(w, http.StatusForbidden, nil, message)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
