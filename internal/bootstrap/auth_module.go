package bootstrap

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/delivery/http"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/service"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/usecase"
)

// AuthModule implements the FeatureModule interface for Auth feature
type AuthModule struct {
	db         *gorm.DB
	config     map[string]interface{}
	usecase    usecase.AuthUsecase
	jwkService *service.JWKService
}

// NewAuthModule creates a new instance of AuthModule
func NewAuthModule(db *gorm.DB, config map[string]interface{}) *AuthModule {
	return &AuthModule{
		db:     db,
		config: config,
	}
}

// Initialize sets up the auth module
func (m *AuthModule) Initialize() error {
	// Initialize JWK service with key rotation
	rotationPeriod := 24 * time.Hour // Rotate keys every 24 hours
	jwkService, err := service.NewJWKService(rotationPeriod)
	if err != nil {
		return err
	}
	m.jwkService = jwkService

	// Initialize repository
	userRepo := postgres.NewUserRepository(m.db)

	// Initialize usecase
	m.usecase = usecase.NewAuthUsecase(userRepo, jwkService)

	return nil
}

// RegisterRoutes registers the auth routes
func (m *AuthModule) RegisterRoutes(router fiber.Router) {
	handler := http.NewAuthHandler(m.usecase)
	http.RegisterRoutes(router, handler, []byte(m.config["jwt_secret"].(string)))
}

// GetProtectedMiddleware returns the JWT authentication middleware
func (m *AuthModule) GetProtectedMiddleware() fiber.Handler {
	return http.AuthMiddleware([]byte(m.config["jwt_secret"].(string)))
}
