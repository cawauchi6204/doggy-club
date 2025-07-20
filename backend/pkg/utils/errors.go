package utils

import (
	"errors"
	"fmt"
	"net/http"
)

// Common errors
var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrInternalServer      = errors.New("internal server error")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrTokenExpired        = errors.New("token expired")
	ErrInvalidToken        = errors.New("invalid token")
)

// APIError represents an API error response
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements error interface
func (e APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(code string, message string, details interface{}) APIError {
	return APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// HTTPError maps error to HTTP status code
func HTTPError(err error) (int, APIError) {
	switch {
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest, NewAPIError("INVALID_INPUT", err.Error(), nil)
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, NewAPIError("UNAUTHORIZED", err.Error(), nil)
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, NewAPIError("FORBIDDEN", err.Error(), nil)
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, NewAPIError("NOT_FOUND", err.Error(), nil)
	case errors.Is(err, ErrConflict):
		return http.StatusConflict, NewAPIError("CONFLICT", err.Error(), nil)
	case errors.Is(err, ErrInvalidCredentials):
		return http.StatusUnauthorized, NewAPIError("INVALID_CREDENTIALS", "Invalid email or password", nil)
	case errors.Is(err, ErrTokenExpired):
		return http.StatusUnauthorized, NewAPIError("TOKEN_EXPIRED", "Token has expired", nil)
	case errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized, NewAPIError("INVALID_TOKEN", "Invalid token", nil)
	default:
		return http.StatusInternalServerError, NewAPIError("INTERNAL_ERROR", "An internal error occurred", nil)
	}
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationError creates validation error response
func NewValidationError(errors []ValidationError) APIError {
	return NewAPIError("VALIDATION_ERROR", "Validation failed", errors)
}

// WrapError wraps an error with context
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}