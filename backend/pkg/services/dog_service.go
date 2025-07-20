package services

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/doggyclub/backend/pkg/models"
)

type DogService struct {
	db *gorm.DB
}

func NewDogService(db *gorm.DB) *DogService {
	return &DogService{
		db: db,
	}
}

// CreateDogRequest represents dog creation request
type CreateDogRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=50"`
	Breed    string `json:"breed" validate:"required,min=1,max=50"`
	Age      int    `json:"age" validate:"required,min=0,max=30"`
	PhotoURL string `json:"photo_url"`
	Bio      string `json:"bio" validate:"max=500"`
}

// UpdateDogRequest represents dog update request
type UpdateDogRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
	Breed    *string `json:"breed,omitempty" validate:"omitempty,min=1,max=50"`
	Age      *int    `json:"age,omitempty" validate:"omitempty,min=0,max=30"`
	PhotoURL *string `json:"photo_url,omitempty"`
	Bio      *string `json:"bio,omitempty" validate:"omitempty,max=500"`
}

// CreateDog creates a new dog profile
func (s *DogService) CreateDog(userID uuid.UUID, req CreateDogRequest) (*models.Dog, error) {
	// Check user exists
	var userCount int64
	if err := s.db.Model(&models.User{}).Where("id = ?", userID).Count(&userCount).Error; err != nil {
		return nil, errors.New("failed to check user existence")
	}
	if userCount == 0 {
		return nil, errors.New("user not found")
	}

	// Create dog
	dog := models.Dog{
		UserID:   userID,
		Name:     req.Name,
		Breed:    req.Breed,
		Age:      req.Age,
		PhotoURL: req.PhotoURL,
		Bio:      req.Bio,
	}

	if err := s.db.Create(&dog).Error; err != nil {
		return nil, errors.New("failed to create dog")
	}

	return &dog, nil
}

// GetUserDogs returns all dogs for a user
func (s *DogService) GetUserDogs(userID uuid.UUID) ([]models.Dog, error) {
	var dogs []models.Dog
	if err := s.db.Where("user_id = ?", userID).Find(&dogs).Error; err != nil {
		return nil, errors.New("failed to get user dogs")
	}

	return dogs, nil
}

// GetDog returns a single dog by ID
func (s *DogService) GetDog(dogID uuid.UUID, userID uuid.UUID) (*models.Dog, error) {
	var dog models.Dog
	query := s.db.Model(&models.Dog{})
	
	// If userID is provided, ensure the dog belongs to the user or is public
	if userID != uuid.Nil {
		query = query.Where("(user_id = ? OR id IN (SELECT id FROM dogs WHERE user_id IN (SELECT id FROM users WHERE visibility = 'public')))", userID)
	}
	
	if err := query.Where("id = ?", dogID).First(&dog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dog not found")
		}
		return nil, errors.New("failed to get dog")
	}

	return &dog, nil
}

// UpdateDog updates dog information
func (s *DogService) UpdateDog(dogID uuid.UUID, userID uuid.UUID, req UpdateDogRequest) (*models.Dog, error) {
	var dog models.Dog
	if err := s.db.Where("id = ? AND user_id = ?", dogID, userID).First(&dog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dog not found")
		}
		return nil, errors.New("failed to find dog")
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Breed != nil {
		updates["breed"] = *req.Breed
	}
	if req.Age != nil {
		updates["age"] = *req.Age
	}
	if req.PhotoURL != nil {
		updates["photo_url"] = *req.PhotoURL
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}

	if len(updates) > 0 {
		if err := s.db.Model(&dog).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update dog")
		}
	}

	// Reload dog
	if err := s.db.Where("id = ?", dogID).First(&dog).Error; err != nil {
		return nil, errors.New("failed to reload dog")
	}

	return &dog, nil
}

// DeleteDog deletes a dog
func (s *DogService) DeleteDog(dogID uuid.UUID, userID uuid.UUID) error {
	var dog models.Dog
	if err := s.db.Where("id = ? AND user_id = ?", dogID, userID).First(&dog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("dog not found")
		}
		return errors.New("failed to find dog")
	}

	if err := s.db.Delete(&dog).Error; err != nil {
		return errors.New("failed to delete dog")
	}

	return nil
}

// SearchPublicDogs searches for public dogs
func (s *DogService) SearchPublicDogs(query string, limit int, offset int) ([]models.Dog, int64, error) {
	var dogs []models.Dog
	var total int64

	// Build search query
	searchQuery := s.db.Model(&models.Dog{}).
		Joins("JOIN users ON dogs.user_id = users.id").
		Where("users.visibility = 'public' AND (dogs.name ILIKE ? OR dogs.breed ILIKE ?)", "%"+query+"%", "%"+query+"%")

	// Count total results
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count dogs")
	}

	// Get paginated results
	if err := searchQuery.Preload("User").Limit(limit).Offset(offset).Find(&dogs).Error; err != nil {
		return nil, 0, errors.New("failed to search dogs")
	}

	return dogs, total, nil
}

// GetPublicDogs returns public dogs with pagination
func (s *DogService) GetPublicDogs(limit int, offset int) ([]models.Dog, int64, error) {
	var dogs []models.Dog
	var total int64

	// Build query for public dogs
	query := s.db.Model(&models.Dog{}).
		Joins("JOIN users ON dogs.user_id = users.id").
		Where("users.visibility = 'public'")

	// Count total results
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count public dogs")
	}

	// Get paginated results
	if err := query.Preload("User").Order("dogs.created_at DESC").Limit(limit).Offset(offset).Find(&dogs).Error; err != nil {
		return nil, 0, errors.New("failed to get public dogs")
	}

	return dogs, total, nil
}