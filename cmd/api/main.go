package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/bootstrap"
)

func main() {
	// Load configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Setup database connection
	dbConfig := viper.GetStringMapString("database.postgres")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig["host"],
		dbConfig["port"],
		dbConfig["user"],
		dbConfig["password"],
		dbConfig["dbname"],
		dbConfig["sslmode"],
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	// Setup Fiber app with optimized config
	app := fiber.New(fiber.Config{
		Prefork:              true,  // Enable multiple processes
		DisableKeepalive:     false, // Keep connections alive
		ReadBufferSize:       4096,  // Optimize buffer sizes
		WriteBufferSize:      4096,
		CompressedFileSuffix: ".gz", // Enable compression
		ProxyHeader:          fiber.HeaderXForwardedFor,
		EnablePrintRoutes:    true, // Print routes on startup
	})

	// Middleware
	app.Use(recover.New())  // Recover from panics
	app.Use(logger.New())   // Request logging
	app.Use(compress.New()) // Response compression
	app.Use(cors.New())     // CORS support

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        1000,            // Max 1000 requests
		Expiration: 1 * time.Minute, // Per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP
		},
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Create config map for bootstrap
	config := viper.AllSettings()

	// Initialize bootstrap
	appBootstrap := bootstrap.NewAppBootstrap(db, app, config)

	// API routes group
	api := app.Group("/api/v1")

	// Bootstrap all features
	if err := appBootstrap.Bootstrap(api); err != nil {
		log.Fatalf("Failed to bootstrap application: %s", err)
	}

	// Start server
	serverConfig := viper.GetStringMapString("server")
	addr := fmt.Sprintf("%s:%s", serverConfig["host"], serverConfig["port"])
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
