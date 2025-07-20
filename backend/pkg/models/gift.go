package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Gift represents a virtual gift sent between dogs
type Gift struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SenderDogID    uuid.UUID `gorm:"type:uuid;not null;index" json:"sender_dog_id"`
	ReceiverDogID  uuid.UUID `gorm:"type:uuid;not null;index" json:"receiver_dog_id"`
	GiftType       string    `gorm:"type:varchar(50);not null" json:"gift_type"`
	Message        string    `gorm:"type:text" json:"message"`
	SentAt         time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"sent_at"`

	// Relationships
	SenderDog   Dog `gorm:"foreignKey:SenderDogID;constraint:OnDelete:CASCADE" json:"sender_dog,omitempty"`
	ReceiverDog Dog `gorm:"foreignKey:ReceiverDogID;constraint:OnDelete:CASCADE" json:"receiver_dog,omitempty"`
}

// BeforeCreate sets the ID before creating the gift
func (g *Gift) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the Gift model
func (Gift) TableName() string {
	return "gifts"
}

// Available gift types
var GiftTypes = []string{
	"bone",
	"ball",
	"treat",
	"toy",
	"heart",
	"star",
	"crown",
	"flower",
	"stick",
	"frisbee",
}

// IsValidGiftType checks if the gift type is valid
func IsValidGiftType(giftType string) bool {
	for _, validType := range GiftTypes {
		if validType == giftType {
			return true
		}
	}
	return false
}