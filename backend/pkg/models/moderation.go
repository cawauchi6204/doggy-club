package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Report represents a user report about content or behavior
type Report struct {
	ID               string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ReporterID       string    `json:"reporter_id" gorm:"type:uuid;not null;index"`
	Reporter         User      `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
	ReportedUserID   *string   `json:"reported_user_id,omitempty" gorm:"type:uuid;index"`
	ReportedUser     *User     `json:"reported_user,omitempty" gorm:"foreignKey:ReportedUserID"`
	ContentType      string    `json:"content_type" gorm:"not null;index"` // post, comment, user, dog
	ContentID        *string   `json:"content_id,omitempty" gorm:"type:uuid;index"`
	ReasonCategory   string    `json:"reason_category" gorm:"not null;index"`
	ReasonDetails    string    `json:"reason_details,omitempty"`
	Description      string    `json:"description" gorm:"not null"`
	Status           string    `json:"status" gorm:"not null;index;default:'pending'"`
	Priority         string    `json:"priority" gorm:"not null;default:'medium'"`
	ReviewedBy       *string   `json:"reviewed_by,omitempty" gorm:"type:uuid"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	Resolution       *string   `json:"resolution,omitempty"`
	ResolutionNotes  *string   `json:"resolution_notes,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// BlockedUser represents a user blocking relationship
type BlockedUser struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BlockerID string    `json:"blocker_id" gorm:"type:uuid;not null;index"`
	Blocker   User      `json:"blocker,omitempty" gorm:"foreignKey:BlockerID"`
	BlockedID string    `json:"blocked_id" gorm:"type:uuid;not null;index"`
	Blocked   User      `json:"blocked,omitempty" gorm:"foreignKey:BlockedID"`
	Reason    *string   `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// UserSuspension represents a user suspension/ban
type UserSuspension struct {
	ID              string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          string     `json:"user_id" gorm:"type:uuid;not null;index"`
	User            User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SuspendedBy     string     `json:"suspended_by" gorm:"type:uuid;not null"`
	SuspendedByUser User       `json:"suspended_by_user,omitempty" gorm:"foreignKey:SuspendedBy"`
	Type            string     `json:"type" gorm:"not null"` // warning, temporary_suspension, permanent_ban
	Reason          string     `json:"reason" gorm:"not null"`
	Duration        *int       `json:"duration,omitempty"` // Duration in hours, null for permanent
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	IsActive        bool       `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ContentFilter represents automated content filtering rules
type ContentFilter struct {
	ID           string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name         string    `json:"name" gorm:"not null"`
	Type         string    `json:"type" gorm:"not null"` // keyword, pattern, ml_model
	Pattern      string    `json:"pattern" gorm:"not null"`
	Action       string    `json:"action" gorm:"not null"` // flag, block, auto_delete
	Severity     string    `json:"severity" gorm:"not null;default:'medium'"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	Description  string    `json:"description,omitempty"`
	CreatedBy    string    `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ModerationAction represents actions taken by moderators
type ModerationAction struct {
	ID           string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ModeratorID  string    `json:"moderator_id" gorm:"type:uuid;not null;index"`
	Moderator    User      `json:"moderator,omitempty" gorm:"foreignKey:ModeratorID"`
	TargetType   string    `json:"target_type" gorm:"not null"` // user, post, comment
	TargetID     string    `json:"target_id" gorm:"type:uuid;not null;index"`
	ActionType   string    `json:"action_type" gorm:"not null;index"`
	Reason       string    `json:"reason" gorm:"not null"`
	Details      *string   `json:"details,omitempty"`
	RelatedReportID *string `json:"related_report_id,omitempty" gorm:"type:uuid"`
	CreatedAt    time.Time `json:"created_at"`
}

// SafetySettings represents user safety and privacy settings
type SafetySettings struct {
	ID                    string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID                string    `json:"user_id" gorm:"type:uuid;unique;not null"`
	User                  User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	RestrictedMode        bool      `json:"restricted_mode" gorm:"default:false"`
	AllowDirectMessages   bool      `json:"allow_direct_messages" gorm:"default:true"`
	AllowTagging          bool      `json:"allow_tagging" gorm:"default:true"`
	RequireFollowApproval bool      `json:"require_follow_approval" gorm:"default:false"`
	HideFromSearch        bool      `json:"hide_from_search" gorm:"default:false"`
	BlockExplicitContent  bool      `json:"block_explicit_content" gorm:"default:true"`
	MinAge                *int      `json:"min_age,omitempty"` // Minimum age of users who can interact
	AllowedLocations      []string  `json:"allowed_locations" gorm:"type:text[]"` // Allowed countries/regions
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// BeforeCreate sets the ID for report
func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate sets the ID for blocked user
func (bu *BlockedUser) BeforeCreate(tx *gorm.DB) error {
	if bu.ID == "" {
		bu.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate sets the ID for user suspension
func (us *UserSuspension) BeforeCreate(tx *gorm.DB) error {
	if us.ID == "" {
		us.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate sets the ID for content filter
func (cf *ContentFilter) BeforeCreate(tx *gorm.DB) error {
	if cf.ID == "" {
		cf.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate sets the ID for moderation action
func (ma *ModerationAction) BeforeCreate(tx *gorm.DB) error {
	if ma.ID == "" {
		ma.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate sets the ID for safety settings
func (ss *SafetySettings) BeforeCreate(tx *gorm.DB) error {
	if ss.ID == "" {
		ss.ID = uuid.New().String()
	}
	return nil
}

// IsActive checks if suspension is currently active
func (us *UserSuspension) IsCurrentlyActive() bool {
	if !us.IsActive {
		return false
	}
	if us.ExpiresAt != nil && time.Now().After(*us.ExpiresAt) {
		return false
	}
	return true
}

// Report status constants
var ReportStatus = struct {
	Pending    string
	Reviewing  string
	Resolved   string
	Dismissed  string
	Escalated  string
}{
	Pending:   "pending",
	Reviewing: "reviewing",
	Resolved:  "resolved",
	Dismissed: "dismissed",
	Escalated: "escalated",
}

// Report reason categories
var ReportReason = struct {
	Spam           string
	Harassment     string
	InappropriateContent string
	FakeProfile    string
	Violence       string
	SelfHarm       string
	IntellectualProperty string
	Privacy        string
	Other          string
}{
	Spam:                 "spam",
	Harassment:           "harassment",
	InappropriateContent: "inappropriate_content",
	FakeProfile:          "fake_profile",
	Violence:             "violence",
	SelfHarm:             "self_harm",
	IntellectualProperty: "intellectual_property",
	Privacy:              "privacy",
	Other:                "other",
}

// Report priority levels
var ReportPriority = struct {
	Low      string
	Medium   string
	High     string
	Critical string
}{
	Low:      "low",
	Medium:   "medium",
	High:     "high",
	Critical: "critical",
}

// Suspension types
var SuspensionType = struct {
	Warning            string
	TemporarySuspension string
	PermanentBan       string
}{
	Warning:             "warning",
	TemporarySuspension: "temporary_suspension",
	PermanentBan:        "permanent_ban",
}

// Content filter types
var FilterType = struct {
	Keyword  string
	Pattern  string
	MLModel  string
}{
	Keyword: "keyword",
	Pattern: "pattern",
	MLModel: "ml_model",
}

// Content filter actions
var FilterAction = struct {
	Flag       string
	Block      string
	AutoDelete string
}{
	Flag:       "flag",
	Block:      "block",
	AutoDelete: "auto_delete",
}

// Moderation action types
var ModerationActionType = struct {
	ContentRemoved     string
	UserWarned         string
	UserSuspended      string
	UserBanned         string
	ContentApproved    string
	ReportDismissed    string
}{
	ContentRemoved:  "content_removed",
	UserWarned:      "user_warned",
	UserSuspended:   "user_suspended",
	UserBanned:      "user_banned",
	ContentApproved: "content_approved",
	ReportDismissed: "report_dismissed",
}