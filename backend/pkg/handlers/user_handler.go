package handlers

import (
	"net/http"
	"strconv"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/middleware"
	"github.com/doggyclub/backend/pkg/services"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService *services.UserService
	cfg         config.Config
}

func NewUserHandler(db *gorm.DB, redis *redis.Client, cfg config.Config) *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(db, redis, cfg),
		cfg:         cfg,
	}
}

// GetProfile returns current user's profile
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, profile)
}

// UpdateProfile updates user profile
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req services.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	profile, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, profile)
}

// UpdatePrivacySettings updates privacy settings
func (h *UserHandler) UpdatePrivacySettings(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req services.UpdatePrivacySettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.userService.UpdatePrivacySettings(userID, req); err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Privacy settings updated successfully"})
}

// UpdateNotificationPreferences updates notification preferences
// Disabled for simplified schema
func (h *UserHandler) UpdateNotificationPreferences(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Notification preferences not implemented in simplified schema"})
	
	/*
	// Original implementation would be:
	userID := middleware.GetUserID(c)
	var req services.UpdateNotificationPreferencesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := h.userService.UpdateNotificationPreferences(userID, req); err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Notification preferences updated successfully"})
	*/
}

// GetUserCurrency returns user's currency balance
// Disabled for simplified schema
func (h *UserHandler) GetUserCurrency(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Currency system not implemented in simplified schema"})
	
	/*
	// Original implementation would be:
	userID := middleware.GetUserID(c)
	currency, err := h.userService.GetUserCurrency(userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}
	return c.JSON(http.StatusOK, currency)
	*/
}

// DeleteAccount deletes user account
func (h *UserHandler) DeleteAccount(c echo.Context) error {
	userID := middleware.GetUserID(c)

	if err := h.userService.DeleteAccount(userID); err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Account deleted successfully"})
}

// SearchUsers searches for users (admin function)
func (h *UserHandler) SearchUsers(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Query parameter 'q' is required"})
	}

	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offsetStr := c.QueryParam("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, total, err := h.userService.SearchUsers(query, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	return c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers user routes
func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	users := e.Group("/api/users", middleware.AuthMiddleware(h.cfg.JWT))

	// Profile management
	users.GET("/profile", h.GetProfile)
	users.PUT("/profile", h.UpdateProfile)
	users.DELETE("/profile", h.DeleteAccount)

	// Settings
	users.PUT("/privacy", h.UpdatePrivacySettings)
	users.PUT("/notifications", h.UpdateNotificationPreferences)

	// Currency
	users.GET("/currency", h.GetUserCurrency)

	// Search (for admin/public use)
	users.GET("/search", h.SearchUsers)
}