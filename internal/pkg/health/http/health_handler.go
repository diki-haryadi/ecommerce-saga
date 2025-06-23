package http

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/health"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	health *health.Health
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(healthManager *health.Health) *HealthHandler {
	return &HealthHandler{
		health: healthManager,
	}
}

// RegisterRoutes registers the health check routes
func RegisterRoutes(router fiber.Router, handler *HealthHandler) {
	health := router.Group("/health")
	health.Get("/", handler.GetHealth)
	health.Get("/liveness", handler.GetLiveness)
	health.Get("/readiness", handler.GetReadiness)
}

// GetHealth handles the health check endpoint
func (h *HealthHandler) GetHealth(c *fiber.Ctx) error {
	status := h.health.GetStatus(c.Context())
	if status.Status == health.StatusDown {
		return c.Status(fiber.StatusServiceUnavailable).JSON(status)
	}
	return c.JSON(status)
}

// GetLiveness handles the liveness probe endpoint
func (h *HealthHandler) GetLiveness(c *fiber.Ctx) error {
	// Liveness probe just checks if the application is running
	return c.JSON(fiber.Map{
		"status":    "UP",
		"timestamp": time.Now(),
	})
}

// GetReadiness handles the readiness probe endpoint
func (h *HealthHandler) GetReadiness(c *fiber.Ctx) error {
	// Readiness probe checks if the application is ready to handle requests
	if h.health.IsHealthy(c.Context()) {
		return c.JSON(fiber.Map{
			"status":    "UP",
			"timestamp": time.Now(),
		})
	}

	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"status":    "DOWN",
		"timestamp": time.Now(),
	})
}
