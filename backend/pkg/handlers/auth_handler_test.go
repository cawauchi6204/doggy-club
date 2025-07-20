package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/doggyclub/backend/pkg/models"
	"github.com/doggyclub/backend/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthHandler_Register(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	// Setup handler
	authHandler := NewAuthHandler(ctx.DB, ctx.Redis, ctx.Config)
	authHandler.RegisterRoutes(ctx.Echo)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name: "Valid registration",
			requestBody: `{
				"email": "test@example.com",
				"password": "password123",
				"username": "TestUser"
			}`,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)

				assert.Contains(t, response, "access_token")
				assert.Contains(t, response, "refresh_token")
				assert.Contains(t, response, "user")

				user := response["user"].(map[string]interface{})
				assert.Equal(t, "test@example.com", user["email"])
				assert.Equal(t, "TestUser", user["username"])
			},
		},
		{
			name: "Invalid email format",
			requestBody: `{
				"email": "invalid-email",
				"password": "password123",
				"username": "TestUser"
			}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "Missing required fields",
			requestBody: `{
				"email": "test2@example.com"
			}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"invalid": json}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := testutils.PerformRequest(ctx.Echo, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, rec.Body.String())
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	// Setup handler
	authHandler := NewAuthHandler(ctx.DB, ctx.Redis, ctx.Config)
	authHandler.RegisterRoutes(ctx.Echo)

	// Create a test user first
	user := testutils.CreateTestUser(t, ctx.DB)
	
	// Update the existing user with hashed password using bcrypt directly
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)
	ctx.DB.Model(user).Update("password_hash", string(hashedPassword))

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name: "Valid login",
			requestBody: `{
				"email": "` + user.Email + `",
				"password": "password123"
			}`,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)

				assert.Contains(t, response, "access_token")
				assert.Contains(t, response, "refresh_token")
				assert.Contains(t, response, "user")
			},
		},
		{
			name: "Invalid credentials",
			requestBody: `{
				"email": "` + user.Email + `",
				"password": "wrongpassword"
			}`,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "Nonexistent user",
			requestBody: `{
				"email": "nonexistent@example.com",
				"password": "password123"
			}`,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := testutils.PerformRequest(ctx.Echo, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, rec.Body.String())
			}
		})
	}
}

func TestAuthHandler_GetProfile(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	// Setup handler with auth middleware
	authHandler := NewAuthHandler(ctx.DB, ctx.Redis, ctx.Config)
	authHandler.RegisterRoutes(ctx.Echo)

	// Create test user
	user := testutils.CreateTestUser(t, ctx.DB)

	tests := []struct {
		name           string
		withAuth       bool
		userID         string
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:           "Valid authenticated request",
			withAuth:       true,
			userID:         user.ID.String(),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)

				assert.Contains(t, response, "id")
				assert.Contains(t, response, "email")
				assert.Contains(t, response, "username")
				assert.Equal(t, user.ID.String(), response["id"])
			},
		},
		{
			name:           "Unauthenticated request",
			withAuth:       false,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request

			if tt.withAuth {
				req = testutils.CreateAuthenticatedRequest(t, http.MethodGet, "/api/auth/profile", nil, tt.userID, ctx.Config.JWT.Secret)
			} else {
				req = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
			}

			rec := testutils.PerformRequest(ctx.Echo, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, rec.Body.String())
			}
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	// Setup handler
	authHandler := NewAuthHandler(ctx.DB, ctx.Redis, ctx.Config)
	authHandler.RegisterRoutes(ctx.Echo)

	// Create test user and refresh token
	user := testutils.CreateTestUser(t, ctx.DB)
	
	// Create refresh token manually
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     "test-refresh-token-" + user.ID.String(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err := ctx.DB.Create(refreshTokenModel).Error
	require.NoError(t, err)
	refreshToken := refreshTokenModel.Token

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name: "Valid refresh token",
			requestBody: `{
				"refresh_token": "` + refreshToken + `"
			}`,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)

				assert.Contains(t, response, "access_token")
				assert.Contains(t, response, "refresh_token")
				assert.Contains(t, response, "user")
			},
		},
		{
			name: "Invalid refresh token",
			requestBody: `{
				"refresh_token": "invalid-token"
			}`,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := testutils.PerformRequest(ctx.Echo, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, rec.Body.String())
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	// Setup handler
	authHandler := NewAuthHandler(ctx.DB, ctx.Redis, ctx.Config)
	authHandler.RegisterRoutes(ctx.Echo)

	// Create test user and refresh token
	user := testutils.CreateTestUser(t, ctx.DB)
	
	// Create refresh token manually
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     "test-logout-refresh-token-" + user.ID.String(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	err := ctx.DB.Create(refreshTokenModel).Error
	require.NoError(t, err)
	refreshToken := refreshTokenModel.Token

	// Test logout
	requestBody := `{
		"refresh_token": "` + refreshToken + `"
	}`

	req := testutils.CreateAuthenticatedRequest(t, http.MethodPost, "/api/auth/logout", strings.NewReader(requestBody), user.ID.String(), ctx.Config.JWT.Secret)
	req.Header.Set("Content-Type", "application/json")

	rec := testutils.PerformRequest(ctx.Echo, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "message")
}

// Integration test for complete auth flow
func TestAuthHandler_CompleteAuthFlow(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	// Setup handler
	authHandler := NewAuthHandler(ctx.DB, ctx.Redis, ctx.Config)
	authHandler.RegisterRoutes(ctx.Echo)

	// Step 1: Register
	registerBody := `{
		"email": "flow@example.com",
		"password": "password123",
		"username": "FlowUser"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(registerBody))
	req.Header.Set("Content-Type", "application/json")
	rec := testutils.PerformRequest(ctx.Echo, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var registerResponse map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &registerResponse)
	require.NoError(t, err)

	accessToken := registerResponse["access_token"].(string)
	refreshToken := registerResponse["refresh_token"].(string)

	// Step 2: Access protected resource
	req = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	rec = testutils.PerformRequest(ctx.Echo, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// Step 3: Refresh token
	refreshBody := `{
		"refresh_token": "` + refreshToken + `"
	}`

	req = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", strings.NewReader(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	rec = testutils.PerformRequest(ctx.Echo, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var refreshResponse map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &refreshResponse)
	require.NoError(t, err)

	newAccessToken := refreshResponse["access_token"].(string)
	newRefreshToken := refreshResponse["refresh_token"].(string)

	// Verify new tokens are different
	assert.NotEqual(t, accessToken, newAccessToken)
	assert.NotEqual(t, refreshToken, newRefreshToken)

	// Step 4: Logout
	logoutBody := `{
		"refresh_token": "` + newRefreshToken + `"
	}`

	req = httptest.NewRequest(http.MethodPost, "/api/auth/logout", strings.NewReader(logoutBody))
	req.Header.Set("Authorization", "Bearer "+newAccessToken)
	req.Header.Set("Content-Type", "application/json")
	rec = testutils.PerformRequest(ctx.Echo, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// Step 5: Try to use old refresh token (should fail)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", strings.NewReader(logoutBody))
	req.Header.Set("Content-Type", "application/json")
	rec = testutils.PerformRequest(ctx.Echo, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// Note: Rate limiting test disabled since cache middleware is disabled