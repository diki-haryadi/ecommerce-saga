package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/errors"
	httpresponse "github.com/diki-haryadi/ecommerce-saga/internal/pkg/http/response"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authUsecase  usecase.AuthUsecase
	errorHandler errors.ErrorHandler
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase:  authUsecase,
		errorHandler: errors.NewErrorHandler(),
	}
}

// RegisterRoutes registers all auth-related routes
func RegisterRoutes(router fiber.Router, handler *AuthHandler, jwtSecret []byte) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.Refresh)
	auth.Get("/.well-known/jwks.json", handler.GetJWKS)

	// Protected routes
	auth.Use(AuthMiddleware(jwtSecret))
	auth.Put("/password", handler.UpdatePassword)
}

// Register handles POST /auth/register request
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	if err := h.authUsecase.Register(req.Email, req.Password); err != nil {
		switch err {
		case usecase.ErrUserExists:
			return h.errorHandler.Handle(c, errors.NewConflictError(err.Error()))
		case usecase.ErrInvalidCredentials:
			return h.errorHandler.Handle(c, errors.NewValidationError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.Created(c, "User registered successfully", nil)
}

// Login handles POST /auth/login request
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	tokens, err := h.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case usecase.ErrInvalidCredentials:
			return h.errorHandler.Handle(c, errors.NewAuthenticationError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Login successful", fiber.Map{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600, // 1 hour
	})
}

// Refresh handles POST /auth/refresh request
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req request.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	tokens, err := h.authUsecase.RefreshToken(req.Token)
	if err != nil {
		switch err {
		case usecase.ErrInvalidToken:
			return h.errorHandler.Handle(c, errors.NewAuthenticationError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Token refreshed successfully", fiber.Map{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600, // 1 hour
	})
}

// UpdatePassword handles password updates
func (h *AuthHandler) UpdatePassword(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewAuthenticationError("Invalid user ID"))
	}

	var req request.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, errors.NewValidationError("Invalid request format"))
	}

	if err := h.authUsecase.UpdatePassword(userID, req.CurrentPassword, req.NewPassword); err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return h.errorHandler.Handle(c, errors.NewNotFoundError(err.Error()))
		case usecase.ErrInvalidCredentials:
			return h.errorHandler.Handle(c, errors.NewAuthenticationError(err.Error()))
		default:
			return h.errorHandler.Handle(c, errors.NewInternalError(err))
		}
	}

	return httpresponse.OK(c, "Password updated successfully", nil)
}

// GetJWKS returns the JWKS (JSON Web Key Set)
func (h *AuthHandler) GetJWKS(c *fiber.Ctx) error {
	jwks, err := h.authUsecase.GetJWKS()
	if err != nil {
		return h.errorHandler.Handle(c, errors.NewInternalError(err))
	}

	return httpresponse.OK(c, "JWKS retrieved successfully", fiber.Map{
		"keys": jwks,
	})
}
