package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/domain/entity"
)

// UserRepository implements the repository.UserRepository interface
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create saves a new user to the database
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email address
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete removes a user from the database
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

// UpdateRefreshToken updates the refresh token for a user
func (r *UserRepository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).
		Update("refresh_token", token).Error
}

// GetByRefreshToken retrieves a user by their refresh token
func (r *UserRepository) GetByRefreshToken(ctx context.Context, token string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("refresh_token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid refresh token")
		}
		return nil, err
	}
	return &user, nil
}
