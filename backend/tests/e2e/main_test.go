package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	baseURL = "http://localhost:8080"
	timeout = 30 * time.Second
)

type E2ETestSuite struct {
	suite.Suite
	client      *http.Client
	authToken   string
	userID      string
	dogID       string
	postID      string
	encounterID string
}

func TestE2ETestSuite(t *testing.T) {
	// Skip if BASE_URL is not set (for CI environments)
	if os.Getenv("SKIP_E2E_TESTS") == "true" {
		t.Skip("Skipping E2E tests")
	}

	suite.Run(t, new(E2ETestSuite))
}

func (suite *E2ETestSuite) SetupSuite() {
	suite.client = &http.Client{
		Timeout: timeout,
	}

	// Wait for server to be ready
	suite.waitForServer()
}

func (suite *E2ETestSuite) waitForServer() {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := suite.client.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return
		}
		time.Sleep(1 * time.Second)
	}
	suite.T().Fatal("Server did not start within timeout")
}

// Test 1: Health Check
func (suite *E2ETestSuite) TestHealthCheck() {
	resp, err := suite.client.Get(baseURL + "/health")
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var health map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&health)
	suite.NoError(err)
	suite.Equal("healthy", health["status"])
}

// Test 2: User Registration
func (suite *E2ETestSuite) TestUserRegistration() {
	registerData := map[string]interface{}{
		"email":     fmt.Sprintf("test_%d@example.com", time.Now().Unix()),
		"password":  "SecurePassword123!",
		"firstName": "Test",
		"lastName":  "User",
		"birthday":  "1990-01-01",
	}

	body, _ := json.Marshal(registerData)
	resp, err := suite.client.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(body))
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var authResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	suite.NoError(err)

	suite.Contains(authResp, "access_token")
	suite.Contains(authResp, "user")

	// Store auth token for subsequent tests
	suite.authToken = authResp["access_token"].(string)
	user := authResp["user"].(map[string]interface{})
	suite.userID = user["id"].(string)
}

// Test 3: User Login
func (suite *E2ETestSuite) TestUserLogin() {
	// First register a user
	email := fmt.Sprintf("login_test_%d@example.com", time.Now().Unix())
	password := "SecurePassword123!"

	registerData := map[string]interface{}{
		"email":     email,
		"password":  password,
		"firstName": "Login",
		"lastName":  "Test",
		"birthday":  "1990-01-01",
	}

	body, _ := json.Marshal(registerData)
	resp, err := suite.client.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(body))
	suite.NoError(err)
	resp.Body.Close()

	// Now test login
	loginData := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	body, _ = json.Marshal(loginData)
	resp, err = suite.client.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var authResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	suite.NoError(err)

	suite.Contains(authResp, "access_token")
	suite.Contains(authResp, "user")
}

// Test 4: Get User Profile
func (suite *E2ETestSuite) TestGetUserProfile() {
	req, _ := http.NewRequest("GET", baseURL+"/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var user map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	suite.NoError(err)

	suite.Contains(user, "id")
	suite.Contains(user, "email")
	suite.Contains(user, "firstName")
}

// Test 5: Update User Profile
func (suite *E2ETestSuite) TestUpdateUserProfile() {
	updateData := map[string]interface{}{
		"firstName": "Updated",
		"lastName":  "Name",
		"bio":       "This is my updated bio",
	}

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", baseURL+"/users/me", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var user map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	suite.NoError(err)

	suite.Equal("Updated", user["firstName"])
	suite.Equal("This is my updated bio", user["bio"])
}

// Test 6: Create Dog Profile
func (suite *E2ETestSuite) TestCreateDogProfile() {
	dogData := map[string]interface{}{
		"name":        "Buddy",
		"breed":       "Golden Retriever",
		"age":         3,
		"size":        "large",
		"personality": "Friendly and energetic",
		"bio":         "Loves to play fetch and swim",
	}

	body, _ := json.Marshal(dogData)
	req, _ := http.NewRequest("POST", baseURL+"/dogs", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var dog map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&dog)
	suite.NoError(err)

	suite.Contains(dog, "id")
	suite.Equal("Buddy", dog["name"])
	suite.Equal("Golden Retriever", dog["breed"])

	// Store dog ID for subsequent tests
	suite.dogID = dog["id"].(string)
}

// Test 7: Get Dog Profile
func (suite *E2ETestSuite) TestGetDogProfile() {
	req, _ := http.NewRequest("GET", baseURL+"/dogs/"+suite.dogID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var dog map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&dog)
	suite.NoError(err)

	suite.Equal(suite.dogID, dog["id"])
	suite.Equal("Buddy", dog["name"])
}

// Test 8: Create Post
func (suite *E2ETestSuite) TestCreatePost() {
	postData := map[string]interface{}{
		"content": "Having a great day at the dog park! 🐕",
		"images":  []string{"https://example.com/image1.jpg"},
	}

	body, _ := json.Marshal(postData)
	req, _ := http.NewRequest("POST", baseURL+"/posts", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var post map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&post)
	suite.NoError(err)

	suite.Contains(post, "id")
	suite.Equal("Having a great day at the dog park! 🐕", post["content"])

	// Store post ID for subsequent tests
	suite.postID = post["id"].(string)
}

// Test 9: Get Posts Feed
func (suite *E2ETestSuite) TestGetPostsFeed() {
	req, _ := http.NewRequest("GET", baseURL+"/posts?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var feed map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&feed)
	suite.NoError(err)

	suite.Contains(feed, "posts")
	suite.Contains(feed, "total")

	posts := feed["posts"].([]interface{})
	suite.GreaterOrEqual(len(posts), 1)
}

// Test 10: Like Post
func (suite *E2ETestSuite) TestLikePost() {
	req, _ := http.NewRequest("POST", baseURL+"/posts/"+suite.postID+"/like", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	suite.NoError(err)

	suite.Contains(result, "liked")
	suite.Equal(true, result["liked"])
}

// Test 11: Comment on Post
func (suite *E2ETestSuite) TestCommentOnPost() {
	commentData := map[string]interface{}{
		"content": "What a cute dog! 🥰",
	}

	body, _ := json.Marshal(commentData)
	req, _ := http.NewRequest("POST", baseURL+"/posts/"+suite.postID+"/comments", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var comment map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&comment)
	suite.NoError(err)

	suite.Contains(comment, "id")
	suite.Equal("What a cute dog! 🥰", comment["content"])
}

// Test 12: Encounter Detection
func (suite *E2ETestSuite) TestEncounterDetection() {
	encounterData := map[string]interface{}{
		"latitude":  37.7749,
		"longitude": -122.4194,
		"accuracy":  10.0,
	}

	body, _ := json.Marshal(encounterData)
	req, _ := http.NewRequest("POST", baseURL+"/encounters/detect", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	suite.NoError(err)

	suite.Contains(result, "encounters")
	// Should be empty since we're the only user in the test
	encounters := result["encounters"].([]interface{})
	suite.Equal(0, len(encounters))
}

// Test 13: Get User's Encounters
func (suite *E2ETestSuite) TestGetUserEncounters() {
	req, _ := http.NewRequest("GET", baseURL+"/encounters?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	suite.NoError(err)

	suite.Contains(result, "encounters")
	suite.Contains(result, "total")
}

// Test 14: Send Gift (will fail without another user, but tests validation)
func (suite *E2ETestSuite) TestSendGiftValidation() {
	giftData := map[string]interface{}{
		"receiverId": "nonexistent-user-id",
		"giftType":   "bone",
		"message":    "For being such a good dog!",
	}

	body, _ := json.Marshal(giftData)
	req, _ := http.NewRequest("POST", baseURL+"/gifts/send", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Should fail because receiver doesn't exist
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

// Test 15: Get User's Gifts
func (suite *E2ETestSuite) TestGetUserGifts() {
	req, _ := http.NewRequest("GET", baseURL+"/gifts?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	suite.NoError(err)

	suite.Contains(result, "sent")
	suite.Contains(result, "received")
}

// Test 16: Get Notifications
func (suite *E2ETestSuite) TestGetNotifications() {
	req, _ := http.NewRequest("GET", baseURL+"/notifications?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	suite.NoError(err)

	suite.Contains(result, "notifications")
	suite.Contains(result, "total")
	suite.Contains(result, "unreadCount")
}

// Test 17: Update Notification Preferences
func (suite *E2ETestSuite) TestUpdateNotificationPreferences() {
	prefsData := map[string]interface{}{
		"encounters":      true,
		"gifts":          true,
		"posts":          false,
		"comments":       true,
		"likes":          false,
		"quietHoursStart": "22:00",
		"quietHoursEnd":   "08:00",
	}

	body, _ := json.Marshal(prefsData)
	req, _ := http.NewRequest("PUT", baseURL+"/notifications/preferences", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var prefs map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&prefs)
	suite.NoError(err)

	suite.Equal(true, prefs["encounters"])
	suite.Equal(false, prefs["posts"])
}

// Test 18: Subscription Info (should show free tier)
func (suite *E2ETestSuite) TestGetSubscriptionInfo() {
	req, _ := http.NewRequest("GET", baseURL+"/subscriptions/status", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var subscription map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&subscription)
	suite.NoError(err)

	suite.Contains(subscription, "tier")
	suite.Contains(subscription, "status")
	// New user should be on free tier
	suite.Equal("free", subscription["tier"])
}

// Test 19: Unauthorized Access
func (suite *E2ETestSuite) TestUnauthorizedAccess() {
	req, _ := http.NewRequest("GET", baseURL+"/users/me", nil)
	// No authorization header

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
}

// Test 20: Invalid Token
func (suite *E2ETestSuite) TestInvalidToken() {
	req, _ := http.NewRequest("GET", baseURL+"/users/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
}

// Test 21: Rate Limiting (basic check)
func (suite *E2ETestSuite) TestRateLimiting() {
	// Make multiple rapid requests to test rate limiting
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", baseURL+"/health", nil)
		resp, err := suite.client.Do(req)
		suite.NoError(err)
		resp.Body.Close()
		
		if i < 4 {
			suite.Equal(http.StatusOK, resp.StatusCode)
		}
		// Note: Rate limiting might kick in, but we don't assert it
		// as it depends on the specific rate limiting configuration
	}
}

// Test 22: Delete Post
func (suite *E2ETestSuite) TestDeletePost() {
	req, _ := http.NewRequest("DELETE", baseURL+"/posts/"+suite.postID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify post is deleted
	req, _ = http.NewRequest("GET", baseURL+"/posts/"+suite.postID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

// Test 23: Delete Dog Profile
func (suite *E2ETestSuite) TestDeleteDogProfile() {
	req, _ := http.NewRequest("DELETE", baseURL+"/dogs/"+suite.dogID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify dog is deleted
	req, _ = http.NewRequest("GET", baseURL+"/dogs/"+suite.dogID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err = suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

// Test 24: Metrics Endpoint (internal access)
func (suite *E2ETestSuite) TestMetricsEndpoint() {
	req, _ := http.NewRequest("GET", baseURL+"/metrics", nil)

	resp, err := suite.client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Metrics endpoint should be accessible (in test environment)
	// In production, this would be restricted by network/firewall rules
	suite.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)

	// Should contain Prometheus metrics
	suite.Contains(string(body), "http_requests_total")
}

// Cleanup function
func (suite *E2ETestSuite) TearDownSuite() {
	// Additional cleanup if needed
	// Note: In a real E2E test, you might want to clean up test data
	// But for this demo, we'll rely on the test database being reset
}

// Helper function to create authenticated request
func (suite *E2ETestSuite) createAuthenticatedRequest(method, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, baseURL+url, body)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

// Helper function to decode JSON response
func (suite *E2ETestSuite) decodeJSONResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

// Benchmark test for basic operations
func BenchmarkE2EBasicOperations(b *testing.B) {
	if os.Getenv("SKIP_E2E_TESTS") == "true" {
		b.Skip("Skipping E2E benchmarks")
	}

	client := &http.Client{Timeout: timeout}

	// Register a test user
	registerData := map[string]interface{}{
		"email":     fmt.Sprintf("bench_%d@example.com", time.Now().Unix()),
		"password":  "SecurePassword123!",
		"firstName": "Bench",
		"lastName":  "User",
		"birthday":  "1990-01-01",
	}

	body, _ := json.Marshal(registerData)
	resp, err := client.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		b.Fatal(err)
	}
	defer resp.Body.Close()

	var authResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&authResp)
	token := authResp["access_token"].(string)

	b.ResetTimer()

	b.Run("HealthCheck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resp, err := client.Get(baseURL + "/health")
			if err != nil {
				b.Error(err)
			}
			resp.Body.Close()
		}
	})

	b.Run("GetUserProfile", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", baseURL+"/users/me", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := client.Do(req)
			if err != nil {
				b.Error(err)
			}
			resp.Body.Close()
		}
	})

	b.Run("GetPosts", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", baseURL+"/posts?limit=10&offset=0", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := client.Do(req)
			if err != nil {
				b.Error(err)
			}
			resp.Body.Close()
		}
	})
}