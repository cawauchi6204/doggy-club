package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/doggyclub/backend/pkg/models"
	"github.com/doggyclub/backend/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	authService := NewAuthService(ctx.DB, ctx.Config.JWT.Secret)

	tests := []struct {
		name        string
		req         RegisterRequest
		expectError bool
		errorCode   string
	}{
		{
			name: "Valid registration",
			req: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Username: "TestUser",
			},
			expectError: false,
		},
		{
			name: "Invalid email",
			req: RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
				Username: "TestUser",
			},
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
		{
			name: "Password too short",
			req: RegisterRequest{
				Email:    "test2@example.com",
				Password: "123",
				Username: "TestUser",
			},
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
		{
			name: "Missing username",
			req: RegisterRequest{
				Email:    "test3@example.com",
				Password: "password123",
				Username: "",
			},
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Register(tt.req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				require.NoError(t, err)
				require.NotNil(t, response)

				// Verify response structure
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.NotEmpty(t, response.User.ID)
				assert.Equal(t, tt.req.Email, response.User.Email)
				assert.Equal(t, tt.req.Username, response.User.Username)

				// Verify user exists in database
				user := testutils.AssertUserExists(t, ctx.DB, response.User.ID.String())
				assert.Equal(t, tt.req.Email, user.Email)
				assert.Equal(t, tt.req.Username, user.Username)

				// Verify password is hashed
				err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(tt.req.Password))
				assert.NoError(t, err, "Password should be correctly hashed")

				// Verify refresh token exists
				var refreshToken models.RefreshToken
				err = ctx.DB.Where("user_id = ?", user.ID).First(&refreshToken).Error
				assert.NoError(t, err, "Refresh token should exist")
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	authService := NewAuthService(ctx.DB, ctx.Config.JWT.Secret)

	// Create a test user first
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &models.User{
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Username:     "TestUser",
	}
	err = ctx.DB.Create(user).Error
	require.NoError(t, err)

	tests := []struct {
		name        string
		req         LoginRequest
		expectError bool
		errorCode   string
	}{
		{
			name: "Valid login",
			req: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "Invalid email",
			req: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectError: true,
			errorCode:   "INVALID_CREDENTIALS",
		},
		{
			name: "Invalid password",
			req: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			expectError: true,
			errorCode:   "INVALID_CREDENTIALS",
		},
		{
			name: "Empty email",
			req: LoginRequest{
				Email:    "",
				Password: "password123",
			},
			expectError: true,
			errorCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				require.NoError(t, err)
				require.NotNil(t, response)

				// Verify response structure
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.Equal(t, user.ID, response.User.ID)
				assert.Equal(t, user.Email, response.User.Email)
			}
		})
	}
}

// Note: RefreshToken functionality is handled at the handler level
// This test is removed since the service doesn't expose this method

// Note: Logout functionality is handled at the handler level
// This test is removed since the service doesn't expose this method

func TestAuthService_ValidateToken(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	authService := NewAuthService(ctx.DB, ctx.Config.JWT.Secret)
	user := testutils.CreateTestUser(t, ctx.DB)

	// Generate a valid token
	validToken, err := testutils.GenerateTestJWT(user.ID.String(), ctx.Config.JWT.Secret)
	require.NoError(t, err)

	// Generate an invalid token with wrong secret
	invalidToken, err := testutils.GenerateTestJWT(user.ID.String(), "wrong-secret")
	require.NoError(t, err)

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "Valid token",
			token:       validToken,
			expectError: false,
		},
		{
			name:        "Invalid token",
			token:       invalidToken,
			expectError: true,
		},
		{
			name:        "Malformed token",
			token:       "invalid.token.format",
			expectError: true,
		},
		{
			name:        "Empty token",
			token:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := authService.ValidateToken(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, user.ID, userID)
			}
		})
	}
}

func TestAuthService_DuplicateEmailRegistration(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	authService := NewAuthService(ctx.DB, ctx.Config.JWT.Secret)

	// Create first user
	req1 := RegisterRequest{
		Email:    "duplicate@example.com",
		Password: "password123",
		Username: "User1",
	}

	response1, err := authService.Register(req1)
	require.NoError(t, err)
	require.NotNil(t, response1)

	// Try to create second user with same email
	req2 := RegisterRequest{
		Email:    "duplicate@example.com",
		Password: "password456",
		Username: "User2",
	}

	response2, err := authService.Register(req2)
	assert.Error(t, err)
	assert.Nil(t, response2)

	// Verify only one user exists
	var userCount int64
	ctx.DB.Model(&models.User{}).Where("email = ?", "duplicate@example.com").Count(&userCount)
	assert.Equal(t, int64(1), userCount)
}

func TestAuthService_TokenExpiration(t *testing.T) {
	ctx := testutils.SetupTestContext(t)
	defer ctx.TeardownTestContext(t)

	// Set short expiration for testing
	ctx.Config.JWT.ExpireHours = 1 * time.Millisecond
	authService := NewAuthService(ctx.DB, ctx.Config.JWT.Secret)

	user := testutils.CreateTestUser(t, ctx.DB)

	// Generate token
	token, err := authService.generateAccessToken(user.ID, user.Email)
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(2 * time.Millisecond)

	// Try to validate expired token
	userID, err := authService.ValidateToken(token)
	assert.Error(t, err)
	assert.Empty(t, userID)
}

// Benchmark tests
func BenchmarkAuthService_Register(b *testing.B) {
	testutils.RunBenchmarkWithContext(b, func(b *testing.B, ctx *testutils.TestContext) {
		authService := NewAuthService(ctx.DB, ctx.Config.JWT.Secret)

		for i := 0; i < b.N; i++ {
			req := RegisterRequest{
				Email:    fmt.Sprintf("bench%d@example.com", i),
				Password: "password123",
				Username: fmt.Sprintf("User%d", i),
			}

			_, err := authService.Register(req)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAuthService_Login(b *testing.B) {
	testutils.RunBenchmarkWithContext(b, func(b *testing.B, ctx *testutils.TestContext) {
		authService := NewAuthService(ctx.DB, ctx.Config.JWT.Secret)

		// Create test user
		user := testutils.CreateTestUser(&testing.T{}, ctx.DB)
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		ctx.DB.Model(user).Update("password_hash", string(hashedPassword))

		req := LoginRequest{
			Email:    user.Email,
			Password: password,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := authService.Login(req)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAuthService_ValidateToken(b *testing.B) {
	testutils.RunBenchmarkWithContext(b, func(b *testing.B, ctx *testutils.TestContext) {
		authService := NewAuthService(ctx.DB, ctx.Config.JWT.Secret)
		user := testutils.CreateTestUser(&testing.T{}, ctx.DB)

		token, err := testutils.GenerateTestJWT(user.ID.String(), ctx.Config.JWT.Secret)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := authService.ValidateToken(token)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}