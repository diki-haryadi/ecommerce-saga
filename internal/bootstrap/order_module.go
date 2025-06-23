package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// OrderModule implements the FeatureModule interface for Order feature
type OrderModule struct {
	db     *gorm.DB
	config map[string]interface{}
}

// NewOrderModule creates a new instance of OrderModule
func NewOrderModule(db *gorm.DB, config map[string]interface{}) *OrderModule {
	return &OrderModule{
		db:     db,
		config: config,
	}
}

// Initialize sets up the order module
func (m *OrderModule) Initialize() error {
	// TODO: Initialize order repositories
	// TODO: Initialize saga orchestrator
	// TODO: Initialize order usecase with saga support
	return nil
}

// RegisterRoutes registers the order routes
func (m *OrderModule) RegisterRoutes(router fiber.Router) {
	// Order routes will be implemented here
	orderGroup := router.Group("/orders")

	orderGroup.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Order feature coming soon"})
	})
}
