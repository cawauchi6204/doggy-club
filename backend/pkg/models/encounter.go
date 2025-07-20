package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"gorm.io/gorm"
)

// DetectionMethod represents how the encounter was detected
type DetectionMethod string

const (
	DetectionMethodGPS       DetectionMethod = "gps"
	DetectionMethodBluetooth DetectionMethod = "bluetooth"
)

// Encounter represents a meeting between two dogs
type Encounter struct {
	ID               uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Dog1ID           uuid.UUID       `gorm:"type:uuid;not null;index" json:"dog1_id"`
	Dog2ID           uuid.UUID       `gorm:"type:uuid;not null;index" json:"dog2_id"`
	Location         orb.Point       `gorm:"type:geography(POINT)" json:"location"`
	DetectionMethod  DetectionMethod `gorm:"type:varchar(20);not null" json:"detection_method"`
	Timestamp        time.Time       `gorm:"default:CURRENT_TIMESTAMP;index" json:"timestamp"`

	// Relationships
	Dog1 Dog `gorm:"foreignKey:Dog1ID;constraint:OnDelete:CASCADE" json:"dog1,omitempty"`
	Dog2 Dog `gorm:"foreignKey:Dog2ID;constraint:OnDelete:CASCADE" json:"dog2,omitempty"`
}

// BeforeCreate sets the ID before creating the encounter
func (e *Encounter) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the Encounter model
func (Encounter) TableName() string {
	return "encounters"
}

// GetLocationWKT returns the location as WKT string
func (e *Encounter) GetLocationWKT() string {
	return wkt.MarshalString(e.Location)
}

// SetLocationFromCoords sets the location from latitude and longitude
func (e *Encounter) SetLocationFromCoords(lat, lng float64) {
	e.Location = orb.Point{lng, lat}
}

// GetLatLng returns the latitude and longitude
func (e *Encounter) GetLatLng() (float64, float64) {
	return e.Location[1], e.Location[0] // lat, lng
}

// DeviceLocation represents real-time location tracking for cache
type DeviceLocation struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DogID     uuid.UUID `gorm:"type:uuid;not null;index" json:"dog_id"`
	Location  orb.Point `gorm:"type:geography(POINT);index:idx_location,type:gist" json:"location"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"updated_at"`

	// Relationship
	Dog Dog `gorm:"foreignKey:DogID;constraint:OnDelete:CASCADE" json:"dog,omitempty"`
}

// BeforeCreate sets the ID before creating the device location
func (dl *DeviceLocation) BeforeCreate(tx *gorm.DB) error {
	if dl.ID == uuid.Nil {
		dl.ID = uuid.New()
	}
	return nil
}

// BeforeUpdate updates the timestamp
func (dl *DeviceLocation) BeforeUpdate(tx *gorm.DB) error {
	dl.UpdatedAt = time.Now()
	return nil
}

// TableName returns the table name for the DeviceLocation model
func (DeviceLocation) TableName() string {
	return "device_locations"
}

// GetLocationWKT returns the location as WKT string
func (dl *DeviceLocation) GetLocationWKT() string {
	return wkt.MarshalString(dl.Location)
}

// SetLocationFromCoords sets the location from latitude and longitude
func (dl *DeviceLocation) SetLocationFromCoords(lat, lng float64) {
	dl.Location = orb.Point{lng, lat}
}

// GetLatLng returns the latitude and longitude
func (dl *DeviceLocation) GetLatLng() (float64, float64) {
	return dl.Location[1], dl.Location[0] // lat, lng
}