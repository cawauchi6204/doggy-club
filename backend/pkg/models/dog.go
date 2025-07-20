package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Dog represents a dog profile in the system
type Dog struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name     string    `gorm:"type:varchar(50);not null" json:"name" validate:"required,max=50"`
	Breed    string    `gorm:"type:varchar(50)" json:"breed"`
	Age      int       `gorm:"type:integer" json:"age" validate:"min=0,max=30"`
	PhotoURL string    `gorm:"type:varchar(255)" json:"photo_url"`
	Bio      string    `gorm:"type:text" json:"bio"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	User         User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Posts        []Post         `gorm:"foreignKey:DogID;constraint:OnDelete:CASCADE" json:"posts,omitempty"`
	Likes        []Like         `gorm:"foreignKey:DogID;constraint:OnDelete:CASCADE" json:"likes,omitempty"`
	Comments     []Comment      `gorm:"foreignKey:DogID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
	Encounters1  []Encounter    `gorm:"foreignKey:Dog1ID;constraint:OnDelete:CASCADE" json:"encounters1,omitempty"`
	Encounters2  []Encounter    `gorm:"foreignKey:Dog2ID;constraint:OnDelete:CASCADE" json:"encounters2,omitempty"`
	SentGifts    []Gift         `gorm:"foreignKey:SenderDogID;constraint:OnDelete:CASCADE" json:"sent_gifts,omitempty"`
	ReceivedGifts []Gift        `gorm:"foreignKey:ReceiverDogID;constraint:OnDelete:CASCADE" json:"received_gifts,omitempty"`
	DeviceLocations []DeviceLocation `gorm:"foreignKey:DogID;constraint:OnDelete:CASCADE" json:"device_locations,omitempty"`
	Followers    []Follower     `gorm:"foreignKey:FollowedDogID;constraint:OnDelete:CASCADE" json:"followers,omitempty"`
	Following    []Follower     `gorm:"foreignKey:FollowerDogID;constraint:OnDelete:CASCADE" json:"following,omitempty"`
}

// BeforeCreate sets the ID before creating the dog
func (d *Dog) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the Dog model
func (Dog) TableName() string {
	return "dogs"
}

// GetAge calculates the age of the dog based on current time
func (d *Dog) GetAge() int {
	return d.Age
}

// IsOwner checks if the given user ID is the owner of this dog
func (d *Dog) IsOwner(userID uuid.UUID) bool {
	return d.UserID == userID
}