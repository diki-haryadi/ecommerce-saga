package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/service"
)

// JWTMiddleware creates a middleware for JWT authentication
func JWTMiddleware(jwkService *service.JWKService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		// Validate the token
		claims, err := jwkService.ValidateToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Store user ID in context
		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}

// Protected creates a middleware that requires authentication
func Protected(jwkService *service.JWKService) fiber.Handler {
	return JWTMiddleware(jwkService)
}
