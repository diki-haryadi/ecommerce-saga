package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/repository"
	saga "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/shared/messaging"
)

func main() {
	// Initialize dependencies
	messageBroker, err := initMessageBroker()
	if err != nil {
		log.Fatalf("Failed to initialize message broker: %v", err)
	}
	defer messageBroker.Close()

	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	sagaRepo := repository.NewPostgresSagaRepository(db)

	// Run migrations
	if err := runMigrations(sagaRepo); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create saga orchestrator
	orchestrator := saga.NewSagaOrchestrator(
		sagaRepo,
		messageBroker,
		5*time.Minute, // Step timeout
		3,             // Max retries
	)

	// Create context that listens for termination signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start processing sagas
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Process any pending sagas
				if err := orchestrator.ProcessPendingSagas(ctx); err != nil {
					log.Printf("Error processing pending sagas: %v", err)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down saga orchestrator...")

	// Give some time for in-flight operations to complete
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownCancel()

	select {
	case <-shutdownCtx.Done():
		log.Println("Shutdown timeout exceeded")
	case <-time.After(100 * time.Millisecond):
		log.Println("Saga orchestrator shutdown complete")
	}
}

func initMessageBroker() (messaging.MessageBroker, error) {
	rabbitmqURI := os.Getenv("RABBITMQ_URI")
	if rabbitmqURI == "" {
		rabbitmqURI = "amqp://guest:guest@localhost:5672/"
	}
	return messaging.NewRabbitMQ(rabbitmqURI)
}

func initDatabase() (*gorm.DB, error) {
	// Get database configuration from environment variables
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPass := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbName := getEnvOrDefault("DB_NAME", "ecommerce")

	// Create DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func runMigrations(repo *repository.PostgresSagaRepository) error {
	db, err := repo.DB().DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Run migrations
	for _, migration := range repo.Migrations() {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
