package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"gorm.io/gorm"

	"github.com/doggyclub/backend/pkg/models"
)

type EncounterService struct {
	db *gorm.DB
}

func NewEncounterService(db *gorm.DB) *EncounterService {
	return &EncounterService{
		db: db,
	}
}

// LocationUpdateRequest represents a location update request
type LocationUpdateRequest struct {
	DogID     uuid.UUID `json:"dog_id" validate:"required"`
	Latitude  float64   `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64   `json:"longitude" validate:"required,min=-180,max=180"`
}

// EncounterDetectionRequest represents an encounter detection request
type EncounterDetectionRequest struct {
	DogID    uuid.UUID                  `json:"dog_id" validate:"required"`
	Location orb.Point                  `json:"location"`
	Method   models.DetectionMethod     `json:"method" validate:"required"`
	Metadata map[string]interface{}     `json:"metadata,omitempty"`
}

// UpdateDeviceLocation updates or creates a device location for a dog
func (s *EncounterService) UpdateDeviceLocation(req LocationUpdateRequest) error {
	// Check if dog exists
	var dog models.Dog
	if err := s.db.Where("id = ?", req.DogID).First(&dog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("dog not found")
		}
		return errors.New("failed to find dog")
	}

	// Create point from coordinates
	location := orb.Point{req.Longitude, req.Latitude}

	// Update or create device location
	deviceLocation := models.DeviceLocation{
		DogID:     req.DogID,
		Location:  location,
		UpdatedAt: time.Now(),
	}

	// Use UPSERT to update or insert
	if err := s.db.Where("dog_id = ?", req.DogID).
		Assign(map[string]interface{}{
			"location":   location,
			"updated_at": time.Now(),
		}).
		FirstOrCreate(&deviceLocation).Error; err != nil {
		return errors.New("failed to update device location")
	}

	return nil
}

// DetectEncounters finds potential encounters within a radius
func (s *EncounterService) DetectEncounters(dogID uuid.UUID, radiusMeters float64) ([]models.Encounter, error) {
	// Get current dog's location
	var currentLocation models.DeviceLocation
	if err := s.db.Where("dog_id = ?", dogID).First(&currentLocation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dog location not found")
		}
		return nil, errors.New("failed to get dog location")
	}

	// Find nearby dogs using PostGIS ST_DWithin
	var nearbyDogs []models.DeviceLocation
	if err := s.db.Raw(`
		SELECT dog_id, location, updated_at 
		FROM device_locations 
		WHERE dog_id != ? 
		AND ST_DWithin(
			location::geography, 
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, 
			?
		)
		AND updated_at > NOW() - INTERVAL '1 hour'
	`, dogID, currentLocation.Location[0], currentLocation.Location[1], radiusMeters).
		Scan(&nearbyDogs).Error; err != nil {
		return nil, errors.New("failed to find nearby dogs")
	}

	var encounters []models.Encounter
	for _, nearbyDog := range nearbyDogs {
		// Check if encounter already exists recently (within 30 minutes)
		var existingCount int64
		if err := s.db.Model(&models.Encounter{}).
			Where("((dog1_id = ? AND dog2_id = ?) OR (dog1_id = ? AND dog2_id = ?)) AND timestamp > ?",
				dogID, nearbyDog.DogID, nearbyDog.DogID, dogID, time.Now().Add(-30*time.Minute)).
			Count(&existingCount).Error; err != nil {
			continue // Skip this dog if we can't check
		}

		if existingCount == 0 {
			// Create new encounter
			encounter := models.Encounter{
				Dog1ID:    dogID,
				Dog2ID:    nearbyDog.DogID,
				Location:  currentLocation.Location,
				Timestamp: time.Now(),
			}

			if err := s.db.Create(&encounter).Error; err != nil {
				continue // Skip if creation fails
			}

			encounters = append(encounters, encounter)
		}
	}

	return encounters, nil
}

// CreateBluetoothEncounter creates an encounter detected via Bluetooth
func (s *EncounterService) CreateBluetoothEncounter(req EncounterDetectionRequest) (*models.Encounter, error) {
	// Validate that both dogs exist
	var dog1, dog2 models.Dog
	if err := s.db.Where("id = ?", req.DogID).First(&dog1).Error; err != nil {
		return nil, errors.New("dog not found")
	}

	// For Bluetooth, we need the other dog ID from metadata
	otherDogIDStr, ok := req.Metadata["other_dog_id"].(string)
	if !ok {
		return nil, errors.New("other_dog_id required for Bluetooth encounters")
	}

	otherDogID, err := uuid.Parse(otherDogIDStr)
	if err != nil {
		return nil, errors.New("invalid other_dog_id format")
	}

	if err := s.db.Where("id = ?", otherDogID).First(&dog2).Error; err != nil {
		return nil, errors.New("other dog not found")
	}

	// Check for recent encounter
	var existingCount int64
	if err := s.db.Model(&models.Encounter{}).
		Where("((dog1_id = ? AND dog2_id = ?) OR (dog1_id = ? AND dog2_id = ?)) AND timestamp > ?",
			req.DogID, otherDogID, otherDogID, req.DogID, time.Now().Add(-30*time.Minute)).
		Count(&existingCount).Error; err != nil {
		return nil, errors.New("failed to check existing encounters")
	}

	if existingCount > 0 {
		return nil, errors.New("encounter already recorded recently")
	}

	// Create encounter
	encounter := models.Encounter{
		Dog1ID:    req.DogID,
		Dog2ID:    otherDogID,
		Location:  req.Location,
		Timestamp: time.Now(),
	}

	if err := s.db.Create(&encounter).Error; err != nil {
		return nil, errors.New("failed to create encounter")
	}

	return &encounter, nil
}

// GetDogEncounters returns encounters for a specific dog
func (s *EncounterService) GetDogEncounters(dogID uuid.UUID, limit int, offset int) ([]models.Encounter, int64, error) {
	var encounters []models.Encounter
	var total int64

	// Count total encounters
	if err := s.db.Model(&models.Encounter{}).
		Where("dog1_id = ? OR dog2_id = ?", dogID, dogID).
		Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count encounters")
	}

	// Get paginated encounters with related dogs
	if err := s.db.Preload("Dog1").Preload("Dog2").
		Where("dog1_id = ? OR dog2_id = ?", dogID, dogID).
		Order("timestamp DESC").
		Limit(limit).Offset(offset).
		Find(&encounters).Error; err != nil {
		return nil, 0, errors.New("failed to get encounters")
	}

	return encounters, total, nil
}

// GetNearbyDogs returns dogs within a radius of a given location
func (s *EncounterService) GetNearbyDogs(dogID uuid.UUID, latitude, longitude float64, radiusMeters float64) ([]models.Dog, error) {
	var dogs []models.Dog

	// Find dogs within radius using PostGIS
	if err := s.db.Raw(`
		SELECT DISTINCT d.* 
		FROM dogs d
		JOIN device_locations dl ON d.id = dl.dog_id
		JOIN users u ON d.user_id = u.id
		WHERE d.id != ?
		AND u.visibility = 'public'
		AND ST_DWithin(
			dl.location::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
		AND dl.updated_at > NOW() - INTERVAL '1 hour'
		ORDER BY ST_Distance(
			dl.location::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography
		)
	`, dogID, longitude, latitude, radiusMeters, longitude, latitude).
		Find(&dogs).Error; err != nil {
		return nil, errors.New("failed to find nearby dogs")
	}

	return dogs, nil
}

// GetLocationHistory returns location history for a dog
func (s *EncounterService) GetLocationHistory(dogID uuid.UUID, hours int) ([]models.DeviceLocation, error) {
	var locations []models.DeviceLocation

	// Get location history (this would require storing historical locations)
	// For now, just return current location
	if err := s.db.Where("dog_id = ? AND updated_at > ?", 
		dogID, time.Now().Add(-time.Duration(hours)*time.Hour)).
		Order("updated_at DESC").
		Find(&locations).Error; err != nil {
		return nil, errors.New("failed to get location history")
	}

	return locations, nil
}

// CleanupOldLocations removes location data older than specified duration
func (s *EncounterService) CleanupOldLocations(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	
	if err := s.db.Where("updated_at < ?", cutoff).Delete(&models.DeviceLocation{}).Error; err != nil {
		return errors.New("failed to cleanup old locations")
	}

	return nil
}