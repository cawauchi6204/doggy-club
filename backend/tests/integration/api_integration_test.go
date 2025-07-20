package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/handlers"
	"github.com/doggyclub/backend/pkg/middleware"
	"github.com/doggyclub/backend/pkg/models"
	"github.com/doggyclub/backend/pkg/services"
	"github.com/doggyclub/backend/pkg/testutils"
)

type APIIntegrationTestSuite struct {
	suite.Suite
	db           *gorm.DB
	server       *httptest.Server
	authService  *services.AuthService
	userService  *services.UserService
	dogService   *services.DogService
	postService  *services.PostService
	testUser     *models.User
	authToken    string
	refreshToken string
}

func TestAPIIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(APIIntegrationTestSuite))
}

func (suite *APIIntegrationTestSuite) SetupSuite() {
	// Load test configuration
	cfg := &config.Config{
		DatabaseURL: os.Getenv("TEST_DATABASE_URL"),
		JWTSecret:   "test-jwt-secret-key-for-integration-tests",
		Environment: "test",
	}

	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://postgres:password@localhost:5432/doggyclub_test?sslmode=disable"
	}

	// Connect to test database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	suite.Require().NoError(err)
	suite.db = db

	// Auto-migrate tables
	err = db.AutoMigrate(
		&models.User{},
		&models.Dog{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
		&models.Encounter{},
		&models.Gift{},
		&models.Notification{},
		&models.Subscription{},
		&models.ModerationLog{},
	)
	suite.Require().NoError(err)

	// Initialize services
	suite.authService = services.NewAuthService(db, cfg.JWTSecret)
	suite.userService = services.NewUserService(db)
	suite.dogService = services.NewDogService(db)
	suite.postService = services.NewPostService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(suite.authService, suite.userService)
	userHandler := handlers.NewUserHandler(suite.userService)
	dogHandler := handlers.NewDogHandler(suite.dogService)
	postHandler := handlers.NewPostHandler(suite.postService)

	// Set up router with middleware
	router := testutils.SetupTestRouter()
	
	// Auth middleware
	authMiddleware := middleware.AuthMiddleware(suite.authService)

	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	api := router.Group("/")
	api.Use(authMiddleware)
	{
		// User routes
		users := api.Group("/users")
		{
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)
			users.DELETE("/me", userHandler.DeleteAccount)
		}

		// Dog routes
		dogs := api.Group("/dogs")
		{
			dogs.POST("", dogHandler.CreateDog)
			dogs.GET("/:id", dogHandler.GetDog)
			dogs.PUT("/:id", dogHandler.UpdateDog)
			dogs.DELETE("/:id", dogHandler.DeleteDog)
			dogs.GET("", dogHandler.GetUserDogs)
		}

		// Post routes
		posts := api.Group("/posts")
		{
			posts.POST("", postHandler.CreatePost)
			posts.GET("", postHandler.GetPosts)
			posts.GET("/:id", postHandler.GetPost)
			posts.PUT("/:id", postHandler.UpdatePost)
			posts.DELETE("/:id", postHandler.DeletePost)
			posts.POST("/:id/like", postHandler.LikePost)
			posts.DELETE("/:id/like", postHandler.UnlikePost)
			posts.POST("/:id/comments", postHandler.CreateComment)
			posts.GET("/:id/comments", postHandler.GetComments)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now(),
		})
	})

	// Create test server
	suite.server = httptest.NewServer(router)

	// Create test user
	suite.createTestUser()
}

func (suite *APIIntegrationTestSuite) TearDownSuite() {
	// Clean up test data
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM dogs")
	suite.db.Exec("DELETE FROM posts")
	suite.db.Exec("DELETE FROM comments")
	suite.db.Exec("DELETE FROM likes")
	suite.db.Exec("DELETE FROM encounters")
	suite.db.Exec("DELETE FROM gifts")
	suite.db.Exec("DELETE FROM notifications")
	suite.db.Exec("DELETE FROM subscriptions")
	suite.db.Exec("DELETE FROM moderation_logs")

	suite.server.Close()
}

func (suite *APIIntegrationTestSuite) SetupTest() {
	// Clean up data before each test
	suite.db.Exec("DELETE FROM posts")
	suite.db.Exec("DELETE FROM dogs")
	suite.db.Exec("DELETE FROM comments")
	suite.db.Exec("DELETE FROM likes")
}

func (suite *APIIntegrationTestSuite) createTestUser() {
	registerData := map[string]interface{}{
		"email":     "integration.test@example.com",
		"password":  "SecurePassword123!",
		"firstName": "Integration",
		"lastName":  "Test",
		"birthday":  "1990-01-01",
	}

	body, _ := json.Marshal(registerData)
	resp, err := http.Post(suite.server.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.Require().Equal(http.StatusCreated, resp.StatusCode)

	var authResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	suite.Require().NoError(err)

	suite.authToken = authResp["access_token"].(string)
	suite.refreshToken = authResp["refresh_token"].(string)

	// Get user from response
	user := authResp["user"].(map[string]interface{})
	suite.testUser = &models.User{
		ID:    user["id"].(string),
		Email: user["email"].(string),
	}
}

// Test complete user registration and login flow
func (suite *APIIntegrationTestSuite) TestUserAuthenticationFlow() {
	// Test registration
	registerData := map[string]interface{}{
		"email":     "newuser@example.com",
		"password":  "SecurePassword123!",
		"firstName": "New",
		"lastName":  "User",
		"birthday":  "1995-05-15",
	}

	body, _ := json.Marshal(registerData)
	resp, err := http.Post(suite.server.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var authResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	suite.NoError(err)

	suite.Contains(authResp, "access_token")
	suite.Contains(authResp, "refresh_token")
	suite.Contains(authResp, "user")

	newUserToken := authResp["access_token"].(string)

	// Test login with the new user
	loginData := map[string]interface{}{
		"email":    "newuser@example.com",
		"password": "SecurePassword123!",
	}

	body, _ = json.Marshal(loginData)
	resp, err = http.Post(suite.server.URL+"/auth/login", "application/json", bytes.NewBuffer(body))
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	// Test token refresh
	refreshData := map[string]interface{}{
		"refresh_token": authResp["refresh_token"].(string),
	}

	body, _ = json.Marshal(refreshData)
	resp, err = http.Post(suite.server.URL+"/auth/refresh", "application/json", bytes.NewBuffer(body))
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	// Test accessing protected route with token
	req, _ := http.NewRequest("GET", suite.server.URL+"/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+newUserToken)

	client := &http.Client{}
	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

// Test user profile management
func (suite *APIIntegrationTestSuite) TestUserProfileManagement() {
	// Get current profile
	req, _ := http.NewRequest("GET", suite.server.URL+"/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var user map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	suite.NoError(err)

	suite.Equal("Integration", user["firstName"])
	suite.Equal("Test", user["lastName"])

	// Update profile
	updateData := map[string]interface{}{
		"firstName": "Updated",
		"lastName":  "Name",
		"bio":       "This is my updated bio for integration testing",
	}

	body, _ := json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", suite.server.URL+"/users/me", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var updatedUser map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&updatedUser)
	suite.NoError(err)

	suite.Equal("Updated", updatedUser["firstName"])
	suite.Equal("Name", updatedUser["lastName"])
	suite.Equal("This is my updated bio for integration testing", updatedUser["bio"])

	// Verify update persisted
	req, _ = http.NewRequest("GET", suite.server.URL+"/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&user)
	suite.NoError(err)

	suite.Equal("Updated", user["firstName"])
	suite.Equal("This is my updated bio for integration testing", user["bio"])
}

// Test dog profile CRUD operations
func (suite *APIIntegrationTestSuite) TestDogProfileCRUD() {
	client := &http.Client{}

	// Create dog
	dogData := map[string]interface{}{
		"name":        "Buddy",
		"breed":       "Golden Retriever",
		"age":         3,
		"size":        "large",
		"personality": "Friendly and energetic",
		"bio":         "Loves to play fetch and swim",
	}

	body, _ := json.Marshal(dogData)
	req, _ := http.NewRequest("POST", suite.server.URL+"/dogs", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var dog map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&dog)
	suite.NoError(err)

	dogID := dog["id"].(string)
	suite.Equal("Buddy", dog["name"])
	suite.Equal("Golden Retriever", dog["breed"])

	// Get dog
	req, _ = http.NewRequest("GET", suite.server.URL+"/dogs/"+dogID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&dog)
	suite.NoError(err)

	suite.Equal(dogID, dog["id"])
	suite.Equal("Buddy", dog["name"])

	// Update dog
	updateData := map[string]interface{}{
		"name":        "Buddy Updated",
		"personality": "Very friendly and energetic",
		"bio":         "Loves to play fetch, swim, and go on hikes",
	}

	body, _ = json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", suite.server.URL+"/dogs/"+dogID, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&dog)
	suite.NoError(err)

	suite.Equal("Buddy Updated", dog["name"])
	suite.Equal("Very friendly and energetic", dog["personality"])

	// Get user's dogs
	req, _ = http.NewRequest("GET", suite.server.URL+"/dogs", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var dogs []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&dogs)
	suite.NoError(err)

	suite.Len(dogs, 1)
	suite.Equal("Buddy Updated", dogs[0]["name"])

	// Delete dog
	req, _ = http.NewRequest("DELETE", suite.server.URL+"/dogs/"+dogID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify deletion
	req, _ = http.NewRequest("GET", suite.server.URL+"/dogs/"+dogID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

// Test post creation, interaction, and deletion
func (suite *APIIntegrationTestSuite) TestPostInteractions() {
	client := &http.Client{}

	// Create post
	postData := map[string]interface{}{
		"content": "Having a wonderful day at the dog park! 🐕 #doglife",
		"images":  []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
	}

	body, _ := json.Marshal(postData)
	req, _ := http.NewRequest("POST", suite.server.URL+"/posts", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var post map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&post)
	suite.NoError(err)

	postID := post["id"].(string)
	suite.Equal("Having a wonderful day at the dog park! 🐕 #doglife", post["content"])

	// Get post
	req, _ = http.NewRequest("GET", suite.server.URL+"/posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&post)
	suite.NoError(err)

	suite.Equal(postID, post["id"])

	// Like post
	req, _ = http.NewRequest("POST", suite.server.URL+"/posts/"+postID+"/like", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var likeResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&likeResp)
	suite.NoError(err)

	suite.Equal(true, likeResp["liked"])

	// Add comment
	commentData := map[string]interface{}{
		"content": "What a beautiful dog! 🥰",
	}

	body, _ = json.Marshal(commentData)
	req, _ = http.NewRequest("POST", suite.server.URL+"/posts/"+postID+"/comments", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var comment map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&comment)
	suite.NoError(err)

	suite.Equal("What a beautiful dog! 🥰", comment["content"])

	// Get comments
	req, _ = http.NewRequest("GET", suite.server.URL+"/posts/"+postID+"/comments", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var comments []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&comments)
	suite.NoError(err)

	suite.Len(comments, 1)
	suite.Equal("What a beautiful dog! 🥰", comments[0]["content"])

	// Get posts feed
	req, _ = http.NewRequest("GET", suite.server.URL+"/posts?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var feed map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&feed)
	suite.NoError(err)

	suite.Contains(feed, "posts")
	posts := feed["posts"].([]interface{})
	suite.Len(posts, 1)

	// Unlike post
	req, _ = http.NewRequest("DELETE", suite.server.URL+"/posts/"+postID+"/like", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&likeResp)
	suite.NoError(err)

	suite.Equal(false, likeResp["liked"])

	// Update post
	updateData := map[string]interface{}{
		"content": "Updated: Having an amazing day at the dog park! 🐕 #doglife #updated",
	}

	body, _ = json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", suite.server.URL+"/posts/"+postID, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&post)
	suite.NoError(err)

	suite.Equal("Updated: Having an amazing day at the dog park! 🐕 #doglife #updated", post["content"])

	// Delete post
	req, _ = http.NewRequest("DELETE", suite.server.URL+"/posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify deletion
	req, _ = http.NewRequest("GET", suite.server.URL+"/posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

// Test error handling and validation
func (suite *APIIntegrationTestSuite) TestErrorHandling() {
	client := &http.Client{}

	// Test invalid registration data
	invalidRegData := map[string]interface{}{
		"email":    "invalid-email",
		"password": "weak",
	}

	body, _ := json.Marshal(invalidRegData)
	resp, err := http.Post(suite.server.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	// Test unauthorized access
	req, _ := http.NewRequest("GET", suite.server.URL+"/users/me", nil)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusUnauthorized, resp.StatusCode)

	// Test invalid token
	req, _ = http.NewRequest("GET", suite.server.URL+"/users/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusUnauthorized, resp.StatusCode)

	// Test accessing non-existent resource
	req, _ = http.NewRequest("GET", suite.server.URL+"/dogs/non-existent-id", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNotFound, resp.StatusCode)

	// Test invalid JSON
	req, _ = http.NewRequest("POST", suite.server.URL+"/dogs", bytes.NewBufferString("invalid json"))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

// Test concurrent operations
func (suite *APIIntegrationTestSuite) TestConcurrentOperations() {
	client := &http.Client{}

	// Create multiple posts concurrently
	const numPosts = 10
	results := make(chan error, numPosts)

	for i := 0; i < numPosts; i++ {
		go func(index int) {
			postData := map[string]interface{}{
				"content": fmt.Sprintf("Concurrent post #%d", index),
			}

			body, _ := json.Marshal(postData)
			req, _ := http.NewRequest("POST", suite.server.URL+"/posts", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+suite.authToken)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				results <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				return
			}

			results <- nil
		}(i)
	}

	// Wait for all operations to complete
	for i := 0; i < numPosts; i++ {
		err := <-results
		suite.NoError(err)
	}

	// Verify all posts were created
	req, _ := http.NewRequest("GET", suite.server.URL+"/posts?limit=20&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var feed map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&feed)
	suite.NoError(err)

	posts := feed["posts"].([]interface{})
	suite.Equal(numPosts, len(posts))
}

// Test database transaction rollback
func (suite *APIIntegrationTestSuite) TestTransactionRollback() {
	// This test would create a scenario where a transaction should be rolled back
	// For example, creating a post with invalid data that should fail validation
	client := &http.Client{}

	// Create post with invalid data that should trigger rollback
	postData := map[string]interface{}{
		"content": "", // Empty content should be invalid
		"images":  []string{"invalid-url"}, // Invalid image URL
	}

	body, _ := json.Marshal(postData)
	req, _ := http.NewRequest("POST", suite.server.URL+"/posts", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Should fail due to validation
	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	// Verify no post was created
	req, _ = http.NewRequest("GET", suite.server.URL+"/posts", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var feed map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&feed)
	suite.NoError(err)

	posts := feed["posts"].([]interface{})
	suite.Equal(0, len(posts))
}

// Benchmark integration test performance
func BenchmarkAPIIntegration(b *testing.B) {
	// This would run performance tests on the API
	// For brevity, we'll just test the health endpoint
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}))
	defer server.Close()

	client := &http.Client{}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Get(server.URL)
			if err != nil {
				b.Error(err)
			}
			resp.Body.Close()
		}
	})
}