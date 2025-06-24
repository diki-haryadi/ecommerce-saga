package bootstrap

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	cartRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/product/service"
)

// CartModule implements the FeatureModule interface for Cart feature
type CartModule struct {
	db          *gorm.DB
	config      *CartConfig
	cartUseCase *usecase.CartUsecase
}

type CartConfig struct {
	CartExpiry time.Duration
}

// NewCartModule creates a new instance of CartModule
func NewCartModule(db *gorm.DB, config map[string]interface{}) *CartModule {
	return &CartModule{
		db: db,
		config: &CartConfig{
			CartExpiry: time.Duration(config["cart_expiry_hours"].(float64)) * time.Hour,
		},
	}
}

// Initialize sets up the cart module
func (m *CartModule) Initialize() error {
	// Initialize repositories and services
	cartRepo := cartRepo.NewCartRepository(m.db)
	productService := service.NewProductService(m.db)

	// Initialize cart usecase with dependencies
	m.cartUseCase = usecase.NewCartUsecase(
		cartRepo,
		productService,
		m.config.CartExpiry,
	)

	return nil
}

// RegisterRoutes registers the cart routes
func (m *CartModule) RegisterRoutes(router fiber.Router) {
	cartGroup := router.Group("/cart")

	cartGroup.Post("/items", m.addItem)
	cartGroup.Delete("/items/:id", m.removeItem)
	cartGroup.Put("/items/:id", m.updateItem)
	cartGroup.Get("/", m.getCart)
	cartGroup.Delete("/", m.clearCart)
}

// Route handlers
func (m *CartModule) addItem(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Add item endpoint coming soon"})
}

func (m *CartModule) removeItem(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Remove item endpoint coming soon"})
}

func (m *CartModule) updateItem(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Update item endpoint coming soon"})
}

func (m *CartModule) getCart(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get cart endpoint coming soon"})
}

func (m *CartModule) clearCart(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Clear cart endpoint coming soon"})
}
