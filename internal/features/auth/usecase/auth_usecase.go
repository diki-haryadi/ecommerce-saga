package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/domain"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

// AuthUsecase defines the interface for authentication use cases
type AuthUsecase interface {
	Register(email, password string) error
	Login(email, password string) (*TokenPair, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
	UpdatePassword(userID uuid.UUID, currentPassword, newPassword string) error
	GetJWKS() ([]byte, error)
}

// JWKService defines the interface for JWT operations
type JWKService interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken() (string, error)
	GetJWKS() ([]byte, error)
}

type authUsecase struct {
	userRepo   repository.UserRepository
	jwkService JWKService
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthUsecase(userRepo repository.UserRepository, jwkService JWKService) AuthUsecase {
	return &authUsecase{
		userRepo:   userRepo,
		jwkService: jwkService,
	}
}

// Register creates a new user account
func (u *authUsecase) Register(email, password string) error {
	ctx := context.Background()

	// Validate email and password
	if err := domain.ValidateEmail(email); err != nil {
		return err
	}
	if err := domain.ValidatePassword(password); err != nil {
		return err
	}

	// Check if user exists
	existingUser, _ := u.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return ErrUserExists
	}

	// Create user
	user, err := entity.NewUser(email, password)
	if err != nil {
		return err
	}

	return u.userRepo.Create(ctx, user)
}

// Login authenticates a user and returns access and refresh tokens
func (u *authUsecase) Login(email, password string) (*TokenPair, error) {
	ctx := context.Background()

	// Find user
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !user.ValidatePassword(password) {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := u.jwkService.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.jwkService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Store refresh token
	if err := u.userRepo.UpdateRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken generates new access token using refresh token
func (u *authUsecase) RefreshToken(refreshToken string) (*TokenPair, error) {
	ctx := context.Background()

	// Find user by refresh token
	user, err := u.userRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Generate new tokens
	accessToken, err := u.jwkService.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := u.jwkService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Update refresh token
	if err := u.userRepo.UpdateRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// UpdatePassword updates a user's password
func (u *authUsecase) UpdatePassword(userID uuid.UUID, currentPassword, newPassword string) error {
	ctx := context.Background()

	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify current password
	if !user.ValidatePassword(currentPassword) {
		return ErrInvalidCredentials
	}

	// Validate new password
	if err := domain.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Update password
	if err := user.UpdatePassword(newPassword); err != nil {
		return err
	}

	return u.userRepo.Update(ctx, user)
}

// GetJWKS returns the public JWK set
func (u *authUsecase) GetJWKS() ([]byte, error) {
	return u.jwkService.GetJWKS()
}
