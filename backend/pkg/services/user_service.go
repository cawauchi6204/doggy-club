package services

import (
	"errors"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/models"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserService struct {
	db    *gorm.DB
	redis *redis.Client
	cfg   config.Config
}

func NewUserService(db *gorm.DB, redis *redis.Client, cfg config.Config) *UserService {
	return &UserService{
		db:    db,
		redis: redis,
		cfg:   cfg,
	}
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Nickname     *string `json:"nickname,omitempty" validate:"omitempty,min=2,max=50"`
	ProfileImage *string `json:"profile_image,omitempty"`
}

// UpdatePrivacySettingsRequest represents privacy settings update
type UpdatePrivacySettingsRequest struct {
	ShareLocation    *bool     `json:"share_location,omitempty"`
	ShareProfile     *bool     `json:"share_profile,omitempty"`
	AllowMessages    *bool     `json:"allow_messages,omitempty"`
	BlockedUserIDs   *[]string `json:"blocked_user_ids,omitempty"`
	VisibleFields    *[]string `json:"visible_fields,omitempty"`
}

// User notification preferences removed from this service

// GetProfile returns user profile with dogs
func (s *UserService) GetProfile(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Dogs").Preload("Subscription").Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, utils.WrapError(err, "failed to get user profile")
	}

	// Clear sensitive information
	user.PasswordHash = ""
	
	return &user, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID string, req UpdateProfileRequest) (*models.User, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(utils.FormatValidationErrors(err))
	}

	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, utils.WrapError(err, "failed to find user")
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.ProfileImage != nil {
		updates["profile_image"] = *req.ProfileImage
	}

	if len(updates) > 0 {
		if err := s.db.Model(&user).Updates(updates).Error; err != nil {
			return nil, utils.WrapError(err, "failed to update profile")
		}
	}

	// Reload user with dogs
	if err := s.db.Preload("Dogs").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, utils.WrapError(err, "failed to reload user")
	}

	user.PasswordHash = ""
	return &user, nil
}

// UpdatePrivacySettings updates privacy settings
func (s *UserService) UpdatePrivacySettings(userID string, req UpdatePrivacySettingsRequest) error {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrUserNotFound
		}
		return utils.WrapError(err, "failed to find user")
	}

	// Update privacy settings
	updates := make(map[string]interface{})
	if req.ShareLocation != nil {
		updates["share_location"] = *req.ShareLocation
	}
	if req.ShareProfile != nil {
		updates["share_profile"] = *req.ShareProfile
	}
	if req.AllowMessages != nil {
		updates["allow_messages"] = *req.AllowMessages
	}
	if req.BlockedUserIDs != nil {
		updates["blocked_user_ids"] = *req.BlockedUserIDs
	}
	if req.VisibleFields != nil {
		updates["visible_fields"] = *req.VisibleFields
	}

	if len(updates) > 0 {
		if err := s.db.Model(&user).Updates(updates).Error; err != nil {
			return utils.WrapError(err, "failed to update privacy settings")
		}
	}

	return nil
}

// Notification preferences moved to notification service

// DeleteAccount soft deletes user account
func (s *UserService) DeleteAccount(userID string) error {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrUserNotFound
		}
		return utils.WrapError(err, "failed to find user")
	}

	// Soft delete user (GORM will handle the deleted_at field)
	if err := s.db.Delete(&user).Error; err != nil {
		return utils.WrapError(err, "failed to delete user")
	}

	// Delete all refresh tokens
	s.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{})

	return nil
}

// Currency system removed for simplified schema

// SearchUsers searches for users by nickname or email (admin function)
func (s *UserService) SearchUsers(query string, limit int, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total results
	countQuery := s.db.Model(&models.User{}).Where("nickname ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%")
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to count users")
	}

	// Get paginated results
	if err := s.db.Where("nickname ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(limit).Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to search users")
	}

	// Clear sensitive information
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, total, nil
}