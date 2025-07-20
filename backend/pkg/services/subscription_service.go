package services

import (
	"errors"
	"time"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/models"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionService struct {
	db  *gorm.DB
	cfg config.Config
}

func NewSubscriptionService(db *gorm.DB, cfg config.Config) *SubscriptionService {
	return &SubscriptionService{
		db:  db,
		cfg: cfg,
	}
}

// CreateSubscriptionRequest represents subscription creation request
type CreateSubscriptionRequest struct {
	PlanID string `json:"plan_id" validate:"required"`
}

// UpdateSubscriptionRequest represents subscription update request
type UpdateSubscriptionRequest struct {
	PlanID *string `json:"plan_id,omitempty"`
}

// CreateSubscription creates a new subscription (simplified)
func (s *SubscriptionService) CreateSubscription(userID string, req CreateSubscriptionRequest) (*models.UserSubscription, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(utils.FormatValidationErrors(err))
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	planUUID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return nil, errors.New("invalid plan ID format")
	}

	// Get subscription plan
	var plan models.SubscriptionPlan
	if err := s.db.Where("id = ?", planUUID).First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.WrapError(err, "failed to find subscription plan")
	}

	// Check if user already has an active subscription
	var existingSubscription models.UserSubscription
	err = s.db.Where("user_id = ? AND status = ?", userUUID, models.SubscriptionStatusActive).First(&existingSubscription).Error
	if err == nil {
		return nil, errors.New("user already has an active subscription")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.WrapError(err, "failed to check existing subscription")
	}

	// Create local subscription record (simplified)
	subscription := models.UserSubscription{
		UserID:    userUUID,
		PlanID:    planUUID,
		Status:    models.SubscriptionStatusActive,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, plan.DurationMonths, 0),
	}

	if err := s.db.Create(&subscription).Error; err != nil {
		return nil, utils.WrapError(err, "failed to create subscription")
	}

	// Load the plan details
	if err := s.db.Preload("Plan").Where("id = ?", subscription.ID).First(&subscription).Error; err != nil {
		return nil, utils.WrapError(err, "failed to load subscription")
	}

	return &subscription, nil
}

// GetUserSubscription returns user's current subscription
func (s *SubscriptionService) GetUserSubscription(userID string) (*models.UserSubscription, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var subscription models.UserSubscription
	if err := s.db.Preload("Plan").Where("user_id = ? AND status = ?", userUUID, models.SubscriptionStatusActive).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.WrapError(err, "failed to get subscription")
	}

	return &subscription, nil
}

// GetSubscriptionPlans returns all available subscription plans
func (s *SubscriptionService) GetSubscriptionPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	if err := s.db.Find(&plans).Error; err != nil {
		return nil, utils.WrapError(err, "failed to get subscription plans")
	}

	return plans, nil
}

// UpdateSubscription updates user's subscription plan
func (s *SubscriptionService) UpdateSubscription(userID string, req UpdateSubscriptionRequest) (*models.UserSubscription, error) {
	if req.PlanID == nil {
		return nil, errors.New("plan_id is required")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	planUUID, err := uuid.Parse(*req.PlanID)
	if err != nil {
		return nil, errors.New("invalid plan ID format")
	}

	// Get current subscription
	var subscription models.UserSubscription
	if err := s.db.Where("user_id = ? AND status = ?", userUUID, models.SubscriptionStatusActive).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.WrapError(err, "failed to find subscription")
	}

	// Get new plan
	var newPlan models.SubscriptionPlan
	if err := s.db.Where("id = ?", planUUID).First(&newPlan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.WrapError(err, "failed to find plan")
	}

	// Update subscription
	subscription.PlanID = planUUID
	subscription.EndDate = time.Now().AddDate(0, newPlan.DurationMonths, 0)

	if err := s.db.Save(&subscription).Error; err != nil {
		return nil, utils.WrapError(err, "failed to update subscription")
	}

	// Reload with plan details
	if err := s.db.Preload("Plan").Where("id = ?", subscription.ID).First(&subscription).Error; err != nil {
		return nil, utils.WrapError(err, "failed to reload subscription")
	}

	return &subscription, nil
}

// CancelSubscription cancels user's subscription
func (s *SubscriptionService) CancelSubscription(userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	var subscription models.UserSubscription
	if err := s.db.Where("user_id = ? AND status = ?", userUUID, models.SubscriptionStatusActive).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound
		}
		return utils.WrapError(err, "failed to find subscription")
	}

	// Update status to canceled
	subscription.Status = models.SubscriptionStatusCanceled
	if err := s.db.Save(&subscription).Error; err != nil {
		return utils.WrapError(err, "failed to cancel subscription")
	}

	return nil
}

// CheckSubscriptionStatus checks and updates expired subscriptions
func (s *SubscriptionService) CheckSubscriptionStatus() error {
	// Update expired subscriptions
	if err := s.db.Model(&models.UserSubscription{}).
		Where("status = ? AND end_date < ?", models.SubscriptionStatusActive, time.Now()).
		Update("status", models.SubscriptionStatusCanceled).Error; err != nil {
		return utils.WrapError(err, "failed to update expired subscriptions")
	}

	return nil
}