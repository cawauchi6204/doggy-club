package middleware

import (
	"context"
	"strings"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware creates JWT authentication middleware
func AuthMiddleware(cfg config.JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, map[string]string{"error": "Missing authorization header"})
			}

			// Extract token
			tokenString := utils.ExtractToken(authHeader)
			if tokenString == "" {
				return c.JSON(401, map[string]string{"error": "Invalid authorization header"})
			}

			// Validate token
			claims, err := utils.ValidateToken(tokenString, cfg)
			if err != nil {
				return c.JSON(401, map[string]string{"error": "Invalid or expired token"})
			}

			// Set user context
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

// OptionalAuthMiddleware for endpoints that work with or without auth
func OptionalAuthMiddleware(cfg config.JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				tokenString := utils.ExtractToken(authHeader)
				if tokenString != "" {
					claims, err := utils.ValidateToken(tokenString, cfg)
					if err == nil {
						c.Set("user_id", claims.UserID)
						c.Set("email", claims.Email)
					}
				}
			}
			return next(c)
		}
	}
}

// GetUserID gets user ID from context
func GetUserID(c echo.Context) string {
	if userID, ok := c.Get("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetEmail gets email from context
func GetEmail(c echo.Context) string {
	if email, ok := c.Get("email").(string); ok {
		return email
	}
	return ""
}

// RequireRole middleware for role-based access control
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// For now, we'll just check if user is authenticated
			// Role checking can be implemented later
			userID := GetUserID(c)
			if userID == "" {
				return c.JSON(403, map[string]string{"error": "Forbidden"})
			}
			return next(c)
		}
	}
}

// RateLimitMiddleware implements rate limiting
func RateLimitMiddleware(limit int, window int) echo.MiddlewareFunc {
	// For now, just pass through
	// Implement actual rate limiting with Redis later
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get client IP
			ip := c.RealIP()
			if ip == "" {
				ip = strings.Split(c.Request().RemoteAddr, ":")[0]
			}
			
			// For now, just pass through
			// Implement actual rate limiting with Redis later
			_ = ip // Use the variable to avoid unused error
			return next(c)
		}
	}
}

// CORSConfig returns CORS configuration
func CORSConfig() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if c.Request().Method == "OPTIONS" {
				return c.NoContent(204)
			}
			
			return next(c)
		}
	}
}

// RequestIDMiddleware adds request ID to context
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = utils.GenerateUUID()
			}
			
			c.Set("request_id", requestID)
			c.Response().Header().Set("X-Request-ID", requestID)
			
			return next(c)
		}
	}
}

// ContextMiddleware adds common context values
func ContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create context with values
			ctx := context.Background()
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}