package handlers

import (
	"net/http"
	"strconv"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/middleware"
	"github.com/doggyclub/backend/pkg/services"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PostHandler struct {
	postService *services.PostService
	cfg         config.Config
}

func NewPostHandler(db *gorm.DB, redis *redis.Client, cfg config.Config) *PostHandler {
	return &PostHandler{
		postService: services.NewPostService(db, redis, cfg),
		cfg:         cfg,
	}
}

// CreatePost creates a new post
func (h *PostHandler) CreatePost(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req services.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	post, err := h.postService.CreatePost(userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusCreated, post)
}

// GetTimeline returns posts for user's timeline
func (h *PostHandler) GetTimeline(c echo.Context) error {
	userID := middleware.GetUserID(c)

	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offsetStr := c.QueryParam("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	posts, total, err := h.postService.GetTimeline(userID, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"posts":  posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	return c.JSON(http.StatusOK, response)
}

// GetPost returns a single post
func (h *PostHandler) GetPost(c echo.Context) error {
	userID := middleware.GetUserID(c)
	postID := c.Param("postId")

	post, err := h.postService.GetPost(postID, userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, post)
}

// UpdatePost updates a post
func (h *PostHandler) UpdatePost(c echo.Context) error {
	userID := middleware.GetUserID(c)
	postID := c.Param("postId")

	var req services.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	post, err := h.postService.UpdatePost(postID, userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, post)
}

// DeletePost deletes a post
func (h *PostHandler) DeletePost(c echo.Context) error {
	userID := middleware.GetUserID(c)
	postID := c.Param("postId")

	if err := h.postService.DeletePost(postID, userID); err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Post deleted successfully"})
}

// LikePost likes or unlikes a post
func (h *PostHandler) LikePost(c echo.Context) error {
	userID := middleware.GetUserID(c)
	postID := c.Param("postId")

	liked, err := h.postService.LikePost(postID, userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"liked": liked,
		"message": func() string {
			if liked {
				return "Post liked successfully"
			}
			return "Post unliked successfully"
		}(),
	})
}

// AddComment adds a comment to a post
func (h *PostHandler) AddComment(c echo.Context) error {
	userID := middleware.GetUserID(c)
	postID := c.Param("postId")

	var req services.CommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	comment, err := h.postService.AddComment(postID, userID, req)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusCreated, comment)
}

// GetComments returns comments for a post
func (h *PostHandler) GetComments(c echo.Context) error {
	postID := c.Param("postId")

	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offsetStr := c.QueryParam("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	comments, total, err := h.postService.GetComments(postID, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"comments": comments,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	}

	return c.JSON(http.StatusOK, response)
}

// FollowDog follows or unfollows a dog
func (h *PostHandler) FollowDog(c echo.Context) error {
	userID := middleware.GetUserID(c)
	dogID := c.Param("dogId")

	following, err := h.postService.FollowDog(dogID, userID)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"following": following,
		"message": func() string {
			if following {
				return "Dog followed successfully"
			}
			return "Dog unfollowed successfully"
		}(),
	})
}

// SearchPosts searches for posts
func (h *PostHandler) SearchPosts(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Query parameter 'q' is required"})
	}

	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offsetStr := c.QueryParam("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	posts, total, err := h.postService.SearchPosts(query, limit, offset)
	if err != nil {
		status, apiErr := utils.HTTPError(err)
		return c.JSON(status, map[string]interface{}{"error": apiErr})
	}

	response := map[string]interface{}{
		"posts":  posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	return c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers post routes
func (h *PostHandler) RegisterRoutes(e *echo.Echo) {
	posts := e.Group("/api/posts", middleware.AuthMiddleware(h.cfg.JWT))

	// Post management
	posts.POST("", h.CreatePost)
	posts.GET("/timeline", h.GetTimeline)
	posts.GET("/:postId", h.GetPost)
	posts.PUT("/:postId", h.UpdatePost)
	posts.DELETE("/:postId", h.DeletePost)

	// Post interactions
	posts.POST("/:postId/like", h.LikePost)
	posts.POST("/:postId/comments", h.AddComment)
	posts.GET("/:postId/comments", h.GetComments)

	// Following
	posts.POST("/dogs/:dogId/follow", h.FollowDog)

	// Public routes
	postsPublic := e.Group("/api/posts")
	postsPublic.GET("/search", h.SearchPosts, middleware.OptionalAuthMiddleware(h.cfg.JWT))
}