package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/models"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type NotificationService struct {
	db               *gorm.DB
	cfg              config.Config
	messagingClient  *messaging.Client
}

func NewNotificationService(db *gorm.DB, cfg config.Config) *NotificationService {
	var messagingClient *messaging.Client
	
	// Initialize Firebase if credentials are provided (simplified for now)
	// In real implementation, would check for Firebase config
	if false { // Disabled for simplified version
		opt := option.WithCredentialsFile("firebase-credentials.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err == nil {
			messagingClient, _ = app.Messaging(context.Background())
		}
	}

	return &NotificationService{
		db:              db,
		cfg:             cfg,
		messagingClient: messagingClient,
	}
}

// RegisterDeviceTokenRequest represents device token registration request
type RegisterDeviceTokenRequest struct {
	DeviceToken string                `json:"device_token" validate:"required"`
	DeviceType  models.DeviceType     `json:"device_type" validate:"required"`
}

// RegisterDeviceToken registers a device token for push notifications
func (s *NotificationService) RegisterDeviceToken(userID string, req RegisterDeviceTokenRequest) (*models.DeviceToken, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Check if device already exists
	var existingDevice models.DeviceToken
	err = s.db.Where("user_id = ? AND token = ?", userUUID, req.DeviceToken).First(&existingDevice).Error
	
	if err == nil {
		// Update existing device
		existingDevice.DeviceType = req.DeviceType
		existingDevice.LastActive = time.Now()
		if err := s.db.Save(&existingDevice).Error; err != nil {
			return nil, errors.New("failed to update device token")
		}
		return &existingDevice, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check existing device")
	}

	// Create new device
	device := models.DeviceToken{
		UserID:     userUUID,
		Token:      req.DeviceToken,
		DeviceType: req.DeviceType,
		LastActive: time.Now(),
	}

	if err := s.db.Create(&device).Error; err != nil {
		return nil, errors.New("failed to create device token")
	}

	return &device, nil
}

// SendNotification sends a push notification to a specific user
func (s *NotificationService) SendNotification(userID string, notificationType models.NotificationType, message string) error {
	if s.messagingClient == nil {
		// If Firebase is not configured, just log the notification
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return errors.New("invalid user ID format")
		}
		notification := models.Notification{
			UserID:  userUUID,
			Type:    notificationType,
			Message: message,
			SentAt:  time.Now(),
		}
		return s.db.Create(&notification).Error
	}

	// Get user's device tokens
	var deviceTokens []models.DeviceToken
	if err := s.db.Where("user_id = ?", userID).Find(&deviceTokens).Error; err != nil {
		return errors.New("failed to get device tokens")
	}

	if len(deviceTokens) == 0 {
		// No device tokens, just save notification
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return errors.New("invalid user ID format")
		}
		notification := models.Notification{
			UserID:  userUUID,
			Type:    notificationType,
			Message: message,
			SentAt:  time.Now(),
		}
		return s.db.Create(&notification).Error
	}

	// Send to all device tokens
	var tokens []string
	for _, device := range deviceTokens {
		tokens = append(tokens, device.Token)
	}

	// Create Firebase message
	firebaseMessage := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: "DoggyClub",
			Body:  message,
		},
		Data: map[string]string{
			"type":    string(notificationType),
			"user_id": userID,
		},
	}

	// Send notification
	ctx := context.Background()
	response, err := s.messagingClient.SendMulticast(ctx, firebaseMessage)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	// Handle failed tokens
	if response.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range response.Responses {
			if !resp.Success {
				failedTokens = append(failedTokens, tokens[idx])
			}
		}
		
		// Remove failed tokens from database
		if len(failedTokens) > 0 {
			s.db.Where("token IN ?", failedTokens).Delete(&models.DeviceToken{})
		}
	}

	// Save notification to database
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}
	notification := models.Notification{
		UserID:  userUUID,
		Type:    notificationType,
		Message: message,
		SentAt:  time.Now(),
	}
	return s.db.Create(&notification).Error
}

// GetUserNotifications returns user's notifications
func (s *NotificationService) GetUserNotifications(userID string, limit int, offset int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	// Count total notifications
	if err := s.db.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count notifications")
	}

	// Get notifications
	if err := s.db.Where("user_id = ?", userID).
		Order("sent_at DESC").
		Limit(limit).Offset(offset).
		Find(&notifications).Error; err != nil {
		return nil, 0, errors.New("failed to get notifications")
	}

	return notifications, total, nil
}

// SendEncounterNotification sends notification about a new encounter
func (s *NotificationService) SendEncounterNotification(userID string, dogName string) error {
	message := fmt.Sprintf("Your dog had an encounter with %s!", dogName)
	return s.SendNotification(userID, models.NotificationTypeEncounter, message)
}

// SendGiftNotification sends notification about receiving a gift
func (s *NotificationService) SendGiftNotification(userID string, giftType string, senderDogName string) error {
	message := fmt.Sprintf("%s sent your dog a %s!", senderDogName, giftType)
	return s.SendNotification(userID, models.NotificationTypeGift, message)
}

// SendLikeNotification sends notification about a post like
func (s *NotificationService) SendLikeNotification(userID string, likerDogName string) error {
	message := fmt.Sprintf("%s liked your post!", likerDogName)
	return s.SendNotification(userID, models.NotificationTypeLike, message)
}