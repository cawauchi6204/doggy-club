package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceType represents the type of device
type DeviceType string

const (
	DeviceTypeIOS     DeviceType = "ios"
	DeviceTypeAndroid DeviceType = "android"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEncounter NotificationType = "encounter"
	NotificationTypeGift      NotificationType = "gift"
	NotificationTypeLike      NotificationType = "like"
)

// DeviceToken represents a device token for push notifications
type DeviceToken struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Token      string     `gorm:"type:varchar(255);not null" json:"token"`
	DeviceType DeviceType `gorm:"type:varchar(20);not null" json:"device_type"`
	LastActive time.Time  `gorm:"index" json:"last_active"`

	// Relationship
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// BeforeCreate sets the ID before creating the device token
func (dt *DeviceToken) BeforeCreate(tx *gorm.DB) error {
	if dt.ID == uuid.Nil {
		dt.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the DeviceToken model
func (DeviceToken) TableName() string {
	return "device_tokens"
}

// Notification represents a notification log
type Notification struct {
	ID      uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID  uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	Type    NotificationType `gorm:"type:varchar(20);not null" json:"type"`
	Message string           `gorm:"type:text;not null" json:"message"`
	SentAt  time.Time        `gorm:"default:CURRENT_TIMESTAMP;index" json:"sent_at"`

	// Relationship
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// BeforeCreate sets the ID before creating the notification
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the Notification model
func (Notification) TableName() string {
	return "notifications"
}