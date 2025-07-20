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

type DogHandler struct {
	dogService *services.DogService
	cfg        config.Config
}

func NewDogHandler(db *gorm.DB, redis *redis.Client, cfg config.Config) *DogHandler {
	return &DogHandler{
		dogService: services.NewDogService(db),
		cfg:        cfg,
	}
}

// CreateDog creates a new dog profile
func (h *DogHandler) CreateDog(c echo.Context) error {
	userID := middleware.GetUserID(c)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var req services.CreateDogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	dog, err := h.dogService.CreateDog(userUUID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusCreated, dog)
}

// GetUserDogs returns all dogs for the current user
func (h *DogHandler) GetUserDogs(c echo.Context) error {
	userID := middleware.GetUserID(c)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	dogs, err := h.dogService.GetUserDogs(userUUID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"dogs": dogs})
}

// GetDog returns a single dog by ID
func (h *DogHandler) GetDog(c echo.Context) error {
	userID := middleware.GetUserID(c)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	dogID := c.Param("dogId")
	dogUUID, err := uuid.Parse(dogID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid dog ID"})
	}

	dog, err := h.dogService.GetDog(dogUUID, userUUID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, dog)
}

// UpdateDog updates dog information
func (h *DogHandler) UpdateDog(c echo.Context) error {
	userID := middleware.GetUserID(c)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	dogID := c.Param("dogId")
	dogUUID, err := uuid.Parse(dogID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid dog ID"})
	}

	var req services.UpdateDogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	dog, err := h.dogService.UpdateDog(dogUUID, userUUID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, dog)
}

// DeleteDog deletes a dog profile
func (h *DogHandler) DeleteDog(c echo.Context) error {
	userID := middleware.GetUserID(c)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	dogID := c.Param("dogId")
	dogUUID, err := uuid.Parse(dogID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid dog ID"})
	}

	if err := h.dogService.DeleteDog(dogUUID, userUUID); err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Dog deleted successfully"})
}

// AddVaccinationRecord adds a vaccination record to a dog
// Disabled for simplified schema
func (h *DogHandler) AddVaccinationRecord(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Vaccination records not implemented in simplified schema"})
	
	/* Original implementation would be:
	userID := middleware.GetUserID(c)
	dogID := c.Param("dogId")

	var req services.VaccinationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	vaccination, err := h.dogService.AddVaccinationRecord(dogID, userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusCreated, vaccination)
	*/
}

// GetVaccinationRecords returns vaccination records for a dog
// Disabled for simplified schema
func (h *DogHandler) GetVaccinationRecords(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Vaccination records not implemented in simplified schema"})
}

// UpdateVaccinationRecord updates a vaccination record
// Disabled for simplified schema
func (h *DogHandler) UpdateVaccinationRecord(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Vaccination records not implemented in simplified schema"})
}

// DeleteVaccinationRecord deletes a vaccination record
// Disabled for simplified schema
func (h *DogHandler) DeleteVaccinationRecord(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Vaccination records not implemented in simplified schema"})
}

// SearchPublicDogs searches for public dogs
func (h *DogHandler) SearchPublicDogs(c echo.Context) error {
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

	dogs, total, err := h.dogService.SearchPublicDogs(query, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"dogs":   dogs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	return c.JSON(http.StatusOK, response)
}

// GetPersonalityTraits returns available personality traits
func (h *DogHandler) GetPersonalityTraits(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"traits": []string{
			"friendly", "playful", "energetic", "calm", "protective",
			"intelligent", "loyal", "independent", "social", "gentle",
			"curious", "brave", "shy", "affectionate", "alert",
		},
	})
}

// RegisterRoutes registers dog routes
func (h *DogHandler) RegisterRoutes(e *echo.Echo) {
	dogs := e.Group("/api/dogs", middleware.AuthMiddleware(h.cfg.JWT))

	// Dog management
	dogs.POST("", h.CreateDog)
	dogs.GET("", h.GetUserDogs)
	dogs.GET("/:dogId", h.GetDog)
	dogs.PUT("/:dogId", h.UpdateDog)
	dogs.DELETE("/:dogId", h.DeleteDog)

	// Vaccination management
	dogs.POST("/:dogId/vaccinations", h.AddVaccinationRecord)
	dogs.GET("/:dogId/vaccinations", h.GetVaccinationRecords)
	dogs.PUT("/:dogId/vaccinations/:vaccinationId", h.UpdateVaccinationRecord)
	dogs.DELETE("/:dogId/vaccinations/:vaccinationId", h.DeleteVaccinationRecord)

	// Public routes
	dogsPublic := e.Group("/api/dogs")
	dogsPublic.GET("/search", h.SearchPublicDogs, middleware.OptionalAuthMiddleware(h.cfg.JWT))
	dogsPublic.GET("/personality-traits", h.GetPersonalityTraits)
}