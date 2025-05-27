package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ParseValidationErrors(err error) string {
	var validationErrors []string
	customErrorMessage := map[string]string{
		"FirstName.required":       "First name is required",
		"LastName.required":        "Last name is required",
		"Email.required":           "Email is required",
		"Email.email":              "Email format is invalid",
		"Password.required":        "Password is required",
		"Password.strong_password": "Password must be at least 8 characters, include uppercase, lowercase, number, and special character",
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range errs {
			key := fmt.Sprintf("%s.%s", fe.Field(), fe.Tag())
			if msg, exists := customErrorMessage[key]; exists {
				validationErrors = append(validationErrors, msg)
			} else {
				validationErrors = append(validationErrors,
					fmt.Sprintf("Field %s failed validation on '%s'", fe.Field(), fe.Tag()))
			}
		}
	}

	return strings.Join(validationErrors, ", ")
}

func ValidatePassword(val validator.FieldLevel) bool {
	password := val.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString
		hasSpecial = regexp.MustCompile(`[\W_]`).MatchString
	)

	return hasUpper(password) && hasLower(password) && hasNumber(password) && hasSpecial(password)
}
