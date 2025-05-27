package user

import (
	"caloria-backend/internal/helper/hash"
	"caloria-backend/internal/helper/ip"
	"caloria-backend/internal/helper/response"
	"caloria-backend/internal/helper/token"
	"caloria-backend/internal/helper/validation"
	"caloria-backend/internal/model"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func (uc *UserController) FindAll(w http.ResponseWriter, r *http.Request) {
	users := []model.User{}
	result := uc.DB.Raw("SELECT id, first_name, last_name, email FROM users WHERE is_deleted = false").Scan(&users)
	message := "Successfully get all users"

	if result.Error != nil {
		message = "Failed to retrieve users"
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	if len(users) == 0 {
		message = "No users found"
		response.SendJSON(w, http.StatusOK, users, message)
		return
	}

	response.SendJSON(w, http.StatusOK, users, message)
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	message := "Success create user"

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		message = "Failed to decode request body"
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	validate := validator.New()
	validate.RegisterValidation("strong_password", validation.ValidatePassword)

	if err := validate.Struct(user); err != nil {
		message = validation.ParseValidationErrors(err)
		response.SendJSON(w, http.StatusBadRequest, nil, "Validation failed: "+message)
		return
	}

	uuidValueV7, err := uuid.NewV7()
	if err != nil {
		message = "Failed to create uuid"
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	user.ID = uuidValueV7.String()

	hashedPassword, err := hash.HashPassword(string(user.Password))

	if err != nil {
		message = "Failed to hash password"
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	query := "INSERT INTO users (id, first_name, last_name, password, email) VALUES (?, ?, ?, ?, ?) RETURNING id"
	result := uc.DB.Raw(query, user.ID, user.FirstName, user.LastName, hashedPassword, user.Email).Scan(&user.ID)
	if result.Error != nil {
		// check if email exists
		message = "Failed to create user"
		if strings.Contains(result.Error.Error(), "users_email_key") {
			message = "Email already exists"
			response.SendJSON(w, http.StatusInternalServerError, nil, message)
			return
		}
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	responseData := struct {
		ID string `json:"id"`
	}{
		ID: user.ID,
	}

	response.SendJSON(w, http.StatusOK, responseData, message)
}

func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// user := model.User{}
	existingUser := model.User{}
	message := "Success update user"
	query := "SELECT * FROM users WHERE id = ? AND is_deleted = false"
	result := uc.DB.Raw(query, id).Scan(&existingUser)
	if result.Error != nil {
		message = "Failed to retrieve user"
		if strings.Contains(result.Error.Error(), "invalid input syntax for type uuid") {
			message = "Invalid user ID"
			response.SendJSON(w, http.StatusBadRequest, nil, message)
			return
		}
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	if result.RowsAffected == 0 {
		message = "No users found"
		response.SendJSON(w, http.StatusOK, struct{}{}, message)
		return
	}

	if existingUser.IsDeleted {
		message = "User is already deleted"
		response.SendJSON(w, http.StatusBadRequest, struct{}{}, message)
		return
	}

	body := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&existingUser); err != nil {
		message = "Failed to decode request body"
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	if val, ok := body["first_name"].(string); ok && val != "" {
		existingUser.FirstName = val
	}
	if val, ok := body["last_name"].(string); ok && val != "" {
		existingUser.LastName = val
	}
	if val, ok := body["email"].(string); ok && val != "" {
		existingUser.Email = val
	}
	if val, ok := body["password"].(string); ok && val != "" {
		existingUser.Password = val
	}
	if val, ok := body["avatar"].(string); ok {
		existingUser.Avatar = val
	}
	if val, ok := body["phone_number"].(string); ok {
		existingUser.PhoneNumber = val
	}
	if val, ok := body["year_of_birth"].(string); ok {
		existingUser.YearOfBirth = val
	} else {
		existingUser.YearOfBirth = ""
	}
	if val, ok := body["month_of_birth"].(string); ok {
		existingUser.MonthOfBirth = val
	} else {
		existingUser.MonthOfBirth = ""
	}
	if val, ok := body["date_of_birth"].(string); ok {
		existingUser.DateOfBirth = val
	} else {
		existingUser.DateOfBirth = ""
	}
	if val, ok := body["gender"].(string); ok && val != "" {
		existingUser.Gender = &val
	} else {
		existingUser.Gender = nil
	}
	query = `
		UPDATE users SET 
			first_name = ?, 
			last_name = ?, 
			email = ?, 
			password = ?, 
			avatar = ?, 
			phone_number = ?, 
			year_of_birth = ?, 
			month_of_birth = ?, 
			date_of_birth = ?, 
			gender = ?, 
			updated_at = NOW() 
		WHERE id = ?
		`
	result = uc.DB.Exec(query,
		existingUser.FirstName,
		existingUser.LastName,
		existingUser.Email,
		existingUser.Password,
		existingUser.Avatar,
		existingUser.PhoneNumber,
		existingUser.YearOfBirth,
		existingUser.MonthOfBirth,
		existingUser.DateOfBirth,
		existingUser.Gender,
		existingUser.ID,
	)

	if result.Error != nil {
		message = "Failed to update user"
		if strings.Contains(result.Error.Error(), "invalid input syntax for type uuid") {
			message = "Invalid user ID"
			response.SendJSON(w, http.StatusBadRequest, nil, message)
			return
		}
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	validate := validator.New()
	validate.RegisterValidation("strong_password", validation.ValidatePassword)
	if err := validate.Struct(existingUser); err != nil {
		message = validation.ParseValidationErrors(err)
		response.SendJSON(w, http.StatusBadRequest, nil, "Validation failed: "+message)
		return
	}

	responseData := struct {
		ID        string `json:"id"`
		UpdatedAt string `json:"updated_at"`
	}{
		ID:        existingUser.ID,
		UpdatedAt: existingUser.UpdatedAt.String(),
	}
	response.SendJSON(w, http.StatusOK, responseData, message)
}

func (uc *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existingUser := model.User{}
	message := "Success delete user"
	query := "SELECT * FROM users WHERE id = ?"
	result := uc.DB.Raw(query, id).Scan(&existingUser)

	if result.Error != nil {
		message = "Failed to retrieve user"
		if strings.Contains(result.Error.Error(), "invalid input syntax for type uuid") {
			message = "Invalid user ID"
			response.SendJSON(w, http.StatusBadRequest, nil, message)
			return
		}
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	if result.RowsAffected == 0 {
		message = "No users found"
		response.SendJSON(w, http.StatusOK, struct{}{}, message)
		return
	}

	if existingUser.IsDeleted {
		message = "User is already deleted"
		response.SendJSON(w, http.StatusBadRequest, struct{}{}, message)
		return
	}

	query = "UPDATE users SET is_deleted = true, deleted_at = NOW() WHERE id = ?"
	result = uc.DB.Raw(query, id).Scan(&existingUser)

	if result.Error != nil {
		message = "Failed to update user"
		if strings.Contains(result.Error.Error(), "invalid input syntax for type uuid") {
			message = "Invalid user ID"
			response.SendJSON(w, http.StatusBadRequest, nil, message)
			return
		}
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	responseData := struct {
		ID string `json:"id"`
	}{
		ID: existingUser.ID,
	}
	response.SendJSON(w, http.StatusOK, responseData, message)

}

func (uc *UserController) FindById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := model.User{}
	result := uc.DB.Raw("SELECT id, first_name, last_name, email, password FROM users WHERE id = ? AND is_deleted = false", id).Scan(&user)
	message := "Successfully get user"

	if result.Error != nil {
		message = "Failed to retrieve user"
		if strings.Contains(result.Error.Error(), "invalid input syntax for type uuid") {
			message = "Invalid user ID"
			response.SendJSON(w, http.StatusBadRequest, nil, message)
			return
		}
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	if result.RowsAffected == 0 {
		message = "No users found"
		response.SendJSON(w, http.StatusOK, struct{}{}, message)
		return
	}

	if user.IsDeleted {
		message = "User is already deleted"
		response.SendJSON(w, http.StatusBadRequest, struct{}{}, message)
		return
	}

	responseData := struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}

	response.SendJSON(w, http.StatusOK, responseData, message)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	loginReq := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	message := "You are successfully logged in"
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		message = "Failed to decode request body"
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	if strings.TrimSpace(loginReq.Email) == "" || strings.TrimSpace(loginReq.Password) == "" {
		message = "Email and password are required"
		response.SendJSON(w, http.StatusBadRequest, nil, message)
		return
	}

	query := "SELECT id, email, password FROM users WHERE email = ? AND is_deleted = false"
	result := uc.DB.Raw(query, loginReq.Email).Scan((&user))

	if result.Error != nil || result.RowsAffected == 0 {
		message = "Invalid email or password"
	
		response.SendJSON(w, http.StatusUnauthorized, nil, message)
		return
	}

	isValidPassword := hash.CheckPasswordHash(loginReq.Password, user.Password)

	if isValidPassword != nil {
		message = "Invalid email or password"
		response.SendJSON(w, http.StatusUnauthorized, nil, message)
		return
	}
	accessTokenDuration := 1 * time.Minute
	accessToken, err := token.GenerateJWT(string(user.ID), accessTokenDuration)
	if err != nil {
		response.SendJSON(w, http.StatusInternalServerError, nil, "Failed to generate access token")
		return
	}

	refreshTokenDuration := 30 * 24 * time.Hour
	refreshToken, err := token.GenerateJWT(string(user.ID), refreshTokenDuration)
	if err != nil {
		response.SendJSON(w, http.StatusInternalServerError, nil, "Failed to generate access token")
		return
	}

	userAgent := r.UserAgent()
	clientIP := ip.GetClientIP(r)
	encrypt, err := token.EncryptWithPublicKey([]byte(clientIP), "public.pem")
	if err != nil {
		response.SendJSON(w, http.StatusInternalServerError, nil, "Failed to encrypt data")
		return
	}
	base64IP := base64.StdEncoding.EncodeToString(encrypt)

	query = "SELECT * FROM user_tokens WHERE user_id = ?"
	userToken := []model.UserToken{}
	result = uc.DB.Raw(query, user.ID).Scan(&userToken)

	if result.Error != nil {
		message = "Invalid email or password"
		response.SendJSON(w, http.StatusUnauthorized, nil, message)
		return
	}
	query = "INSERT INTO user_tokens (id, access_token, refresh_token, user_agent, ip_address, expires_at, is_revoked, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	uuidValueV7, err := uuid.NewV7()
	if err != nil {
		message = "Failed to create uuid"
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}

	result = uc.DB.Exec(query, uuidValueV7.String(), accessToken, refreshToken, userAgent, base64IP, nil, false, user.ID)
	if result.Error != nil {
		message = "Invalid user token"
		response.SendJSON(w, http.StatusUnauthorized, nil, message)
		return
	}

	responseData := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.SendJSON(w, http.StatusOK, responseData, message)
}

func (uc *UserController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	refreshTokenParam := params.Get("token")
	if refreshTokenParam == "" {
		message := "Token param required"
		response.SendJSON(w, http.StatusBadRequest, nil, message)
		return
	}

	userToken := model.UserToken{}
	query := "SELECT * FROM user_tokens WHERE refresh_token = ?"
	result := uc.DB.Raw(query, refreshTokenParam).Scan(&userToken)
	if result.Error != nil || result.RowsAffected == 0 {
		message := "Invalid refresh token"
	
		response.SendJSON(w, http.StatusUnauthorized, nil, message)
		return
	}

	userId := userToken.ID
	accessToken := userToken.AccessToken
	refreshToken := userToken.RefreshToken

	// generate new token
	accessTokenDuration := 1 * time.Minute
	newAccessToken, err := token.GenerateJWT(string(userId), accessTokenDuration)
	if err != nil {
		message := "Failed to generate new access token"
		response.SendJSON(w, http.StatusInternalServerError, nil, message)
		return
	}
	query = `
	UPDATE user_tokens SET 
		access_token = ?
	WHERE refresh_token = ?
`
	result = uc.DB.Exec(query, newAccessToken, refreshToken)
	if result.Error != nil || result.RowsAffected == 0 {
		message := "Failed to update access token"
		response.SendJSON(w, http.StatusUnauthorized, nil, message)
		return
	}

	responseData := struct {
		ID           string `json:"id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID:           userToken.UserID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.SendJSON(w, http.StatusOK, responseData, "message")
}
