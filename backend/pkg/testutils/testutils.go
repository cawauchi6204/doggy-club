package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/db"
	"github.com/doggyclub/backend/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestDB represents a test database connection
type TestDB struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// TestContext holds common test utilities and data
type TestContext struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Config config.Config
	Echo   *echo.Echo
}

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *TestDB {
	// Use environment variables for test database
	testDBHost := getEnv("TEST_DB_HOST", "localhost")
	testDBPort := getEnv("TEST_DB_PORT", "5432")
	testDBName := getEnv("TEST_DB_NAME", "doggyclub_test")
	testDBUser := getEnv("TEST_DB_USER", "postgres")
	testDBPassword := getEnv("TEST_DB_PASSWORD", "password")

	testConfig := config.Config{
		Database: config.DatabaseConfig{
			Host:     testDBHost,
			Port:     testDBPort,
			User:     testDBUser,
			Password: testDBPassword,
			Name:     testDBName,
			SSLMode:  "disable",
		},
		Redis: config.RedisConfig{
			Host:     getEnv("TEST_REDIS_HOST", "localhost"),
			Port:     getEnv("TEST_REDIS_PORT", "6379"),
			Password: "",
			DB:       1, // Use database 1 for tests
		},
	}

	// Connect to test database
	database, err := db.InitPostgres(testConfig.Database)
	require.NoError(t, err, "Failed to connect to test database")

	// Connect to test Redis
	redisClient, err := db.InitRedis(testConfig.Redis)
	require.NoError(t, err, "Failed to connect to test Redis")

	// Auto-migrate tables
	err = db.AutoMigrate(database)
	require.NoError(t, err, "Failed to run test migrations")

	return &TestDB{
		DB:    database,
		Redis: redisClient,
	}
}

// TeardownTestDB cleans up test database
func (testDB *TestDB) TeardownTestDB(t *testing.T) {
	// Clean up Redis
	testDB.Redis.FlushDB(context.Background())
	testDB.Redis.Close()

	// Clean up database tables
	tables := []string{
		"moderation_actions", "content_filters", "user_suspensions", "blocked_users", "reports",
		"invoices", "payment_methods", "subscriptions", "subscription_features", "subscription_plans",
		"user_devices", "notifications", "notification_preferences",
		"transactions", "user_currencies", "gift_histories", "gifts",
		"follows", "comments", "likes", "posts",
		"encounters", "encounter_settings",
		"vaccination_records", "dogs",
		"refresh_tokens", "users",
		"safety_settings",
	}

	for _, table := range tables {
		testDB.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
	}

	// Close database connection
	if sqlDB, err := testDB.DB.DB(); err == nil {
		sqlDB.Close()
	}
}

// SetupTestContext creates a complete test context
func SetupTestContext(t *testing.T) *TestContext {
	testDB := SetupTestDB(t)

	cfg := config.Config{
		JWT: config.JWTConfig{
			Secret:             "test-secret-key",
			ExpireHours:        24 * time.Hour,
			RefreshExpireHours: 168 * time.Hour,
		},
		Server: config.ServerConfig{
			Environment: "test",
		},
	}

	e := echo.New()

	return &TestContext{
		DB:     testDB.DB,
		Redis:  testDB.Redis,
		Config: cfg,
		Echo:   e,
	}
}

// TeardownTestContext cleans up test context
func (ctx *TestContext) TeardownTestContext(t *testing.T) {
	testDB := &TestDB{DB: ctx.DB, Redis: ctx.Redis}
	testDB.TeardownTestDB(t)
}

// Test data factories

// CreateTestUser creates a test user
func CreateTestUser(t *testing.T, db *gorm.DB) *models.User {
	user := &models.User{
		Username:     fmt.Sprintf("testuser%d", time.Now().UnixNano()),
		Email:        fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
		PasswordHash: "$2a$10$test.hash.password", // Pre-hashed test password
		Visibility:   models.VisibilityPublic,
	}

	err := db.Create(user).Error
	require.NoError(t, err, "Failed to create test user")

	return user
}

// CreateTestDog creates a test dog
func CreateTestDog(t *testing.T, db *gorm.DB, userID string) *models.Dog {
	userUUID, _ := uuid.Parse(userID)
	dog := &models.Dog{
		UserID:   userUUID,
		Name:     "TestDog",
		Breed:    "Golden Retriever",
		Age:      3,
		PhotoURL: "https://example.com/dog.jpg",
		Bio:      "A friendly test dog",
	}

	err := db.Create(dog).Error
	require.NoError(t, err, "Failed to create test dog")

	return dog
}

// CreateTestPost creates a test post
func CreateTestPost(t *testing.T, db *gorm.DB, dogID string) *models.Post {
	dogUUID, _ := uuid.Parse(dogID)
	post := &models.Post{
		DogID:   dogUUID,
		Content: "Test post content",
	}

	err := db.Create(post).Error
	require.NoError(t, err, "Failed to create test post")

	return post
}

// CreateTestGift creates a test gift
func CreateTestGift(t *testing.T, db *gorm.DB, senderDogID, receiverDogID string) *models.Gift {
	senderUUID, _ := uuid.Parse(senderDogID)
	receiverUUID, _ := uuid.Parse(receiverDogID)
	gift := &models.Gift{
		SenderDogID:   senderUUID,
		ReceiverDogID: receiverUUID,
		GiftType:      "bone",
		Message:       "Test gift message",
		SentAt:        time.Now(),
	}

	err := db.Create(gift).Error
	require.NoError(t, err, "Failed to create test gift")

	return gift
}

// CreateTestSubscriptionPlan creates a test subscription plan
func CreateTestSubscriptionPlan(t *testing.T, db *gorm.DB) *models.SubscriptionPlan {
	plan := &models.SubscriptionPlan{
		Name:           "Test Premium",
		Price:          999,
		DurationMonths: 1,
		Features:       `{"unlimited_dogs": true, "premium_support": true}`,
	}

	err := db.Create(plan).Error
	require.NoError(t, err, "Failed to create test subscription plan")

	return plan
}

// Authentication helpers

// GenerateTestJWT generates a test JWT token
func GenerateTestJWT(userID string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(secret))
}

// CreateAuthenticatedRequest creates an HTTP request with authentication
func CreateAuthenticatedRequest(t *testing.T, method, url string, body interface{}, userID string, secret string) *http.Request {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		require.NoError(t, err, "Failed to marshal request body")
	}

	req := httptest.NewRequest(method, url, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Add authentication header
	token, err := GenerateTestJWT(userID, secret)
	require.NoError(t, err, "Failed to generate test JWT")
	req.Header.Set("Authorization", "Bearer "+token)

	return req
}

// HTTP test helpers

// PerformRequest performs an HTTP request and returns the response
func PerformRequest(e *echo.Echo, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// AssertJSONResponse asserts that the response is valid JSON and matches expected status
func AssertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int) map[string]interface{} {
	assert.Equal(t, expectedStatus, rec.Code, "Unexpected status code")
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json", "Response should be JSON")

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal JSON response")

	return response
}

// AssertErrorResponse asserts that the response contains an error
func AssertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedErrorCode string) {
	response := AssertJSONResponse(t, rec, expectedStatus)
	
	assert.Contains(t, response, "error", "Response should contain error field")
	
	if expectedErrorCode != "" {
		errorField, ok := response["error"].(map[string]interface{})
		if ok {
			assert.Equal(t, expectedErrorCode, errorField["code"], "Unexpected error code")
		}
	}
}

// Database test helpers

// CleanupTables truncates specified tables
func CleanupTables(db *gorm.DB, tables ...string) error {
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return err
		}
	}
	return nil
}

// CountRecords counts records in a table
func CountRecords(db *gorm.DB, model interface{}) int64 {
	var count int64
	db.Model(model).Count(&count)
	return count
}

// Test data assertions

// AssertUserExists asserts that a user exists in the database
func AssertUserExists(t *testing.T, db *gorm.DB, userID string) *models.User {
	var user models.User
	err := db.Where("id = ?", userID).First(&user).Error
	require.NoError(t, err, "User should exist in database")
	return &user
}

// AssertUserNotExists asserts that a user does not exist in the database
func AssertUserNotExists(t *testing.T, db *gorm.DB, userID string) {
	var user models.User
	err := db.Where("id = ?", userID).First(&user).Error
	assert.Error(t, err, "User should not exist in database")
	assert.True(t, gorm.ErrRecordNotFound == err, "Should get record not found error")
}

// Time helpers

// TimeWithinDuration asserts that a time is within a duration of now
func AssertTimeWithinDuration(t *testing.T, timeToCheck time.Time, duration time.Duration) {
	now := time.Now()
	diff := now.Sub(timeToCheck)
	if diff < 0 {
		diff = -diff
	}
	assert.True(t, diff <= duration, "Time should be within %v of now, but was %v", duration, diff)
}

// Mocking helpers

// MockRedisClient creates a mock Redis client for testing
type MockRedisClient struct {
	data map[string]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
	}
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	m.data[key] = fmt.Sprintf("%v", value)
	return redis.NewStatusCmd(ctx)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	if value, exists := m.data[key]; exists {
		cmd.SetVal(value)
	} else {
		cmd.SetErr(redis.Nil)
	}
	return cmd
}

// Environment helpers

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Benchmark helpers

// BenchmarkFunc is a helper type for benchmark functions
type BenchmarkFunc func(b *testing.B, ctx *TestContext)

// RunBenchmarkWithContext runs a benchmark with a test context
func RunBenchmarkWithContext(b *testing.B, fn BenchmarkFunc) {
	testDB := SetupTestDB(&testing.T{})
	defer testDB.TeardownTestDB(&testing.T{})

	cfg := config.Config{
		JWT: config.JWTConfig{
			Secret:             "test-secret-key",
			ExpireHours:        24 * time.Hour,
			RefreshExpireHours: 168 * time.Hour,
		},
	}

	ctx := &TestContext{
		DB:     testDB.DB,
		Redis:  testDB.Redis,
		Config: cfg,
		Echo:   echo.New(),
	}

	b.ResetTimer()
	fn(b, ctx)
}

// Load testing helpers

// LoadTestConfig represents configuration for load testing
type LoadTestConfig struct {
	ConcurrentUsers int
	RequestsPerUser int
	RampUpDuration  time.Duration
}

// RunLoadTest runs a simple load test
func RunLoadTest(t *testing.T, config LoadTestConfig, testFunc func()) {
	// This is a basic load test implementation
	// In a real scenario, you'd use more sophisticated tools like k6 or Apache Bench
	
	done := make(chan bool, config.ConcurrentUsers)
	
	for i := 0; i < config.ConcurrentUsers; i++ {
		go func() {
			defer func() { done <- true }()
			
			for j := 0; j < config.RequestsPerUser; j++ {
				testFunc()
				// Small delay to simulate realistic user behavior
				time.Sleep(10 * time.Millisecond)
			}
		}()
		
		// Ramp up delay
		if config.RampUpDuration > 0 {
			time.Sleep(config.RampUpDuration / time.Duration(config.ConcurrentUsers))
		}
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < config.ConcurrentUsers; i++ {
		<-done
	}
}