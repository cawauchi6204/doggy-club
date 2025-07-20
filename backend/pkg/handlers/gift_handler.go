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

type GiftHandler struct {
	giftService *services.GiftService
	cfg         config.Config
}

func NewGiftHandler(db *gorm.DB, redis *redis.Client, cfg config.Config) *GiftHandler {
	return &GiftHandler{
		giftService: services.NewGiftService(db),
		cfg:         cfg,
	}
}

// GetGiftCatalog returns available gifts
func (h *GiftHandler) GetGiftCatalog(c echo.Context) error {
	gifts := h.giftService.GetAvailableGiftTypes()
	giftCatalog := make([]map[string]interface{}, len(gifts))
	for i, giftType := range gifts {
		giftCatalog[i] = map[string]interface{}{
			"type": giftType,
			"name": giftType, // Simplified
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"gifts": giftCatalog,
	})
}

// SendGift sends a gift to another dog
func (h *GiftHandler) SendGift(c echo.Context) error {
	// userID := middleware.GetUserID(c) // Not used in simplified schema

	var req services.SendGiftRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	gift, err := h.giftService.SendGift(req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusCreated, gift)
}

// GetSentGifts returns gifts sent by user
func (h *GiftHandler) GetSentGifts(c echo.Context) error {
	userID := middleware.GetUserID(c)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
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

	gifts, total, err := h.giftService.GetSentGifts(userUUID, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"gifts":  gifts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	return c.JSON(http.StatusOK, response)
}

// GetReceivedGifts returns gifts received by user's dogs
func (h *GiftHandler) GetReceivedGifts(c echo.Context) error {
	userID := middleware.GetUserID(c)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
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

	gifts, total, err := h.giftService.GetReceivedGifts(userUUID, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"gifts":  gifts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	return c.JSON(http.StatusOK, response)
}

// ExchangeGift exchanges a virtual gift for real reward
// Disabled for simplified schema
func (h *GiftHandler) ExchangeGift(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "ExchangeGift not implemented in simplified schema"})
}

// GetGiftRankings returns gift rankings
func (h *GiftHandler) GetGiftRankings(c echo.Context) error {
	period := c.QueryParam("period")
	if period == "" {
		period = "all_time"
	}

	limitStr := c.QueryParam("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	rankings, err := h.giftService.GetPopularGiftTypes(limit)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"rankings": rankings,
		"period":   period,
		"limit":    limit,
	})
}

// PurchaseCurrency allows users to purchase in-app currency
// Disabled for simplified schema
func (h *GiftHandler) PurchaseCurrency(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Currency system not implemented in simplified schema"})
}

// GetTransactionHistory returns user's transaction history
// Disabled for simplified schema
func (h *GiftHandler) GetTransactionHistory(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Transaction history not implemented in simplified schema"})
	
	/*
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

	// Original implementation would be:
	// transactions, total, err := h.giftService.GetTransactionHistory(userID, limit, offset)
	// if err != nil {
	//	status, apiErr := utils.HTTPError(err)
	//	return c.JSON(status, map[string]interface{}{"error": apiErr})
	// }
	// response := map[string]interface{}{
	//	"transactions": transactions,
	//	"total":        total,
	//	"limit":        limit,
	//	"offset":       offset,
	// }
	// return c.JSON(http.StatusOK, response)
	*/
}

// RegisterRoutes registers gift routes
func (h *GiftHandler) RegisterRoutes(e *echo.Echo) {
	gifts := e.Group("/api/gifts", middleware.AuthMiddleware(h.cfg.JWT))

	// Gift catalog
	gifts.GET("/catalog", h.GetGiftCatalog)

	// Gift purchasing and management
	gifts.POST("/send", h.SendGift)
	gifts.GET("/sent", h.GetSentGifts)
	gifts.GET("/received", h.GetReceivedGifts)
	gifts.POST("/exchange", h.ExchangeGift)

	// Currency management
	currency := e.Group("/api/currency", middleware.AuthMiddleware(h.cfg.JWT))
	currency.POST("/purchase", h.PurchaseCurrency)
	currency.GET("/transactions", h.GetTransactionHistory)

	// Public routes
	giftsPublic := e.Group("/api/gifts")
	giftsPublic.GET("/rankings", h.GetGiftRankings)
}