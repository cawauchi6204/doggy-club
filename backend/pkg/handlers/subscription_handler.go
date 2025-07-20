package handlers

import (
	"net/http"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/middleware"
	"github.com/doggyclub/backend/pkg/services"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	subscriptionService *services.SubscriptionService
	cfg                 config.Config
}

func NewSubscriptionHandler(db *gorm.DB, redis *redis.Client, cfg config.Config) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: services.NewSubscriptionService(db, cfg),
		cfg:                 cfg,
	}
}

// GetSubscriptionPlans returns all available subscription plans
func (h *SubscriptionHandler) GetSubscriptionPlans(c echo.Context) error {
	plans, err := h.subscriptionService.GetSubscriptionPlans()
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"plans": plans,
	})
}

// GetUserSubscription gets the user's current subscription
func (h *SubscriptionHandler) GetUserSubscription(c echo.Context) error {
	userID := middleware.GetUserID(c)

	subscription, err := h.subscriptionService.GetUserSubscription(userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	if subscription == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"subscription": nil,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"subscription": subscription,
	})
}

// CreateSubscription creates a new subscription
func (h *SubscriptionHandler) CreateSubscription(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req services.CreateSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	subscription, err := h.subscriptionService.CreateSubscription(userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusCreated, subscription)
}

// UpdateSubscription updates an existing subscription
func (h *SubscriptionHandler) UpdateSubscription(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req services.UpdateSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	subscription, err := h.subscriptionService.UpdateSubscription(userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, subscription)
}

// CancelSubscription cancels the user's subscription
func (h *SubscriptionHandler) CancelSubscription(c echo.Context) error {
	userID := middleware.GetUserID(c)

	// cancelAtPeriodEndStr := c.QueryParam("cancel_at_period_end")
	// cancelAtPeriodEnd := cancelAtPeriodEndStr == "true" // Not used in simplified schema

	err := h.subscriptionService.CancelSubscription(userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Subscription canceled successfully"})
}

// AddPaymentMethod adds a payment method
// Disabled for simplified schema
func (h *SubscriptionHandler) AddPaymentMethod(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Payment methods not implemented in simplified schema"})
}

// GetPaymentMethods returns user's payment methods
// Disabled for simplified schema
func (h *SubscriptionHandler) GetPaymentMethods(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Payment methods not implemented in simplified schema"})
	/*
	userID := middleware.GetUserID(c)
	paymentMethods, err := h.subscriptionService.GetPaymentMethods(userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"payment_methods": paymentMethods,
	})
	*/
}

// RemovePaymentMethod removes a payment method
// Disabled for simplified schema
func (h *SubscriptionHandler) RemovePaymentMethod(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Payment methods not implemented in simplified schema"})
}

// GetInvoices returns user's invoices
// Disabled for simplified schema
func (h *SubscriptionHandler) GetInvoices(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Invoices not implemented in simplified schema"})
}

// CheckFeatureAccess checks if user has access to a premium feature
// Disabled for simplified schema
func (h *SubscriptionHandler) CheckFeatureAccess(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"error": "Feature access checking not implemented in simplified schema"})
}

// RegisterRoutes registers subscription routes
func (h *SubscriptionHandler) RegisterRoutes(e *echo.Echo) {
	subscriptions := e.Group("/api/subscriptions", middleware.AuthMiddleware(h.cfg.JWT))

	// Subscription management
	subscriptions.GET("/plans", h.GetSubscriptionPlans)
	subscriptions.GET("/current", h.GetUserSubscription)
	subscriptions.POST("", h.CreateSubscription)
	subscriptions.PUT("", h.UpdateSubscription)
	subscriptions.DELETE("", h.CancelSubscription)

	// Payment methods
	subscriptions.GET("/payment-methods", h.GetPaymentMethods)
	subscriptions.POST("/payment-methods", h.AddPaymentMethod)
	subscriptions.DELETE("/payment-methods/:paymentMethodId", h.RemovePaymentMethod)

	// Invoices
	subscriptions.GET("/invoices", h.GetInvoices)

	// Feature access
	subscriptions.GET("/features/:featureCode/access", h.CheckFeatureAccess)

	// Public routes
	publicSubscriptions := e.Group("/api/subscriptions")
	publicSubscriptions.GET("/plans", h.GetSubscriptionPlans) // Make plans publicly accessible
}