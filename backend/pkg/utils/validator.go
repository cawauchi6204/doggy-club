package utils

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// Validator instance
var Validate *validator.Validate

func init() {
	Validate = validator.New()
	
	// Register custom validators
	Validate.RegisterValidation("email", validateEmail)
	Validate.RegisterValidation("password", validatePassword)
}

// validateEmail custom email validation
func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return emailRegex.MatchString(strings.ToLower(email))
}

// validatePassword custom password validation
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return len(password) >= 8
}

// ValidateStruct validates a struct
func ValidateStruct(s interface{}) error {
	return Validate.Struct(s)
}

// FormatValidationErrors formats validation errors
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError
	
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors = append(errors, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Message: formatFieldError(e),
			})
		}
	}
	
	return errors
}

// formatFieldError formats field error message
func formatFieldError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "password":
		return "Password must be at least 8 characters"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	default:
		return "Invalid value"
	}
}