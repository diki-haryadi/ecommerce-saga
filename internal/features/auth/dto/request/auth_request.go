package request

// RegisterRequest represents the registration request data
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest represents the login request data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents the refresh token request data
type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// UpdatePasswordRequest represents the password update request data
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}
