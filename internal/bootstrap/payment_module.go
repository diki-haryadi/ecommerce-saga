package bootstrap

import (
	usecase2 "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	orderRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/order/repository/postgres"
	paymentHttp "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/delivery/http"
	paymentRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/eventbus"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/payment/provider"
)

// PaymentModule implements the FeatureModule interface for Payment feature
type PaymentModule struct {
	db             *gorm.DB
	config         *PaymentConfig
	eventBus       *eventbus.EventBus
	paymentUseCase usecase2.Usecase
}

type PaymentConfig struct {
	ProviderType    string // e.g., "stripe", "midtrans"
	APIKey          string
	APISecret       string
	TimeoutDuration time.Duration
	RetryAttempts   int
	WebhookEndpoint string
}

// NewPaymentModule creates a new instance of PaymentModule
func NewPaymentModule(db *gorm.DB, config map[string]interface{}, eventBus *eventbus.EventBus) *PaymentModule {
	return &PaymentModule{
		db: db,
		config: &PaymentConfig{
			ProviderType:    config["payment_provider"].(string),
			APIKey:          config["payment_api_key"].(string),
			APISecret:       config["payment_api_secret"].(string),
			TimeoutDuration: time.Duration(config["payment_timeout_seconds"].(float64)) * time.Second,
			RetryAttempts:   int(config["payment_retry_attempts"].(float64)),
			WebhookEndpoint: config["payment_webhook_endpoint"].(string),
		},
		eventBus: eventBus,
	}
}

// Initialize sets up the payment module
func (m *PaymentModule) Initialize() error {
	// Initialize repositories
	paymentRepo := paymentRepo.NewPaymentRepository(m.db)
	orderRepo := orderRepo.NewOrderRepository(m.db)

	// Initialize payment provider
	paymentProvider, err := provider.NewPaymentProvider(
		m.config.ProviderType,
		provider.Config{
			APIKey:          m.config.APIKey,
			APISecret:       m.config.APISecret,
			TimeoutDuration: m.config.TimeoutDuration,
			RetryAttempts:   m.config.RetryAttempts,
			WebhookEndpoint: m.config.WebhookEndpoint,
		},
	)
	if err != nil {
		return err
	}

	// Initialize payment usecase with dependencies
	m.paymentUseCase = usecase.NewPaymentUsecase(
		paymentRepo,
		orderRepo,
		paymentProvider,
		m.eventBus,
	)

	return nil
}

// RegisterRoutes registers the payment routes
func (m *PaymentModule) RegisterRoutes(router fiber.Router) {
	paymentHttp.RegisterRoutes(router, m.paymentUseCase)
}
