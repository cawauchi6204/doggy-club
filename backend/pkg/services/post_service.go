package services

import (
	"errors"
	"strings"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/models"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PostService struct {
	db    *gorm.DB
	redis *redis.Client
	cfg   config.Config
}

func NewPostService(db *gorm.DB, redis *redis.Client, cfg config.Config) *PostService {
	return &PostService{
		db:    db,
		redis: redis,
		cfg:   cfg,
	}
}

// CreatePostRequest represents post creation request
type CreatePostRequest struct {
	DogID     string   `json:"dog_id" validate:"required"`
	Content   string   `json:"content" validate:"required,max=1000"`
	MediaUrls []string `json:"media_urls" validate:"max=10"`
	MediaType string   `json:"media_type" validate:"oneof=photo video mixed"`
	Hashtags  []string `json:"hashtags" validate:"max=20"`
	Location  *string  `json:"location" validate:"omitempty,max=200"`
	IsPublic  bool     `json:"is_public"`
}

// UpdatePostRequest represents post update request
type UpdatePostRequest struct {
	Content   *string  `json:"content,omitempty" validate:"omitempty,max=1000"`
	MediaUrls *[]string `json:"media_urls,omitempty" validate:"omitempty,max=10"`
	MediaType *string  `json:"media_type,omitempty" validate:"omitempty,oneof=photo video mixed"`
	Hashtags  *[]string `json:"hashtags,omitempty" validate:"omitempty,max=20"`
	Location  *string  `json:"location,omitempty" validate:"omitempty,max=200"`
	IsPublic  *bool    `json:"is_public,omitempty"`
}

// CommentRequest represents comment request
type CommentRequest struct {
	Content  string  `json:"content" validate:"required,max=500"`
	ParentID *string `json:"parent_id,omitempty"`
}

// CreatePost creates a new post
func (s *PostService) CreatePost(userID string, req CreatePostRequest) (*models.Post, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(utils.FormatValidationErrors(err))
	}

	// Check if dog belongs to user
	var dog models.Dog
	if err := s.db.Where("id = ? AND user_id = ?", req.DogID, userID).First(&dog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.WrapError(err, "failed to find dog")
	}

	// Extract hashtags from content
	hashtags := append(req.Hashtags, extractHashtags(req.Content)...)
	hashtags = removeDuplicates(hashtags)

	dogUUID, err := uuid.Parse(req.DogID)
	if err != nil {
		return nil, errors.New("invalid dog ID format")
	}

	// Use the first media URL as the image URL for simplified model
	var imageURL string
	if len(req.MediaUrls) > 0 {
		imageURL = req.MediaUrls[0]
	}

	post := models.Post{
		DogID:    dogUUID,
		Content:  req.Content,
		ImageURL: imageURL,
	}

	if err := s.db.Create(&post).Error; err != nil {
		return nil, utils.WrapError(err, "failed to create post")
	}

	// Load post with dog information
	if err := s.db.Preload("Dog").Where("id = ?", post.ID).First(&post).Error; err != nil {
		return nil, utils.WrapError(err, "failed to reload post")
	}

	return &post, nil
}

// GetTimeline returns posts for user's timeline
func (s *PostService) GetTimeline(userID string, limit int, offset int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Get followed dogs
	var followedDogIDs []string
	s.db.Model(&models.Follower{}).Where("follower_dog_id IN (SELECT id FROM dogs WHERE user_id = ?)", userID).Pluck("followed_dog_id", &followedDogIDs)

	// Get user's own dogs
	var userDogIDs []string
	s.db.Model(&models.Dog{}).Where("user_id = ?", userID).Pluck("id", &userDogIDs)

	// Combine all dog IDs
	allDogIDs := append(followedDogIDs, userDogIDs...)
	if len(allDogIDs) == 0 {
		return posts, 0, nil
	}

	// Count total posts
	if err := s.db.Model(&models.Post{}).Where("dog_id IN ? AND is_public = true", allDogIDs).Count(&total).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to count posts")
	}

	// Get posts with pagination
	if err := s.db.Preload("Dog").
		Where("dog_id IN ? AND is_public = true", allDogIDs).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to get timeline posts")
	}

	return posts, total, nil
}

// GetPost returns a single post by ID
func (s *PostService) GetPost(postID string, userID string) (*models.Post, error) {
	var post models.Post
	query := s.db.Preload("Dog")

	// Check if user can access this post
	if userID != "" {
		// Check if it's a public post or user owns the dog
		query = query.Where("(is_public = true OR dog_id IN (SELECT id FROM dogs WHERE user_id = ?))", userID)
	} else {
		// Public access only
		query = query.Where("is_public = true")
	}

	if err := query.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.WrapError(err, "failed to get post")
	}

	return &post, nil
}

// UpdatePost updates a post
func (s *PostService) UpdatePost(postID string, userID string, req UpdatePostRequest) (*models.Post, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(utils.FormatValidationErrors(err))
	}

	// Check if post belongs to user's dog
	var post models.Post
	if err := s.db.Joins("JOIN dogs ON posts.dog_id = dogs.id").
		Where("posts.id = ? AND dogs.user_id = ?", postID, userID).
		First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.WrapError(err, "failed to find post")
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.Content != nil {
		updates["content"] = *req.Content
		// Extract hashtags from new content
		if req.Hashtags != nil {
			hashtags := append(*req.Hashtags, extractHashtags(*req.Content)...)
			updates["hashtags"] = removeDuplicates(hashtags)
		} else {
			updates["hashtags"] = extractHashtags(*req.Content)
		}
	}
	if req.MediaUrls != nil {
		updates["media_urls"] = *req.MediaUrls
	}
	if req.MediaType != nil {
		updates["media_type"] = *req.MediaType
	}
	if req.Hashtags != nil && req.Content == nil {
		updates["hashtags"] = *req.Hashtags
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	if len(updates) > 0 {
		if err := s.db.Model(&post).Updates(updates).Error; err != nil {
			return nil, utils.WrapError(err, "failed to update post")
		}
	}

	// Reload post with dog information
	if err := s.db.Preload("Dog").Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, utils.WrapError(err, "failed to reload post")
	}

	return &post, nil
}

// DeletePost deletes a post
func (s *PostService) DeletePost(postID string, userID string) error {
	// Check if post belongs to user's dog
	var post models.Post
	if err := s.db.Joins("JOIN dogs ON posts.dog_id = dogs.id").
		Where("posts.id = ? AND dogs.user_id = ?", postID, userID).
		First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound
		}
		return utils.WrapError(err, "failed to find post")
	}

	if err := s.db.Delete(&post).Error; err != nil {
		return utils.WrapError(err, "failed to delete post")
	}

	return nil
}

// LikePost likes or unlikes a post
func (s *PostService) LikePost(postID string, userID string) (bool, error) {
	// Check if post exists and is accessible
	var post models.Post
	if err := s.db.Where("id = ? AND is_public = true", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, utils.ErrNotFound
		}
		return false, utils.WrapError(err, "failed to find post")
	}

	// Check if already liked
	var existingLike models.Like
	err := s.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&existingLike).Error

	if err == nil {
		// Unlike the post
		if err := s.db.Delete(&existingLike).Error; err != nil {
			return false, utils.WrapError(err, "failed to unlike post")
		}

		// Decrement likes count
		s.db.Model(&post).UpdateColumn("likes_count", gorm.Expr("likes_count - 1"))
		return false, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, utils.WrapError(err, "failed to check existing like")
	}

	// Like the post
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return false, errors.New("invalid post ID format")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, errors.New("invalid user ID format")
	}

	// Get user's dog for the like
	var userDog models.Dog
	if err := s.db.Where("user_id = ?", userUUID).First(&userDog).Error; err != nil {
		return false, errors.New("user dog not found")
	}

	like := models.Like{
		PostID: postUUID,
		DogID:  userDog.ID,
	}

	if err := s.db.Create(&like).Error; err != nil {
		return false, utils.WrapError(err, "failed to like post")
	}

	// Increment likes count
	s.db.Model(&post).UpdateColumn("likes_count", gorm.Expr("likes_count + 1"))
	return true, nil
}

// AddComment adds a comment to a post
func (s *PostService) AddComment(postID string, userID string, req CommentRequest) (*models.Comment, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(utils.FormatValidationErrors(err))
	}

	// Check if post exists and is accessible
	var post models.Post
	if err := s.db.Where("id = ? AND is_public = true", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.WrapError(err, "failed to find post")
	}

	// If it's a reply, check if parent comment exists
	if req.ParentID != nil {
		var parentComment models.Comment
		if err := s.db.Where("id = ? AND post_id = ?", *req.ParentID, postID).First(&parentComment).Error; err != nil {
			return nil, utils.ErrNotFound
		}
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Get user's dog for the comment
	var userDog models.Dog
	if err := s.db.Where("user_id = ?", userUUID).First(&userDog).Error; err != nil {
		return nil, errors.New("user dog not found")
	}

	comment := models.Comment{
		PostID:  postUUID,
		DogID:   userDog.ID,
		Content: req.Content,
	}

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, utils.WrapError(err, "failed to create comment")
	}

	// Increment comments count
	s.db.Model(&post).UpdateColumn("comments_count", gorm.Expr("comments_count + 1"))

	// Load comment with user information
	if err := s.db.Preload("User").Where("id = ?", comment.ID).First(&comment).Error; err != nil {
		return nil, utils.WrapError(err, "failed to reload comment")
	}

	return &comment, nil
}

// GetComments returns comments for a post
func (s *PostService) GetComments(postID string, limit int, offset int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	// Count total comments
	if err := s.db.Model(&models.Comment{}).Where("post_id = ? AND parent_id IS NULL", postID).Count(&total).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to count comments")
	}

	// Get comments with replies
	if err := s.db.Preload("User").Preload("Replies.User").
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&comments).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to get comments")
	}

	return comments, total, nil
}

// FollowDog follows or unfollows a dog
func (s *PostService) FollowDog(dogID string, userID string) (bool, error) {
	// Check if dog exists and is public
	var dog models.Dog
	if err := s.db.Where("id = ? AND is_public = true", dogID).First(&dog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, utils.ErrNotFound
		}
		return false, utils.WrapError(err, "failed to find dog")
	}

	// Can't follow own dog
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, errors.New("invalid user ID format")
	}
	if dog.UserID == userUUID {
		return false, utils.NewAPIError("INVALID_ACTION", "Cannot follow your own dog", nil)
	}

	// Get follower dog
	var followerDog models.Dog
	if err := s.db.Where("user_id = ?", userUUID).First(&followerDog).Error; err != nil {
		return false, errors.New("follower dog not found")
	}

	dogUUID, err := uuid.Parse(dogID)
	if err != nil {
		return false, errors.New("invalid dog ID format")
	}

	// Check if already following
	var existingFollow models.Follower
	err = s.db.Where("follower_dog_id = ? AND followed_dog_id = ?", followerDog.ID, dogUUID).First(&existingFollow).Error

	if err == nil {
		// Unfollow the dog
		if err := s.db.Delete(&existingFollow).Error; err != nil {
			return false, utils.WrapError(err, "failed to unfollow dog")
		}
		return false, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, utils.WrapError(err, "failed to check existing follow")
	}

	// Follow the dog
	follow := models.Follower{
		FollowerDogID: followerDog.ID,
		FollowedDogID: dogUUID,
	}

	if err := s.db.Create(&follow).Error; err != nil {
		return false, utils.WrapError(err, "failed to follow dog")
	}

	return true, nil
}

// SearchPosts searches for posts by content or hashtags
func (s *PostService) SearchPosts(query string, limit int, offset int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	searchQuery := "%" + query + "%"

	// Count total results
	countQuery := s.db.Model(&models.Post{}).Where("is_public = true AND (content ILIKE ? OR ? = ANY(hashtags))", searchQuery, query)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to count posts")
	}

	// Get posts
	if err := s.db.Preload("Dog").
		Where("is_public = true AND (content ILIKE ? OR ? = ANY(hashtags))", searchQuery, query).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, utils.WrapError(err, "failed to search posts")
	}

	return posts, total, nil
}

// extractHashtags extracts hashtags from text content
func extractHashtags(content string) []string {
	var hashtags []string
	words := strings.Fields(content)
	
	for _, word := range words {
		if strings.HasPrefix(word, "#") && len(word) > 1 {
			hashtag := strings.ToLower(strings.TrimPrefix(word, "#"))
			// Remove punctuation
			hashtag = strings.TrimRight(hashtag, ".,!?;:")
			if hashtag != "" {
				hashtags = append(hashtags, hashtag)
			}
		}
	}
	
	return hashtags
}

// removeDuplicates removes duplicate strings from slice
func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}