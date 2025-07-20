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

type ModerationHandler struct {
	moderationService *services.ModerationService
	cfg               config.Config
}

func NewModerationHandler(db *gorm.DB, redis *redis.Client, cfg config.Config) *ModerationHandler {
	return &ModerationHandler{
		moderationService: services.NewModerationService(db, redis, cfg),
		cfg:               cfg,
	}
}

// Content reporting
func (h *ModerationHandler) CreateReport(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req services.CreateReportRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	report, err := h.moderationService.CreateReport(userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusCreated, report)
}

// User blocking
func (h *ModerationHandler) BlockUser(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req services.BlockUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.moderationService.BlockUser(userID, req); err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User blocked successfully"})
}

func (h *ModerationHandler) UnblockUser(c echo.Context) error {
	userID := middleware.GetUserID(c)
	blockedUserID := c.Param("userId")

	if err := h.moderationService.UnblockUser(userID, blockedUserID); err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User unblocked successfully"})
}

func (h *ModerationHandler) GetBlockedUsers(c echo.Context) error {
	userID := middleware.GetUserID(c)

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

	blockedUsers, total, err := h.moderationService.GetBlockedUsers(userID, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"blocked_users": blockedUsers,
		"total":         total,
		"limit":         limit,
		"offset":        offset,
	}

	return c.JSON(http.StatusOK, response)
}

// Safety settings
func (h *ModerationHandler) GetSafetySettings(c echo.Context) error {
	userID := middleware.GetUserID(c)

	settings, err := h.moderationService.GetSafetySettings(userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, settings)
}

func (h *ModerationHandler) UpdateSafetySettings(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req services.UpdateSafetySettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	settings, err := h.moderationService.UpdateSafetySettings(userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, settings)
}

// Admin functions
func (h *ModerationHandler) GetReports(c echo.Context) error {
	status := c.QueryParam("status")
	priority := c.QueryParam("priority")

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

	reports, total, err := h.moderationService.GetReports(status, priority, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"reports": reports,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *ModerationHandler) ReviewReport(c echo.Context) error {
	reviewerID := middleware.GetUserID(c)
	reportID := c.Param("reportId")

	var req services.ReviewReportRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	report, err := h.moderationService.ReviewReport(reviewerID, reportID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, report)
}

func (h *ModerationHandler) SuspendUser(c echo.Context) error {
	moderatorID := middleware.GetUserID(c)

	var req services.SuspendUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	suspension, err := h.moderationService.SuspendUser(moderatorID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusCreated, suspension)
}

// Utility endpoints
func (h *ModerationHandler) CheckUserBlocked(c echo.Context) error {
	userID := middleware.GetUserID(c)
	otherUserID := c.Param("userId")

	isBlocked, err := h.moderationService.IsUserBlocked(userID, otherUserID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"is_blocked": isBlocked,
	})
}

func (h *ModerationHandler) CheckUserSuspended(c echo.Context) error {
	userID := c.Param("userId")

	isSuspended, suspension, err := h.moderationService.IsUserSuspended(userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"is_suspended": isSuspended,
	}

	if isSuspended && suspension != nil {
		response["suspension"] = suspension
	}

	return c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers moderation routes
func (h *ModerationHandler) RegisterRoutes(e *echo.Echo) {
	moderation := e.Group("/api/moderation", middleware.AuthMiddleware(h.cfg.JWT))

	// Content reporting
	moderation.POST("/reports", h.CreateReport)

	// User blocking
	moderation.POST("/block", h.BlockUser)
	moderation.DELETE("/block/:userId", h.UnblockUser)
	moderation.GET("/blocked-users", h.GetBlockedUsers)

	// Safety settings
	moderation.GET("/safety-settings", h.GetSafetySettings)
	moderation.PUT("/safety-settings", h.UpdateSafetySettings)

	// Utility endpoints
	moderation.GET("/check-blocked/:userId", h.CheckUserBlocked)

	// Admin routes (would need admin middleware in real app)
	admin := e.Group("/api/admin/moderation", middleware.AuthMiddleware(h.cfg.JWT))
	admin.GET("/reports", h.GetReports)
	admin.PUT("/reports/:reportId/review", h.ReviewReport)
	admin.POST("/suspend", h.SuspendUser)
	admin.GET("/check-suspended/:userId", h.CheckUserSuspended)
}