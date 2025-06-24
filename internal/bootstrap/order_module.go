package bootstrap

import (
	usecase2 "github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	cartRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/delivery/http"
	orderRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/order/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/eventbus"
)

// OrderModule implements the FeatureModule interface for Order feature
type OrderModule struct {
	db           *gorm.DB
	config       *Config
	eventBus     *eventbus.EventBus
	orderUseCase usecase2.Usecase
}

type Config struct {
	MaxOrderItems int
	MinOrderValue float64
}

// NewOrderModule creates a new instance of OrderModule
func NewOrderModule(db *gorm.DB, config *Config, eventBus *eventbus.EventBus) *OrderModule {
	return &OrderModule{
		db:       db,
		config:   config,
		eventBus: eventBus,
	}
}

// Initialize sets up the order module
func (m *OrderModule) Initialize() error {
	// Initialize repositories
	orderRepo := orderRepo.NewOrderRepository(m.db)
	cartRepo := cartRepo.NewCartRepository(m.db)

	// Initialize order usecase with dependencies
	m.orderUseCase = usecase.NewOrderUsecase(orderRepo, cartRepo)

	return nil
}

// RegisterRoutes registers the order routes
func (m *OrderModule) RegisterRoutes(router fiber.Router) {
	handler := http.NewOrderHandler(m.orderUseCase)
	http.RegisterRoutes(router, handler, nil) // nil for authMiddleware since it should be handled at a higher level
}
