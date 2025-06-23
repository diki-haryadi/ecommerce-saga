package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/dto/response"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/usecase"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
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
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request format",
		})
	}

	if err := h.authUsecase.Register(req.Email, req.Password); err != nil {
		switch err {
		case usecase.ErrUserExists:
			return c.Status(fiber.StatusConflict).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		case usecase.ErrInvalidCredentials:
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
				Error: "Internal server error",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{
		Message: "User registered successfully",
	})
}

// Login handles POST /auth/login request
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request format",
		})
	}

	tokens, err := h.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case usecase.ErrInvalidCredentials:
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
				Error: "Internal server error",
			})
		}
	}

	return c.JSON(response.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
	})
}

// Refresh handles POST /auth/refresh request
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req request.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request format",
		})
	}

	tokens, err := h.authUsecase.RefreshToken(req.Token)
	if err != nil {
		switch err {
		case usecase.ErrInvalidToken:
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
				Error: "Internal server error",
			})
		}
	}

	return c.JSON(response.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
	})
}

// UpdatePassword handles password updates
func (h *AuthHandler) UpdatePassword(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Error: "Invalid user ID",
		})
	}

	var req request.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: "Invalid request format",
		})
	}

	if err := h.authUsecase.UpdatePassword(userID, req.CurrentPassword, req.NewPassword); err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		case usecase.ErrInvalidCredentials:
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
				Error: "Internal server error",
			})
		}
	}

	return c.JSON(response.SuccessResponse{
		Message: "Password updated successfully",
	})
}

// GetJWKS returns the JWKS (JSON Web Key Set)
func (h *AuthHandler) GetJWKS(c *fiber.Ctx) error {
	jwks, err := h.authUsecase.GetJWKS()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: "Failed to get JWKS",
		})
	}

	return c.JSON(fiber.Map{
		"keys": jwks,
	})
}
