package handlers

import (
	"net/http"
	"strconv"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/middleware"
	"github.com/doggyclub/backend/pkg/services"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type EncounterHandler struct {
	encounterService *services.EncounterService
	cfg              config.Config
}

func NewEncounterHandler(db *gorm.DB, redis *redis.Client, cfg config.Config) *EncounterHandler {
	return &EncounterHandler{
		encounterService: services.NewEncounterService(db),
		cfg:              cfg,
	}
}

// DetectEncounters handles encounter detection
func (h *EncounterHandler) DetectEncounters(c echo.Context) error {
	// userID := middleware.GetUserID(c) // Not used in this method

	var req struct {
		DogID         string  `json:"dog_id" validate:"required"`
		RadiusMeters  float64 `json:"radius_meters" validate:"required,min=1,max=10000"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	dogUUID, err := uuid.Parse(req.DogID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid dog ID format"})
	}

	response, err := h.encounterService.DetectEncounters(dogUUID, req.RadiusMeters)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, response)
}

// GetEncounterHistory returns dog's encounter history
func (h *EncounterHandler) GetEncounterHistory(c echo.Context) error {
	dogID := c.Param("dogId")
	dogUUID, err := uuid.Parse(dogID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid dog ID format"})
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

	encounters, total, err := h.encounterService.GetDogEncounters(dogUUID, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"encounters": encounters,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	}

	return c.JSON(http.StatusOK, response)
}

// GetEncounterDetails returns detailed encounter information
// Disabled for simplified schema
func (h *EncounterHandler) GetEncounterDetails(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "GetEncounterDetails not implemented in simplified schema"})
}

// UpdateEncounterPreferences updates user's encounter preferences
// Disabled for simplified schema
func (h *EncounterHandler) UpdateEncounterPreferences(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "UpdateEncounterPreferences not implemented in simplified schema"})
	
	/*
	// Original implementation would be:
	userID := middleware.GetUserID(c)
	var req services.UpdateEncounterSettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := h.encounterService.UpdateEncounterPreferences(userID, req); err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Encounter preferences updated successfully"})
	*/
}

// RegisterRoutes registers encounter routes
func (h *EncounterHandler) RegisterRoutes(e *echo.Echo) {
	encounters := e.Group("/api/encounters", middleware.AuthMiddleware(h.cfg.JWT))

	// Encounter detection
	encounters.POST("/detect", h.DetectEncounters)

	// Encounter history
	encounters.GET("/history", h.GetEncounterHistory)
	encounters.GET("/:encounterId/details", h.GetEncounterDetails)

	// Settings
	encounters.PUT("/preferences", h.UpdateEncounterPreferences)
}