package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/domain"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/domain/entity"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockUserRepository) GetByRefreshToken(ctx context.Context, token string) (*entity.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

// MockJWKService is a mock implementation of JWKService
type MockJWKService struct {
	mock.Mock
}

func (m *MockJWKService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWKService) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockJWKService) GetJWKS() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func TestAuthUsecase_Register(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func(*MockUserRepository)
		expectedError error
	}{
		{
			name:     "successful registration",
			email:    "test@example.com",
			password: "Password123!",
			mockSetup: func(repo *MockUserRepository) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "user already exists",
			email:    "existing@example.com",
			password: "Password123!",
			mockSetup: func(repo *MockUserRepository) {
				repo.On("GetByEmail", mock.Anything, "existing@example.com").Return(&entity.User{}, nil)
			},
			expectedError: ErrUserExists,
		},
		{
			name:          "invalid email",
			email:         "invalid-email",
			password:      "Password123!",
			mockSetup:     func(repo *MockUserRepository) {},
			expectedError: domain.ErrInvalidEmail,
		},
		{
			name:          "weak password",
			email:         "test@example.com",
			password:      "weak",
			mockSetup:     func(repo *MockUserRepository) {},
			expectedError: domain.ErrPasswordTooWeak,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			repo := new(MockUserRepository)
			tt.mockSetup(repo)
			usecase := NewAuthUsecase(repo, nil)

			// Execute
			err := usecase.Register(tt.email, tt.password)

			// Assert
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestAuthUsecase_Login(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func(*MockUserRepository, *MockJWKService)
		expectedError error
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "Password123!",
			mockSetup: func(repo *MockUserRepository, jwk *MockJWKService) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(&entity.User{
					ID:           userID,
					Email:        "test@example.com",
					PasswordHash: "$2a$10$abcdefghijklmnopqrstuvwxyz",
				}, nil)
				repo.On("UpdateRefreshToken", mock.Anything, userID, mock.AnythingOfType("string")).Return(nil)
				jwk.On("GenerateAccessToken", userID).Return("access-token", nil)
				jwk.On("GenerateRefreshToken").Return("refresh-token", nil)
			},
			expectedError: nil,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: "Password123!",
			mockSetup: func(repo *MockUserRepository, jwk *MockJWKService) {
				repo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, nil)
			},
			expectedError: ErrInvalidCredentials,
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "WrongPassword123!",
			mockSetup: func(repo *MockUserRepository, jwk *MockJWKService) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(&entity.User{
					ID:           userID,
					Email:        "test@example.com",
					PasswordHash: "$2a$10$abcdefghijklmnopqrstuvwxyz",
				}, nil)
			},
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			repo := new(MockUserRepository)
			jwk := new(MockJWKService)
			tt.mockSetup(repo, jwk)
			usecase := NewAuthUsecase(repo, jwk)

			// Execute
			tokens, err := usecase.Login(tt.email, tt.password)

			// Assert
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, tokens)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokens)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
			repo.AssertExpectations(t)
			jwk.AssertExpectations(t)
		})
	}
}

func TestAuthUsecase_RefreshToken(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name          string
		refreshToken  string
		mockSetup     func(*MockUserRepository, *MockJWKService)
		expectedError error
	}{
		{
			name:         "successful token refresh",
			refreshToken: "valid-refresh-token",
			mockSetup: func(repo *MockUserRepository, jwk *MockJWKService) {
				repo.On("GetByRefreshToken", mock.Anything, "valid-refresh-token").Return(&entity.User{
					ID: userID,
				}, nil)
				repo.On("UpdateRefreshToken", mock.Anything, userID, mock.AnythingOfType("string")).Return(nil)
				jwk.On("GenerateAccessToken", userID).Return("new-access-token", nil)
				jwk.On("GenerateRefreshToken").Return("new-refresh-token", nil)
			},
			expectedError: nil,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid-refresh-token",
			mockSetup: func(repo *MockUserRepository, jwk *MockJWKService) {
				repo.On("GetByRefreshToken", mock.Anything, "invalid-refresh-token").Return(nil, ErrInvalidToken)
			},
			expectedError: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			repo := new(MockUserRepository)
			jwk := new(MockJWKService)
			tt.mockSetup(repo, jwk)
			usecase := NewAuthUsecase(repo, jwk)

			// Execute
			tokens, err := usecase.RefreshToken(tt.refreshToken)

			// Assert
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, tokens)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokens)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
			repo.AssertExpectations(t)
			jwk.AssertExpectations(t)
		})
	}
}
