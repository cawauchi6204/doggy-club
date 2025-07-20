package handlers

import (
	"net/http"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/middleware"
	"github.com/doggyclub/backend/pkg/services"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService *services.AuthService
	cfg         config.Config
}

func NewAuthHandler(db *gorm.DB, redis *redis.Client, cfg config.Config) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(db, cfg.JWT.Secret),
		cfg:         cfg,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	var req services.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	resp, err := h.authService.Register(req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusCreated, resp)
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	var req services.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, resp)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// RefreshToken method not implemented in simplified service
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "RefreshToken not implemented"})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c echo.Context) error {
	// userID := middleware.GetUserID(c) // Not used in simplified implementation
	
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	c.Bind(&req)

	// Logout method not implemented in simplified service
	// In a real implementation, would invalidate refresh tokens

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	userID := middleware.GetUserID(c)
	
	var req struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": utils.NewValidationError(utils.FormatValidationErrors(err)),
		})
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := h.authService.ChangePassword(userUUID, req.OldPassword, req.NewPassword); err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

// ForgotPassword handles password reset request
func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": utils.NewValidationError(utils.FormatValidationErrors(err)),
		})
	}

	// ForgotPassword method not implemented in simplified service
	// In a real implementation, would send password reset email

	return c.JSON(http.StatusOK, map[string]string{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": utils.NewValidationError(utils.FormatValidationErrors(err)),
		})
	}

	// ResetPassword method not implemented in simplified service
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "ResetPassword not implemented"})
	
	// In a real implementation, would handle password reset here:
	// 
	// if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
	//	status, apiErr := utils.HTTPError(err)
	//	return c.JSON(status, map[string]interface{}{"error": apiErr})
	// }
	// return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
}

// GetMe returns current user info
func (h *AuthHandler) GetMe(c echo.Context) error {
	userID := middleware.GetUserID(c)
	
	// This would normally call a user service
	// For now, return user ID
	return c.JSON(http.StatusOK, map[string]string{
		"user_id": userID,
		"email":   middleware.GetEmail(c),
	})
}

// RegisterRoutes registers auth routes
func (h *AuthHandler) RegisterRoutes(e *echo.Echo) {
	auth := e.Group("/api/auth")
	
	// Public routes
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/forgot-password", h.ForgotPassword)
	auth.POST("/reset-password", h.ResetPassword)
	
	// Protected routes
	auth.POST("/logout", h.Logout, middleware.AuthMiddleware(h.cfg.JWT))
	auth.POST("/change-password", h.ChangePassword, middleware.AuthMiddleware(h.cfg.JWT))
	auth.GET("/me", h.GetMe, middleware.AuthMiddleware(h.cfg.JWT))
}