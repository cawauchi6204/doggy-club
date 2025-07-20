package services

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/models"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ModerationService struct {
	db    *gorm.DB
	redis *redis.Client
	cfg   config.Config
}

func NewModerationService(db *gorm.DB, redis *redis.Client, cfg config.Config) *ModerationService {
	return &ModerationService{
		db:    db,
		redis: redis,
		cfg:   cfg,
	}
}

// Request structs
type CreateReportRequest struct {
	ReportedUserID *string `json:"reported_user_id,omitempty"`
	ContentType    string  `json:"content_type" validate:"required,oneof=post comment user dog"`
	ContentID      *string `json:"content_id,omitempty"`
	ReasonCategory string  `json:"reason_category" validate:"required"`
	ReasonDetails  string  `json:"reason_details,omitempty"`
	Description    string  `json:"description" validate:"required,max=1000"`
}

type BlockUserRequest struct {
	BlockedUserID string  `json:"blocked_user_id" validate:"required"`
	Reason        *string `json:"reason,omitempty" validate:"omitempty,max=500"`
}

type UpdateSafetySettingsRequest struct {
	RestrictedMode        *bool    `json:"restricted_mode,omitempty"`
	AllowDirectMessages   *bool    `json:"allow_direct_messages,omitempty"`
	AllowTagging          *bool    `json:"allow_tagging,omitempty"`
	RequireFollowApproval *bool    `json:"require_follow_approval,omitempty"`
	HideFromSearch        *bool    `json:"hide_from_search,omitempty"`
	BlockExplicitContent  *bool    `json:"block_explicit_content,omitempty"`
	MinAge                *int     `json:"min_age,omitempty" validate:"omitempty,min=13,max=100"`
	AllowedLocations      []string `json:"allowed_locations,omitempty"`
}

type ReviewReportRequest struct {
	Status          string  `json:"status" validate:"required,oneof=resolved dismissed escalated"`
	Resolution      *string `json:"resolution,omitempty"`
	ResolutionNotes *string `json:"resolution_notes,omitempty" validate:"omitempty,max=1000"`
}

type SuspendUserRequest struct {
	UserID   string  `json:"user_id" validate:"required"`
	Type     string  `json:"type" validate:"required,oneof=warning temporary_suspension permanent_ban"`
	Reason   string  `json:"reason" validate:"required,max=1000"`
	Duration *int    `json:"duration,omitempty" validate:"omitempty,min=1,max=8760"` // Max 1 year in hours
}

// Content reporting
func (s *ModerationService) CreateReport(reporterID string, req CreateReportRequest) (*models.Report, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(utils.FormatValidationErrors(err))
	}

	// Validate content exists if content report
	if req.ContentType != "user" && req.ContentID != nil {
		if err := s.validateContentExists(req.ContentType, *req.ContentID); err != nil {
			return nil, err
		}
	}

	// Validate reported user exists if user report
	if req.ReportedUserID != nil {
		var reportedUser models.User
		if err := s.db.Where("id = ?", *req.ReportedUserID).First(&reportedUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, utils.ErrNotFound
			}
			return nil, utils.WrapError(err, "failed to validate reported user")
		}

		// Can't report yourself
		if *req.ReportedUserID == reporterID {
			return nil, utils.NewAPIError("SELF_REPORT", "Cannot report yourself", nil)
		}
	}

	// Check for duplicate reports (same reporter, same content within 24 hours)
	var existingReport models.Report
	query := s.db.Where("reporter_id = ? AND content_type = ? AND created_at > ?", 
		reporterID, req.ContentType, time.Now().Add(-24*time.Hour))
	
	if req.ContentID != nil {
		query = query.Where("content_id = ?", *req.ContentID)
	}
	if req.ReportedUserID != nil {
		query = query.Where("reported_user_id = ?", *req.ReportedUserID)
	}

	err := query.First(&existingReport).Error
	if err == nil {
		return nil, utils.NewAPIError("DUPLICATE_REPORT", "You have already reported this content recently", nil)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.WrapError(err, "failed to check for duplicate reports")
	}

	// Determine priority based on reason
	priority := s.determinePriority(req.ReasonCategory)

	// Create report
	report := models.Report{
		ReporterID:     reporterID,
		ReportedUserID: req.ReportedUserID,
		ContentType:    req.ContentType,
		ContentID:      req.ContentID,
		ReasonCategory: req.ReasonCategory,
		ReasonDetails:  req.ReasonDetails,
		Description:    req.Description,
		Status:         models.ReportStatus.Pending,
		Priority:       priority,
	}

	if err := s.db.Create(&report).Error; err != nil {
		return nil, utils.WrapError(err, "failed to create report")
	}

	// Load with relationships
	if err := s.db.Preload("Reporter").Preload("ReportedUser").Where("id = ?", report.ID).First(&report).Error; err != nil {
		return nil, utils.WrapError(err, "failed to reload report")
	}

	// Auto-moderate if high priority
	if priority == models.ReportPriority.Critical {
		go s.autoModerateContent(report)
	}

	return &report, nil
}

// User blocking
func (s *ModerationService) BlockUser(blockerID string, req BlockUserRequest) error {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.NewValidationError(utils.FormatValidationErrors(err))
	}

	// Can't block yourself
	if req.BlockedUserID == blockerID {
		return utils.NewAPIError("SELF_BLOCK", "Cannot block yourself", nil)
	}

	// Check if user exists
	var blockedUser models.User
	if err := s.db.Where("id = ?", req.BlockedUserID).First(&blockedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound
		}
		return utils.WrapError(err, "failed to find user to block")
	}

	// Check if already blocked
	var existingBlock models.BlockedUser
	err := s.db.Where("blocker_id = ? AND blocked_id = ?", blockerID, req.BlockedUserID).First(&existingBlock).Error
	if err == nil {
		return utils.NewAPIError("ALREADY_BLOCKED", "User is already blocked", nil)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.WrapError(err, "failed to check existing block")
	}

	// Create block
	block := models.BlockedUser{
		BlockerID: blockerID,
		BlockedID: req.BlockedUserID,
		Reason:    req.Reason,
	}

	if err := s.db.Create(&block).Error; err != nil {
		return utils.WrapError(err, "failed to block user")
	}

	// Remove any existing follows between the users
	s.db.Where("(follower_dog_id IN (SELECT id FROM dogs WHERE user_id = ?) AND followed_dog_id IN (SELECT id FROM dogs WHERE user_id = ?)) OR (follower_dog_id IN (SELECT id FROM dogs WHERE user_id = ?) AND followed_dog_id IN (SELECT id FROM dogs WHERE user_id = ?))",
		blockerID, req.BlockedUserID, req.BlockedUserID, blockerID).Delete(&models.Follower{})

	return nil
}

func (s *ModerationService) UnblockUser(blockerID string, blockedUserID string) error {
	result := s.db.Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedUserID).Delete(&models.BlockedUser{})
	
	if result.Error != nil {
		return utils.WrapError(result.Error, "failed to unblock user")
	}

	if result.RowsAffected == 0 {
		return utils.ErrNotFound
	}

	return nil
}

func (s *ModerationService) GetBlockedUsers(userID string, limit int, offset int) ([]models.BlockedUser, int64, error) {
	var blockedUsers []models.BlockedUser
	var total int64

	// Count total blocked users
	if err := s.db.Model(&models.BlockedUser{}).Where("blocker_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to count blocked users")
	}

	// Get blocked users with user information
	if err := s.db.Preload("Blocked").Where("blocker_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&blockedUsers).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to get blocked users")
	}

	return blockedUsers, total, nil
}

// Safety settings
func (s *ModerationService) GetSafetySettings(userID string) (*models.SafetySettings, error) {
	var settings models.SafetySettings
	err := s.db.Where("user_id = ?", userID).First(&settings).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create default settings
		settings = models.SafetySettings{
			UserID:                userID,
			RestrictedMode:        false,
			AllowDirectMessages:   true,
			AllowTagging:          true,
			RequireFollowApproval: false,
			HideFromSearch:        false,
			BlockExplicitContent:  true,
		}
		
		if err := s.db.Create(&settings).Error; err != nil {
			return nil, utils.WrapError(err, "failed to create safety settings")
		}
		return &settings, nil
	}
	
	if err != nil {
		return nil, utils.WrapError(err, "failed to get safety settings")
	}

	return &settings, nil
}

func (s *ModerationService) UpdateSafetySettings(userID string, req UpdateSafetySettingsRequest) (*models.SafetySettings, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(utils.FormatValidationErrors(err))
	}

	// Get existing settings
	settings, err := s.GetSafetySettings(userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.RestrictedMode != nil {
		updates["restricted_mode"] = *req.RestrictedMode
	}
	if req.AllowDirectMessages != nil {
		updates["allow_direct_messages"] = *req.AllowDirectMessages
	}
	if req.AllowTagging != nil {
		updates["allow_tagging"] = *req.AllowTagging
	}
	if req.RequireFollowApproval != nil {
		updates["require_follow_approval"] = *req.RequireFollowApproval
	}
	if req.HideFromSearch != nil {
		updates["hide_from_search"] = *req.HideFromSearch
	}
	if req.BlockExplicitContent != nil {
		updates["block_explicit_content"] = *req.BlockExplicitContent
	}
	if req.MinAge != nil {
		updates["min_age"] = *req.MinAge
	}
	if req.AllowedLocations != nil {
		updates["allowed_locations"] = req.AllowedLocations
	}

	if len(updates) > 0 {
		if err := s.db.Model(settings).Updates(updates).Error; err != nil {
			return nil, utils.WrapError(err, "failed to update safety settings")
		}
	}

	// Reload settings
	if err := s.db.Where("user_id = ?", userID).First(settings).Error; err != nil {
		return nil, utils.WrapError(err, "failed to reload safety settings")
	}

	return settings, nil
}

// Admin functions
func (s *ModerationService) GetReports(status string, priority string, limit int, offset int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := s.db.Model(&models.Report{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// Count total reports
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to count reports")
	}

	// Get reports with user information
	if err := query.Preload("Reporter").Preload("ReportedUser").
		Order("priority DESC, created_at DESC").
		Limit(limit).Offset(offset).
		Find(&reports).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to get reports")
	}

	return reports, total, nil
}

func (s *ModerationService) ReviewReport(reviewerID string, reportID string, req ReviewReportRequest) (*models.Report, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(utils.FormatValidationErrors(err))
	}

	// Get report
	var report models.Report
	if err := s.db.Where("id = ?", reportID).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.WrapError(err, "failed to get report")
	}

	// Update report
	now := time.Now()
	updates := map[string]interface{}{
		"status":            req.Status,
		"reviewed_by":       reviewerID,
		"reviewed_at":       now,
		"resolution":        req.Resolution,
		"resolution_notes":  req.ResolutionNotes,
	}

	if err := s.db.Model(&report).Updates(updates).Error; err != nil {
		return nil, utils.WrapError(err, "failed to update report")
	}

	// Create moderation action
	action := models.ModerationAction{
		ModeratorID:     reviewerID,
		TargetType:      report.ContentType,
		TargetID:        *report.ContentID,
		ActionType:      models.ModerationActionType.ReportDismissed,
		Reason:          "Report reviewed",
		RelatedReportID: &reportID,
	}

	if req.Status == models.ReportStatus.Resolved {
		action.ActionType = models.ModerationActionType.ContentApproved
	}

	s.db.Create(&action)

	// Reload report
	if err := s.db.Preload("Reporter").Preload("ReportedUser").Where("id = ?", reportID).First(&report).Error; err != nil {
		return nil, utils.WrapError(err, "failed to reload report")
	}

	return &report, nil
}

func (s *ModerationService) SuspendUser(moderatorID string, req SuspendUserRequest) (*models.UserSuspension, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(utils.FormatValidationErrors(err))
	}

	// Check if user exists
	var user models.User
	if err := s.db.Where("id = ?", req.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.WrapError(err, "failed to find user")
	}

	// Deactivate any existing active suspensions
	s.db.Model(&models.UserSuspension{}).Where("user_id = ? AND is_active = true", req.UserID).Update("is_active", false)

	// Create suspension
	suspension := models.UserSuspension{
		UserID:      req.UserID,
		SuspendedBy: moderatorID,
		Type:        req.Type,
		Reason:      req.Reason,
		Duration:    req.Duration,
		IsActive:    true,
	}

	// Set expiry for temporary suspensions
	if req.Type == models.SuspensionType.TemporarySuspension && req.Duration != nil {
		expiresAt := time.Now().Add(time.Duration(*req.Duration) * time.Hour)
		suspension.ExpiresAt = &expiresAt
	}

	if err := s.db.Create(&suspension).Error; err != nil {
		return nil, utils.WrapError(err, "failed to create suspension")
	}

	// Create moderation action
	var actionType string
	switch req.Type {
	case models.SuspensionType.Warning:
		actionType = models.ModerationActionType.UserWarned
	case models.SuspensionType.TemporarySuspension:
		actionType = models.ModerationActionType.UserSuspended
	case models.SuspensionType.PermanentBan:
		actionType = models.ModerationActionType.UserBanned
	}

	action := models.ModerationAction{
		ModeratorID: moderatorID,
		TargetType:  "user",
		TargetID:    req.UserID,
		ActionType:  actionType,
		Reason:      req.Reason,
	}
	s.db.Create(&action)

	// Load with relationships
	if err := s.db.Preload("User").Preload("SuspendedByUser").Where("id = ?", suspension.ID).First(&suspension).Error; err != nil {
		return nil, utils.WrapError(err, "failed to reload suspension")
	}

	return &suspension, nil
}

// Content filtering
func (s *ModerationService) CheckContentFilter(content string) (bool, string, error) {
	// Get active content filters
	var filters []models.ContentFilter
	if err := s.db.Where("is_active = true").Find(&filters).Error; err != nil {
		return false, "", utils.WrapError(err, "failed to get content filters")
	}

	content = strings.ToLower(content)

	for _, filter := range filters {
		matched := false
		
		switch filter.Type {
		case models.FilterType.Keyword:
			matched = strings.Contains(content, strings.ToLower(filter.Pattern))
		case models.FilterType.Pattern:
			if regex, err := regexp.Compile(filter.Pattern); err == nil {
				matched = regex.MatchString(content)
			}
		}

		if matched {
			return true, filter.Action, nil
		}
	}

	return false, "", nil
}

// Utility functions
func (s *ModerationService) IsUserBlocked(userID string, otherUserID string) (bool, error) {
	var count int64
	if err := s.db.Model(&models.BlockedUser{}).
		Where("(blocker_id = ? AND blocked_id = ?) OR (blocker_id = ? AND blocked_id = ?)",
			userID, otherUserID, otherUserID, userID).
		Count(&count).Error; err != nil {
		return false, utils.WrapError(err, "failed to check if user is blocked")
	}
	return count > 0, nil
}

func (s *ModerationService) IsUserSuspended(userID string) (bool, *models.UserSuspension, error) {
	var suspension models.UserSuspension
	err := s.db.Where("user_id = ? AND is_active = true", userID).First(&suspension).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}
	
	if err != nil {
		return false, nil, utils.WrapError(err, "failed to check user suspension")
	}

	// Check if suspension has expired
	if !suspension.IsCurrentlyActive() {
		// Deactivate expired suspension
		s.db.Model(&suspension).Update("is_active", false)
		return false, nil, nil
	}

	return true, &suspension, nil
}

func (s *ModerationService) validateContentExists(contentType string, contentID string) error {
	switch contentType {
	case "post":
		var post models.Post
		if err := s.db.Where("id = ?", contentID).First(&post).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.ErrNotFound
			}
			return utils.WrapError(err, "failed to validate post")
		}
	case "comment":
		var comment models.Comment
		if err := s.db.Where("id = ?", contentID).First(&comment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.ErrNotFound
			}
			return utils.WrapError(err, "failed to validate comment")
		}
	case "dog":
		var dog models.Dog
		if err := s.db.Where("id = ?", contentID).First(&dog).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.ErrNotFound
			}
			return utils.WrapError(err, "failed to validate dog")
		}
	}
	return nil
}

func (s *ModerationService) determinePriority(reasonCategory string) string {
	switch reasonCategory {
	case models.ReportReason.Violence, models.ReportReason.SelfHarm:
		return models.ReportPriority.Critical
	case models.ReportReason.Harassment, models.ReportReason.InappropriateContent:
		return models.ReportPriority.High
	case models.ReportReason.Spam, models.ReportReason.FakeProfile:
		return models.ReportPriority.Medium
	default:
		return models.ReportPriority.Low
	}
}

func (s *ModerationService) autoModerateContent(report models.Report) {
	// This would implement automatic moderation for critical reports
	// For now, just flag for immediate review
	s.db.Model(&report).Update("priority", models.ReportPriority.Critical)
}