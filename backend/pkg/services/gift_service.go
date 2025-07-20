package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/doggyclub/backend/pkg/models"
)

type GiftService struct {
	db *gorm.DB
}

func NewGiftService(db *gorm.DB) *GiftService {
	return &GiftService{
		db: db,
	}
}

// SendGiftRequest represents gift sending request
type SendGiftRequest struct {
	SenderDogID   uuid.UUID `json:"sender_dog_id" validate:"required"`
	ReceiverDogID uuid.UUID `json:"receiver_dog_id" validate:"required"`
	GiftType      string    `json:"gift_type" validate:"required"`
	Message       string    `json:"message" validate:"max=200"`
}

// GetAvailableGiftTypes returns available gift types
func (s *GiftService) GetAvailableGiftTypes() []string {
	return []string{"bone", "ball", "treat", "toy", "heart", "star", "diamond"}
}

// SendGift sends a virtual gift between dogs
func (s *GiftService) SendGift(req SendGiftRequest) (*models.Gift, error) {
	// Validate that both dogs exist
	var senderDog, receiverDog models.Dog
	if err := s.db.Where("id = ?", req.SenderDogID).First(&senderDog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("sender dog not found")
		}
		return nil, errors.New("failed to find sender dog")
	}

	if err := s.db.Where("id = ?", req.ReceiverDogID).First(&receiverDog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receiver dog not found")
		}
		return nil, errors.New("failed to find receiver dog")
	}

	// Can't send gift to same dog
	if req.SenderDogID == req.ReceiverDogID {
		return nil, errors.New("cannot send gift to the same dog")
	}

	// Validate gift type
	validGiftTypes := s.GetAvailableGiftTypes()
	isValidGiftType := false
	for _, validType := range validGiftTypes {
		if req.GiftType == validType {
			isValidGiftType = true
			break
		}
	}
	if !isValidGiftType {
		return nil, errors.New("invalid gift type")
	}

	// Create gift record
	gift := models.Gift{
		SenderDogID:   req.SenderDogID,
		ReceiverDogID: req.ReceiverDogID,
		GiftType:      req.GiftType,
		Message:       req.Message,
		SentAt:        time.Now(),
	}

	if err := s.db.Create(&gift).Error; err != nil {
		return nil, errors.New("failed to send gift")
	}

	return &gift, nil
}

// GetSentGifts returns gifts sent by dogs owned by user
func (s *GiftService) GetSentGifts(userID uuid.UUID, limit int, offset int) ([]models.Gift, int64, error) {
	var gifts []models.Gift
	var total int64

	// Get user's dog IDs
	var dogIDs []uuid.UUID
	if err := s.db.Model(&models.Dog{}).Where("user_id = ?", userID).Pluck("id", &dogIDs).Error; err != nil {
		return nil, 0, errors.New("failed to get user dogs")
	}

	if len(dogIDs) == 0 {
		return gifts, 0, nil
	}

	// Count total sent gifts
	if err := s.db.Model(&models.Gift{}).Where("sender_dog_id IN ?", dogIDs).Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count sent gifts")
	}

	// Get sent gifts with receiver dog info
	if err := s.db.Preload("SenderDog").Preload("ReceiverDog").
		Where("sender_dog_id IN ?", dogIDs).
		Order("sent_at DESC").
		Limit(limit).Offset(offset).
		Find(&gifts).Error; err != nil {
		return nil, 0, errors.New("failed to get sent gifts")
	}

	return gifts, total, nil
}

// GetReceivedGifts returns gifts received by user's dogs
func (s *GiftService) GetReceivedGifts(userID uuid.UUID, limit int, offset int) ([]models.Gift, int64, error) {
	var gifts []models.Gift
	var total int64

	// Get user's dog IDs
	var dogIDs []uuid.UUID
	if err := s.db.Model(&models.Dog{}).Where("user_id = ?", userID).Pluck("id", &dogIDs).Error; err != nil {
		return nil, 0, errors.New("failed to get user dogs")
	}

	if len(dogIDs) == 0 {
		return gifts, 0, nil
	}

	// Count total received gifts
	if err := s.db.Model(&models.Gift{}).Where("receiver_dog_id IN ?", dogIDs).Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count received gifts")
	}

	// Get received gifts with sender dog info
	if err := s.db.Preload("SenderDog").Preload("ReceiverDog").
		Where("receiver_dog_id IN ?", dogIDs).
		Order("sent_at DESC").
		Limit(limit).Offset(offset).
		Find(&gifts).Error; err != nil {
		return nil, 0, errors.New("failed to get received gifts")
	}

	return gifts, total, nil
}

// GetGiftsByDogID returns all gifts for a specific dog (sent and received)
func (s *GiftService) GetGiftsByDogID(dogID uuid.UUID, limit int, offset int) ([]models.Gift, int64, error) {
	var gifts []models.Gift
	var total int64

	// Count total gifts
	if err := s.db.Model(&models.Gift{}).
		Where("sender_dog_id = ? OR receiver_dog_id = ?", dogID, dogID).
		Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count gifts")
	}

	// Get gifts with sender and receiver info
	if err := s.db.Preload("SenderDog").Preload("ReceiverDog").
		Where("sender_dog_id = ? OR receiver_dog_id = ?", dogID, dogID).
		Order("sent_at DESC").
		Limit(limit).Offset(offset).
		Find(&gifts).Error; err != nil {
		return nil, 0, errors.New("failed to get gifts")
	}

	return gifts, total, nil
}

// GetPopularGiftTypes returns the most popular gift types
func (s *GiftService) GetPopularGiftTypes(limit int) ([]string, error) {
	type GiftTypeCount struct {
		GiftType string `json:"gift_type"`
		Count    int64  `json:"count"`
	}

	var giftCounts []GiftTypeCount
	if err := s.db.Model(&models.Gift{}).
		Select("gift_type, COUNT(*) as count").
		Group("gift_type").
		Order("count DESC").
		Limit(limit).
		Scan(&giftCounts).Error; err != nil {
		return nil, errors.New("failed to get popular gift types")
	}

	var giftTypes []string
	for _, gc := range giftCounts {
		giftTypes = append(giftTypes, gc.GiftType)
	}

	return giftTypes, nil
}

// DeleteGift deletes a gift (only if sent by one of user's dogs)
func (s *GiftService) DeleteGift(userID uuid.UUID, giftID uuid.UUID) error {
	// Check if the gift was sent by one of the user's dogs
	var gift models.Gift
	if err := s.db.Preload("SenderDog").Where("id = ?", giftID).First(&gift).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("gift not found")
		}
		return errors.New("failed to find gift")
	}

	// Check if user owns the sender dog
	if gift.SenderDog.UserID != userID {
		return errors.New("unauthorized: you can only delete gifts sent by your dogs")
	}

	// Delete the gift
	if err := s.db.Delete(&gift).Error; err != nil {
		return errors.New("failed to delete gift")
	}

	return nil
}