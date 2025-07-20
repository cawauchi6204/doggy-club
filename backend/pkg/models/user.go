package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Visibility represents user privacy settings
type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username     string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	Email        string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Visibility   Visibility `gorm:"type:varchar(20);default:'public'" json:"visibility"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Dogs              []Dog              `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"dogs,omitempty"`
	UserSubscriptions []UserSubscription `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"subscriptions,omitempty"`
	DeviceTokens      []DeviceToken      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"device_tokens,omitempty"`
	Notifications     []Notification     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"notifications,omitempty"`
}

// BeforeCreate sets the ID before creating the user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// IsPublic returns true if the user's profile is public
func (u *User) IsPublic() bool {
	return u.Visibility == VisibilityPublic
}

// IsPrivate returns true if the user's profile is private
func (u *User) IsPrivate() bool {
	return u.Visibility == VisibilityPrivate
}

// ToPublicUser returns a user object without sensitive information
func (u *User) ToPublicUser() *User {
	return &User{
		ID:         u.ID,
		Username:   u.Username,
		Visibility: u.Visibility,
		CreatedAt:  u.CreatedAt,
	}
}

// RefreshToken represents a refresh token for JWT authentication
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);not null;index" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// BeforeCreate sets the ID before creating the refresh token
func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the RefreshToken model
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}