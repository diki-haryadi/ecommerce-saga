package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/dto/response"
	"github.com/diki-haryadi/ecommerce-saga/internal/shared/utils"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Error: "No authorization header",
			})
		}

		// Extract token from Bearer header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Error: "Invalid authorization header format",
			})
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenParts[1], secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		}

		// Set user ID in context
		c.Locals("user_id", claims.UserID.String())
		return c.Next()
	}
}
