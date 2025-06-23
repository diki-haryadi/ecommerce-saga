package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AppBootstrap represents the application bootstrap facade
type AppBootstrap struct {
	DB     *gorm.DB
	App    *fiber.App
	Config map[string]interface{}
}

// FeatureModule interface for all feature modules
type FeatureModule interface {
	RegisterRoutes(router fiber.Router)
	Initialize() error
}

// NewAppBootstrap creates a new instance of AppBootstrap
func NewAppBootstrap(db *gorm.DB, app *fiber.App, config map[string]interface{}) *AppBootstrap {
	return &AppBootstrap{
		DB:     db,
		App:    app,
		Config: config,
	}
}

// Bootstrap initializes all features and registers their routes
func (b *AppBootstrap) Bootstrap(apiGroup fiber.Router) error {
	// Initialize feature modules using factory
	modules := b.createFeatureModules()

	// Initialize and register each module
	for _, module := range modules {
		if err := module.Initialize(); err != nil {
			return err
		}
		module.RegisterRoutes(apiGroup)
	}

	return nil
}

// createFeatureModules creates all feature modules using factory pattern
func (b *AppBootstrap) createFeatureModules() []FeatureModule {
	return []FeatureModule{
		NewAuthModule(b.DB, b.Config),
		NewCartModule(b.DB, b.Config),
		NewOrderModule(b.DB, b.Config),
		NewPaymentModule(b.DB, b.Config),
		// Add other feature modules here
	}
}
