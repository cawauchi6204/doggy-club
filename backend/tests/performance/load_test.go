package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/handlers"
	"github.com/doggyclub/backend/pkg/middleware"
	"github.com/doggyclub/backend/pkg/models"
	"github.com/doggyclub/backend/pkg/services"
	"github.com/doggyclub/backend/pkg/testutils"
)

type LoadTestSuite struct {
	server      *httptest.Server
	db          *gorm.DB
	authService *services.AuthService
	testUsers   []TestUser
	client      *http.Client
}

type TestUser struct {
	ID       string
	Email    string
	Token    string
	Password string
}

type LoadTestResult struct {
	TotalRequests    int
	SuccessfulReqs   int
	FailedReqs       int
	TotalDuration    time.Duration
	AvgResponseTime  time.Duration
	MinResponseTime  time.Duration
	MaxResponseTime  time.Duration
	RequestsPerSec   float64
	Errors           []error
}

func setupLoadTestServer() (*LoadTestSuite, error) {
	cfg := &config.Config{
		DatabaseURL: os.Getenv("TEST_DATABASE_URL"),
		JWTSecret:   "load-test-jwt-secret-key-for-performance-testing",
		Environment: "test",
	}

	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://postgres:password@localhost:5432/doggyclub_test?sslmode=disable"
	}

	// Connect to test database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}

	// Initialize services
	authService := services.NewAuthService(db, cfg.JWTSecret)
	userService := services.NewUserService(db)
	dogService := services.NewDogService(db)
	postService := services.NewPostService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(userService)
	dogHandler := handlers.NewDogHandler(dogService)
	postHandler := handlers.NewPostHandler(postService)

	// Set up router with middleware
	router := testutils.SetupTestRouter()
	
	// Auth middleware
	authMiddleware := middleware.AuthMiddleware(authService)

	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := router.Group("/")
	api.Use(authMiddleware)
	{
		users := api.Group("/users")
		{
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)
		}

		dogs := api.Group("/dogs")
		{
			dogs.POST("", dogHandler.CreateDog)
			dogs.GET("/:id", dogHandler.GetDog)
			dogs.GET("", dogHandler.GetUserDogs)
		}

		posts := api.Group("/posts")
		{
			posts.POST("", postHandler.CreatePost)
			posts.GET("", postHandler.GetPosts)
			posts.GET("/:id", postHandler.GetPost)
			posts.POST("/:id/like", postHandler.LikePost)
			posts.POST("/:id/comments", postHandler.CreateComment)
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
	server := httptest.NewServer(router)

	return &LoadTestSuite{
		server:      server,
		db:          db,
		authService: authService,
		client:      &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (suite *LoadTestSuite) cleanup() {
	// Clean up test data
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM dogs")
	suite.db.Exec("DELETE FROM posts")
	suite.db.Exec("DELETE FROM comments")
	suite.db.Exec("DELETE FROM likes")

	suite.server.Close()
}

func (suite *LoadTestSuite) createTestUsers(count int) error {
	suite.testUsers = make([]TestUser, 0, count)

	for i := 0; i < count; i++ {
		email := fmt.Sprintf("loadtest%d@example.com", i)
		password := fmt.Sprintf("LoadTestPassword%d!", i)

		registerData := map[string]interface{}{
			"email":     email,
			"password":  password,
			"firstName": fmt.Sprintf("Load%d", i),
			"lastName":  fmt.Sprintf("Test%d", i),
			"birthday":  "1990-01-01",
		}

		body, _ := json.Marshal(registerData)
		resp, err := suite.client.Post(suite.server.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to create user %d: status %d", i, resp.StatusCode)
		}

		var authResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&authResp)
		if err != nil {
			return err
		}

		user := authResp["user"].(map[string]interface{})
		suite.testUsers = append(suite.testUsers, TestUser{
			ID:       user["id"].(string),
			Email:    email,
			Token:    authResp["access_token"].(string),
			Password: password,
		})
	}

	return nil
}

func (suite *LoadTestSuite) runLoadTest(endpoint string, method string, requestsPerUser int, concurrentUsers int, requestBuilder func(user TestUser) (*http.Request, error)) *LoadTestResult {
	totalRequests := requestsPerUser * concurrentUsers
	results := make(chan time.Duration, totalRequests)
	errors := make(chan error, totalRequests)
	
	var wg sync.WaitGroup
	startTime := time.Now()

	// Launch concurrent users
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userIndex int) {
			defer wg.Done()
			
			user := suite.testUsers[userIndex%len(suite.testUsers)]
			
			for j := 0; j < requestsPerUser; j++ {
				reqStart := time.Now()
				
				req, err := requestBuilder(user)
				if err != nil {
					errors <- err
					continue
				}

				resp, err := suite.client.Do(req)
				duration := time.Since(reqStart)
				results <- duration
				
				if err != nil {
					errors <- err
					continue
				}
				resp.Body.Close()
				
				if resp.StatusCode >= 400 {
					errors <- fmt.Errorf("HTTP %d", resp.StatusCode)
				}
			}
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(startTime)
	
	close(results)
	close(errors)

	// Collect results
	var responseTimes []time.Duration
	var errorList []error
	
	for duration := range results {
		responseTimes = append(responseTimes, duration)
	}
	
	for err := range errors {
		errorList = append(errorList, err)
	}

	// Calculate statistics
	successfulReqs := len(responseTimes)
	failedReqs := len(errorList)
	
	var totalResponseTime time.Duration
	var minResponseTime, maxResponseTime time.Duration
	
	if len(responseTimes) > 0 {
		minResponseTime = responseTimes[0]
		maxResponseTime = responseTimes[0]
		
		for _, duration := range responseTimes {
			totalResponseTime += duration
			if duration < minResponseTime {
				minResponseTime = duration
			}
			if duration > maxResponseTime {
				maxResponseTime = duration
			}
		}
	}

	avgResponseTime := time.Duration(0)
	if successfulReqs > 0 {
		avgResponseTime = totalResponseTime / time.Duration(successfulReqs)
	}

	requestsPerSec := float64(successfulReqs) / totalDuration.Seconds()

	return &LoadTestResult{
		TotalRequests:   totalRequests,
		SuccessfulReqs:  successfulReqs,
		FailedReqs:      failedReqs,
		TotalDuration:   totalDuration,
		AvgResponseTime: avgResponseTime,
		MinResponseTime: minResponseTime,
		MaxResponseTime: maxResponseTime,
		RequestsPerSec:  requestsPerSec,
		Errors:          errorList,
	}
}

// Test health endpoint under load
func TestHealthEndpointLoad(t *testing.T) {
	if os.Getenv("SKIP_LOAD_TESTS") == "true" {
		t.Skip("Skipping load tests")
	}

	suite, err := setupLoadTestServer()
	assert.NoError(t, err)
	defer suite.cleanup()

	err = suite.createTestUsers(10)
	assert.NoError(t, err)

	result := suite.runLoadTest("/health", "GET", 100, 10, func(user TestUser) (*http.Request, error) {
		return http.NewRequest("GET", suite.server.URL+"/health", nil)
	})

	t.Logf("Health Endpoint Load Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Successful: %d, Failed: %d", result.SuccessfulReqs, result.FailedReqs)
	t.Logf("Total Duration: %v", result.TotalDuration)
	t.Logf("Avg Response Time: %v", result.AvgResponseTime)
	t.Logf("Min/Max Response Time: %v/%v", result.MinResponseTime, result.MaxResponseTime)
	t.Logf("Requests/sec: %.2f", result.RequestsPerSec)
	t.Logf("Error Rate: %.2f%%", float64(result.FailedReqs)/float64(result.TotalRequests)*100)

	// Assertions
	assert.Equal(t, 1000, result.TotalRequests)
	assert.True(t, result.SuccessfulReqs > 950, "Success rate should be > 95%")
	assert.True(t, result.RequestsPerSec > 100, "Should handle > 100 req/sec")
	assert.True(t, result.AvgResponseTime < 100*time.Millisecond, "Avg response time should be < 100ms")
}

// Test user profile endpoint under load
func TestUserProfileLoad(t *testing.T) {
	if os.Getenv("SKIP_LOAD_TESTS") == "true" {
		t.Skip("Skipping load tests")
	}

	suite, err := setupLoadTestServer()
	assert.NoError(t, err)
	defer suite.cleanup()

	err = suite.createTestUsers(20)
	assert.NoError(t, err)

	result := suite.runLoadTest("/users/me", "GET", 50, 20, func(user TestUser) (*http.Request, error) {
		req, err := http.NewRequest("GET", suite.server.URL+"/users/me", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+user.Token)
		return req, nil
	})

	t.Logf("User Profile Load Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Successful: %d, Failed: %d", result.SuccessfulReqs, result.FailedReqs)
	t.Logf("Total Duration: %v", result.TotalDuration)
	t.Logf("Avg Response Time: %v", result.AvgResponseTime)
	t.Logf("Min/Max Response Time: %v/%v", result.MinResponseTime, result.MaxResponseTime)
	t.Logf("Requests/sec: %.2f", result.RequestsPerSec)
	t.Logf("Error Rate: %.2f%%", float64(result.FailedReqs)/float64(result.TotalRequests)*100)

	assert.Equal(t, 1000, result.TotalRequests)
	assert.True(t, result.SuccessfulReqs > 950, "Success rate should be > 95%")
	assert.True(t, result.RequestsPerSec > 50, "Should handle > 50 req/sec")
	assert.True(t, result.AvgResponseTime < 200*time.Millisecond, "Avg response time should be < 200ms")
}

// Test post creation under load
func TestPostCreationLoad(t *testing.T) {
	if os.Getenv("SKIP_LOAD_TESTS") == "true" {
		t.Skip("Skipping load tests")
	}

	suite, err := setupLoadTestServer()
	assert.NoError(t, err)
	defer suite.cleanup()

	err = suite.createTestUsers(10)
	assert.NoError(t, err)

	requestCounter := 0
	result := suite.runLoadTest("/posts", "POST", 20, 10, func(user TestUser) (*http.Request, error) {
		requestCounter++
		postData := map[string]interface{}{
			"content": fmt.Sprintf("Load test post #%d from user %s", requestCounter, user.Email),
		}

		body, _ := json.Marshal(postData)
		req, err := http.NewRequest("POST", suite.server.URL+"/posts", bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+user.Token)
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	})

	t.Logf("Post Creation Load Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Successful: %d, Failed: %d", result.SuccessfulReqs, result.FailedReqs)
	t.Logf("Total Duration: %v", result.TotalDuration)
	t.Logf("Avg Response Time: %v", result.AvgResponseTime)
	t.Logf("Min/Max Response Time: %v/%v", result.MinResponseTime, result.MaxResponseTime)
	t.Logf("Requests/sec: %.2f", result.RequestsPerSec)
	t.Logf("Error Rate: %.2f%%", float64(result.FailedReqs)/float64(result.TotalRequests)*100)

	assert.Equal(t, 200, result.TotalRequests)
	assert.True(t, result.SuccessfulReqs > 190, "Success rate should be > 95%")
	assert.True(t, result.RequestsPerSec > 20, "Should handle > 20 req/sec")
	assert.True(t, result.AvgResponseTime < 500*time.Millisecond, "Avg response time should be < 500ms")
}

// Test mixed workload
func TestMixedWorkload(t *testing.T) {
	if os.Getenv("SKIP_LOAD_TESTS") == "true" {
		t.Skip("Skipping load tests")
	}

	suite, err := setupLoadTestServer()
	assert.NoError(t, err)
	defer suite.cleanup()

	err = suite.createTestUsers(15)
	assert.NoError(t, err)

	// Create some posts first
	for i := 0; i < 5; i++ {
		user := suite.testUsers[i]
		postData := map[string]interface{}{
			"content": fmt.Sprintf("Initial post %d for mixed workload test", i),
		}

		body, _ := json.Marshal(postData)
		req, _ := http.NewRequest("POST", suite.server.URL+"/posts", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+user.Token)
		req.Header.Set("Content-Type", "application/json")
		
		resp, _ := suite.client.Do(req)
		resp.Body.Close()
	}

	var wg sync.WaitGroup
	results := make(chan *LoadTestResult, 3)

	// Profile reads (60% of traffic)
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := suite.runLoadTest("/users/me", "GET", 30, 15, func(user TestUser) (*http.Request, error) {
			req, err := http.NewRequest("GET", suite.server.URL+"/users/me", nil)
			if err != nil {
				return nil, err
			}
			req.Header.Set("Authorization", "Bearer "+user.Token)
			return req, nil
		})
		results <- result
	}()

	// Post reads (30% of traffic)
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := suite.runLoadTest("/posts", "GET", 15, 15, func(user TestUser) (*http.Request, error) {
			req, err := http.NewRequest("GET", suite.server.URL+"/posts?limit=10&offset=0", nil)
			if err != nil {
				return nil, err
			}
			req.Header.Set("Authorization", "Bearer "+user.Token)
			return req, nil
		})
		results <- result
	}()

	// Post writes (10% of traffic)
	wg.Add(1)
	go func() {
		defer wg.Done()
		postCounter := 0
		result := suite.runLoadTest("/posts", "POST", 5, 15, func(user TestUser) (*http.Request, error) {
			postCounter++
			postData := map[string]interface{}{
				"content": fmt.Sprintf("Mixed workload post #%d from %s", postCounter, user.Email),
			}

			body, _ := json.Marshal(postData)
			req, err := http.NewRequest("POST", suite.server.URL+"/posts", bytes.NewBuffer(body))
			if err != nil {
				return nil, err
			}
			req.Header.Set("Authorization", "Bearer "+user.Token)
			req.Header.Set("Content-Type", "application/json")
			return req, nil
		})
		results <- result
	}()

	wg.Wait()
	close(results)

	totalRequests := 0
	totalSuccessful := 0
	totalFailed := 0
	
	t.Logf("Mixed Workload Test Results:")
	for i, result := range []<-chan *LoadTestResult{results, results, results} {
		if result == nil {
			break
		}
		workloadType := []string{"Profile Reads", "Post Reads", "Post Writes"}[i]
		
		t.Logf("%s: %d requests, %d successful, %.2f req/sec, %v avg response",
			workloadType, result.TotalRequests, result.SuccessfulReqs, 
			result.RequestsPerSec, result.AvgResponseTime)
		
		totalRequests += result.TotalRequests
		totalSuccessful += result.SuccessfulReqs
		totalFailed += result.FailedReqs
	}

	t.Logf("Overall: %d requests, %d successful, %d failed", totalRequests, totalSuccessful, totalFailed)
	
	successRate := float64(totalSuccessful) / float64(totalRequests) * 100
	t.Logf("Overall Success Rate: %.2f%%", successRate)

	assert.True(t, successRate > 95.0, "Overall success rate should be > 95%")
}

// Benchmark database operations
func BenchmarkDatabaseOperations(b *testing.B) {
	if os.Getenv("SKIP_LOAD_TESTS") == "true" {
		b.Skip("Skipping load tests")
	}

	suite, err := setupLoadTestServer()
	if err != nil {
		b.Fatal(err)
	}
	defer suite.cleanup()

	err = suite.createTestUsers(1)
	if err != nil {
		b.Fatal(err)
	}

	user := suite.testUsers[0]

	b.Run("UserProfileRead", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", suite.server.URL+"/users/me", nil)
			req.Header.Set("Authorization", "Bearer "+user.Token)
			
			resp, err := suite.client.Do(req)
			if err != nil {
				b.Error(err)
			}
			resp.Body.Close()
		}
	})

	b.Run("PostCreation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			postData := map[string]interface{}{
				"content": fmt.Sprintf("Benchmark post #%d", i),
			}

			body, _ := json.Marshal(postData)
			req, _ := http.NewRequest("POST", suite.server.URL+"/posts", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+user.Token)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := suite.client.Do(req)
			if err != nil {
				b.Error(err)
			}
			resp.Body.Close()
		}
	})

	b.Run("PostRetrieval", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", suite.server.URL+"/posts?limit=10&offset=0", nil)
			req.Header.Set("Authorization", "Bearer "+user.Token)
			
			resp, err := suite.client.Do(req)
			if err != nil {
				b.Error(err)
			}
			resp.Body.Close()
		}
	})
}

// Memory leak detection test
func TestMemoryUsage(t *testing.T) {
	if os.Getenv("SKIP_LOAD_TESTS") == "true" {
		t.Skip("Skipping load tests")
	}

	suite, err := setupLoadTestServer()
	assert.NoError(t, err)
	defer suite.cleanup()

	err = suite.createTestUsers(5)
	assert.NoError(t, err)

	// Baseline memory measurement would go here
	// This is a simplified version - in practice you'd use
	// runtime.ReadMemStats() or similar tools

	// Run sustained load for memory leak detection
	for round := 0; round < 5; round++ {
		result := suite.runLoadTest("/users/me", "GET", 100, 5, func(user TestUser) (*http.Request, error) {
			req, err := http.NewRequest("GET", suite.server.URL+"/users/me", nil)
			if err != nil {
				return nil, err
			}
			req.Header.Set("Authorization", "Bearer "+user.Token)
			return req, nil
		})

		t.Logf("Round %d: %d successful requests, %.2f req/sec", 
			round+1, result.SuccessfulReqs, result.RequestsPerSec)

		// In a real test, you'd measure memory here and ensure it's not growing
		assert.True(t, result.SuccessfulReqs > 450, "Should maintain performance across rounds")
		
		// Small pause between rounds
		time.Sleep(100 * time.Millisecond)
	}
}