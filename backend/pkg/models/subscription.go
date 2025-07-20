package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusCanceled  SubscriptionStatus = "canceled"
)

// SubscriptionPlan represents a subscription plan
type SubscriptionPlan struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name           string    `gorm:"type:varchar(50);not null" json:"name"`
	Price          float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	DurationMonths int       `gorm:"type:integer;not null" json:"duration_months"`
	Features       string    `gorm:"type:jsonb" json:"features"` // JSON array of features

	// Relationships
	UserSubscriptions []UserSubscription `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"user_subscriptions,omitempty"`
}

// BeforeCreate sets the ID before creating the subscription plan
func (sp *SubscriptionPlan) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the SubscriptionPlan model
func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

// UserSubscription represents a user's subscription to a plan
type UserSubscription struct {
	ID        uuid.UUID          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	PlanID    uuid.UUID          `gorm:"type:uuid;not null;index" json:"plan_id"`
	StartDate time.Time          `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time          `gorm:"type:date;not null" json:"end_date"`
	Status    SubscriptionStatus `gorm:"type:varchar(20);not null" json:"status"`

	// Relationships
	User User             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Plan SubscriptionPlan `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"plan,omitempty"`
}

// BeforeCreate sets the ID before creating the user subscription
func (us *UserSubscription) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the UserSubscription model
func (UserSubscription) TableName() string {
	return "user_subscriptions"
}

// IsActive checks if the subscription is currently active
func (us *UserSubscription) IsActive() bool {
	return us.Status == SubscriptionStatusActive && time.Now().Before(us.EndDate)
}

// IsExpired checks if the subscription has expired
func (us *UserSubscription) IsExpired() bool {
	return time.Now().After(us.EndDate)
}