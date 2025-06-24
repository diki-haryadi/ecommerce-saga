package auth

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	authhttp "github.com/diki-haryadi/ecommerce-saga/internal/features/auth/delivery/http"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/service"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/testutil"
)

type AuthIntegrationTestSuite struct {
	suite.Suite
	app      *fiber.App
	db       *testutil.TestDB
	usecase  usecase.AuthUsecase
	handler  *authhttp.AuthHandler
	fixtures struct {
		users []entity.User
	}
}

func TestAuthIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}

func (s *AuthIntegrationTestSuite) SetupSuite() {
	// Create test database
	s.db = testutil.NewTestPostgres(s.T())

	// Run migrations
	err := s.db.AutoMigrate(&entity.User{})
	require.NoError(s.T(), err)

	// Load fixtures
	data := testutil.LoadFixture(s.T(), "user.json")
	err = json.Unmarshal(data, &s.fixtures)
	require.NoError(s.T(), err)

	// Create postgres
	userRepo := postgres.NewUserRepository(s.db.DB)

	// Create JWK service
	jwkService, err := service.NewJWKService(24 * time.Hour)
	require.NoError(s.T(), err)

	// Create usecase
	s.usecase = usecase.NewAuthUsecase(userRepo, jwkService)

	// Create handler
	s.handler = authhttp.NewAuthHandler(s.usecase)

	// Create Fiber app
	s.app = fiber.New()
	authhttp.RegisterRoutes(s.app, s.handler, []byte("test-secret"))
}

func (s *AuthIntegrationTestSuite) TearDownSuite() {
	s.db.Cleanup(s.T())
}

func (s *AuthIntegrationTestSuite) TestRegister() {
	tests := []struct {
		name        string
		payload     map[string]interface{}
		expectCode  int
		expectError bool
	}{
		{
			name: "successful registration",
			payload: map[string]interface{}{
				"email":    "newuser@example.com",
				"password": "Password123!",
			},
			expectCode:  fiber.StatusCreated,
			expectError: false,
		},
		{
			name: "existing user",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "Password123!",
			},
			expectCode:  fiber.StatusBadRequest,
			expectError: true,
		},
		{
			name: "invalid email",
			payload: map[string]interface{}{
				"email":    "invalid-email",
				"password": "Password123!",
			},
			expectCode:  fiber.StatusBadRequest,
			expectError: true,
		},
		{
			name: "weak password",
			payload: map[string]interface{}{
				"email":    "newuser@example.com",
				"password": "weak",
			},
			expectCode:  fiber.StatusBadRequest,
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Create request
			body, err := json.Marshal(tt.payload)
			require.NoError(s.T(), err)

			req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := s.app.Test(req)
			require.NoError(s.T(), err)

			// Assert response
			require.Equal(s.T(), tt.expectCode, resp.StatusCode)

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(s.T(), err)

			if tt.expectError {
				require.Contains(s.T(), result, "error")
			} else {
				require.Contains(s.T(), result, "message")
				require.Equal(s.T(), "User registered successfully", result["message"])
			}
		})
	}
}

func (s *AuthIntegrationTestSuite) TestLogin() {
	tests := []struct {
		name        string
		payload     map[string]interface{}
		expectCode  int
		expectError bool
	}{
		{
			name: "successful login",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "Password123!",
			},
			expectCode:  fiber.StatusOK,
			expectError: false,
		},
		{
			name: "invalid credentials",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "WrongPassword123!",
			},
			expectCode:  fiber.StatusUnauthorized,
			expectError: true,
		},
		{
			name: "user not found",
			payload: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "Password123!",
			},
			expectCode:  fiber.StatusUnauthorized,
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Create request
			body, err := json.Marshal(tt.payload)
			require.NoError(s.T(), err)

			req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := s.app.Test(req)
			require.NoError(s.T(), err)

			// Assert response
			require.Equal(s.T(), tt.expectCode, resp.StatusCode)

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(s.T(), err)

			if tt.expectError {
				require.Contains(s.T(), result, "error")
			} else {
				require.Contains(s.T(), result, "access_token")
				require.Contains(s.T(), result, "refresh_token")
			}
		})
	}
}

func (s *AuthIntegrationTestSuite) TestRefreshToken() {
	// First, get a valid refresh token
	loginReq := map[string]interface{}{
		"email":    "test@example.com",
		"password": "Password123!",
	}

	body, err := json.Marshal(loginReq)
	require.NoError(s.T(), err)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	require.NoError(s.T(), err)
	require.Equal(s.T(), fiber.StatusOK, resp.StatusCode)

	var loginResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&loginResult)
	require.NoError(s.T(), err)

	refreshToken := loginResult["refresh_token"].(string)

	tests := []struct {
		name        string
		token       string
		expectCode  int
		expectError bool
	}{
		{
			name:        "successful refresh",
			token:       refreshToken,
			expectCode:  fiber.StatusOK,
			expectError: false,
		},
		{
			name:        "invalid token",
			token:       "invalid-token",
			expectCode:  fiber.StatusUnauthorized,
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Create request
			body, err := json.Marshal(map[string]interface{}{
				"refresh_token": tt.token,
			})
			require.NoError(s.T(), err)

			req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := s.app.Test(req)
			require.NoError(s.T(), err)

			// Assert response
			require.Equal(s.T(), tt.expectCode, resp.StatusCode)

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(s.T(), err)

			if tt.expectError {
				require.Contains(s.T(), result, "error")
			} else {
				require.Contains(s.T(), result, "access_token")
				require.Contains(s.T(), result, "refresh_token")
			}
		})
	}
}
