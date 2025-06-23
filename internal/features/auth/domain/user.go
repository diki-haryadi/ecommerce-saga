package domain

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrPasswordTooWeak = errors.New("password must be at least 8 characters long and contain at least one number, one uppercase letter, and one special character")
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Password     string    `json:"-" gorm:"not null"` // "-" means this field won't be included in JSON
	RefreshToken string    `json:"-" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword checks if the password meets security requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooWeak
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	if !hasNumber || !hasUpper || !hasSpecial {
		return ErrPasswordTooWeak
	}

	return nil
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with a plain password
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
