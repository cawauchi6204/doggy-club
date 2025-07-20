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

type NotificationHandler struct {
	notificationService *services.NotificationService
	cfg                 config.Config
}

func NewNotificationHandler(db *gorm.DB, redis *redis.Client, cfg config.Config) *NotificationHandler {
	return &NotificationHandler{
		notificationService: services.NewNotificationService(db, cfg),
		cfg:                 cfg,
	}
}

// RegisterDevice registers a device for push notifications
func (h *NotificationHandler) RegisterDevice(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req services.RegisterDeviceTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	device, err := h.notificationService.RegisterDeviceToken(userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusCreated, device)
}

// UnregisterDevice removes a device
// Disabled for simplified schema
func (h *NotificationHandler) UnregisterDevice(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "UnregisterDevice not implemented in simplified schema"})
}

// GetNotificationPreferences gets user's notification preferences
// Disabled for simplified schema
func (h *NotificationHandler) GetNotificationPreferences(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Notification preferences not implemented in simplified schema"})
}

// UpdateNotificationPreferences updates user's notification preferences
// Disabled for simplified schema
func (h *NotificationHandler) UpdateNotificationPreferences(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Notification preferences not implemented in simplified schema"})
	/*
	// Original implementation would be:
	userID := middleware.GetUserID(c)
	var req services.UpdateNotificationPreferencesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	prefs, err := h.notificationService.UpdateNotificationPreferences(userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}
	return c.JSON(http.StatusOK, prefs)
	*/
}

// GetNotifications gets user's notifications
func (h *NotificationHandler) GetNotifications(c echo.Context) error {
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

	notifications, total, err := h.notificationService.GetUserNotifications(userID, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"notifications": notifications,
		"total":         total,
		"limit":         limit,
		"offset":        offset,
	}

	return c.JSON(http.StatusOK, response)
}

// MarkNotificationAsRead marks a notification as read
func (h *NotificationHandler) MarkNotificationAsRead(c echo.Context) error {
	// userID := middleware.GetUserID(c)
	// notificationID := c.Param("notificationId")

	// MarkNotificationAsRead not implemented in simplified schema
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "MarkNotificationAsRead not implemented"})
	
	/*
	// Original implementation would be:
	// userID := middleware.GetUserID(c)
	// notificationID := c.Param("notificationId")
	// if err := h.notificationService.MarkNotificationAsRead(userID, notificationID); err != nil {
	//	status, apiErr := utils.HTTPError(err)
	//	return c.JSON(status, map[string]interface{}{"error": apiErr})
	// }
	// return c.JSON(http.StatusOK, map[string]string{"message": "Notification marked as read"})
	*/
}

// MarkAllNotificationsAsRead marks all notifications as read
// Disabled for simplified schema
func (h *NotificationHandler) MarkAllNotificationsAsRead(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "MarkAllNotificationsAsRead not implemented in simplified schema"})
}

// GetUnreadCount gets the count of unread notifications
// Disabled for simplified schema
func (h *NotificationHandler) GetUnreadCount(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "GetUnreadCount not implemented in simplified schema"})
	/*
	userID := middleware.GetUserID(c)
	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"unread_count": count,
	})
	*/
}

// SendNotification sends a notification (admin only)
// Disabled for simplified schema
func (h *NotificationHandler) SendNotification(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "SendNotification not implemented in simplified schema"})
}

// RegisterRoutes registers notification routes
func (h *NotificationHandler) RegisterRoutes(e *echo.Echo) {
	notifications := e.Group("/api/notifications", middleware.AuthMiddleware(h.cfg.JWT))

	// Device management
	notifications.POST("/devices", h.RegisterDevice)
	notifications.DELETE("/devices/:deviceToken", h.UnregisterDevice)

	// Notification preferences
	notifications.GET("/preferences", h.GetNotificationPreferences)
	notifications.PUT("/preferences", h.UpdateNotificationPreferences)

	// Notifications
	notifications.GET("", h.GetNotifications)
	notifications.PUT("/:notificationId/read", h.MarkNotificationAsRead)
	notifications.PUT("/read-all", h.MarkAllNotificationsAsRead)
	notifications.GET("/unread-count", h.GetUnreadCount)

	// Admin routes (would need admin middleware in real app)
	admin := e.Group("/api/admin/notifications", middleware.AuthMiddleware(h.cfg.JWT))
	admin.POST("/send", h.SendNotification)
}