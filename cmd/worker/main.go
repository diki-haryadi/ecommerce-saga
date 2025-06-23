package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/repository"
	saga "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/shared/messaging"
)

type StepMessage struct {
	SagaID  uuid.UUID       `json:"saga_id"`
	OrderID uuid.UUID       `json:"order_id"`
	Step    entity.SagaStep `json:"step"`
}

type Worker struct {
	messageBroker    messaging.MessageBroker
	sagaOrchestrator *saga.SagaOrchestrator
	db               *gorm.DB
}

func main() {
	// Initialize dependencies
	worker, err := initWorker()
	if err != nil {
		log.Fatalf("Failed to initialize worker: %v", err)
	}
	defer worker.messageBroker.Close()

	// Create context that listens for termination signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Subscribe to saga step topics
	topics := []string{
		"saga.CREATE_ORDER",
		"saga.PROCESS_PAYMENT",
		"saga.UPDATE_INVENTORY",
		"saga.compensation.CREATE_ORDER",
		"saga.compensation.PROCESS_PAYMENT",
		"saga.compensation.UPDATE_INVENTORY",
	}

	// Start workers for each topic
	for _, topic := range topics {
		if err := worker.messageBroker.Subscribe(topic, worker.createStepHandler(topic)); err != nil {
			log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
		}
	}

	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down worker...")

	// Give some time for in-flight messages to complete
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownCancel()

	select {
	case <-shutdownCtx.Done():
		log.Println("Shutdown timeout exceeded")
	case <-time.After(100 * time.Millisecond):
		log.Println("Worker shutdown complete")
	}
}

func (w *Worker) createStepHandler(topic string) messaging.MessageHandler {
	return func(ctx context.Context, message []byte) error {
		// Parse message
		var stepMsg StepMessage
		if err := json.Unmarshal(message, &stepMsg); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		// Log step processing
		log.Printf("Processing step %s for saga %s", stepMsg.Step.Name, stepMsg.SagaID)

		// Process step based on topic
		switch topic {
		case "saga.CREATE_ORDER":
			return w.handleCreateOrder(ctx, stepMsg)
		case "saga.PROCESS_PAYMENT":
			return w.handleProcessPayment(ctx, stepMsg)
		case "saga.UPDATE_INVENTORY":
			return w.handleUpdateInventory(ctx, stepMsg)
		case "saga.compensation.CREATE_ORDER":
			return w.handleCreateOrderCompensation(ctx, stepMsg)
		case "saga.compensation.PROCESS_PAYMENT":
			return w.handleProcessPaymentCompensation(ctx, stepMsg)
		case "saga.compensation.UPDATE_INVENTORY":
			return w.handleUpdateInventoryCompensation(ctx, stepMsg)
		default:
			return fmt.Errorf("unknown topic: %s", topic)
		}
	}
}

func (w *Worker) handleCreateOrder(ctx context.Context, msg StepMessage) error {
	// Start a transaction
	tx := w.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Create order logic
	order := struct {
		ID      uuid.UUID `gorm:"type:uuid;primary_key"`
		OrderID uuid.UUID `gorm:"type:uuid"`
		Status  string    `gorm:"type:varchar(50)"`
	}{
		ID:      uuid.New(),
		OrderID: msg.OrderID,
		Status:  "CREATED",
	}

	if err := tx.Create(&order).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Update saga step status
	if err := w.sagaOrchestrator.ProcessStepResult(ctx, msg.SagaID, msg.Step, entity.StepStatusSuccess, nil); err != nil {
		return fmt.Errorf("failed to update saga step: %w", err)
	}

	return tx.Commit().Error
}

func (w *Worker) handleProcessPayment(ctx context.Context, msg StepMessage) error {
	// Start a transaction
	tx := w.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Process payment logic
	payment := struct {
		ID      uuid.UUID `gorm:"type:uuid;primary_key"`
		OrderID uuid.UUID `gorm:"type:uuid"`
		Status  string    `gorm:"type:varchar(50)"`
		Amount  float64   `gorm:"type:decimal(10,2)"`
	}{
		ID:      uuid.New(),
		OrderID: msg.OrderID,
		Status:  "COMPLETED",
		Amount:  100.00, // Example amount
	}

	if err := tx.Create(&payment).Error; err != nil {
		return fmt.Errorf("failed to process payment: %w", err)
	}

	// Update saga step status
	if err := w.sagaOrchestrator.ProcessStepResult(ctx, msg.SagaID, msg.Step, entity.StepStatusSuccess, nil); err != nil {
		return fmt.Errorf("failed to update saga step: %w", err)
	}

	return tx.Commit().Error
}

func (w *Worker) handleUpdateInventory(ctx context.Context, msg StepMessage) error {
	// Start a transaction
	tx := w.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Update inventory logic
	inventory := struct {
		ID       uuid.UUID `gorm:"type:uuid;primary_key"`
		OrderID  uuid.UUID `gorm:"type:uuid"`
		Status   string    `gorm:"type:varchar(50)"`
		Quantity int       `gorm:"type:integer"`
	}{
		ID:       uuid.New(),
		OrderID:  msg.OrderID,
		Status:   "RESERVED",
		Quantity: 1, // Example quantity
	}

	if err := tx.Create(&inventory).Error; err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	// Update saga step status
	if err := w.sagaOrchestrator.ProcessStepResult(ctx, msg.SagaID, msg.Step, entity.StepStatusSuccess, nil); err != nil {
		return fmt.Errorf("failed to update saga step: %w", err)
	}

	return tx.Commit().Error
}

func (w *Worker) handleCreateOrderCompensation(ctx context.Context, msg StepMessage) error {
	// Start a transaction
	tx := w.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Cancel order logic
	if err := tx.Model(struct{ OrderID uuid.UUID }{}).
		Where("order_id = ?", msg.OrderID).
		Update("status", "CANCELLED").Error; err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// Update saga step status
	if err := w.sagaOrchestrator.ProcessStepResult(ctx, msg.SagaID, msg.Step, entity.StepStatusCompensated, nil); err != nil {
		return fmt.Errorf("failed to update saga step: %w", err)
	}

	return tx.Commit().Error
}

func (w *Worker) handleProcessPaymentCompensation(ctx context.Context, msg StepMessage) error {
	// Start a transaction
	tx := w.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Refund payment logic
	if err := tx.Model(struct{ OrderID uuid.UUID }{}).
		Where("order_id = ?", msg.OrderID).
		Update("status", "REFUNDED").Error; err != nil {
		return fmt.Errorf("failed to refund payment: %w", err)
	}

	// Update saga step status
	if err := w.sagaOrchestrator.ProcessStepResult(ctx, msg.SagaID, msg.Step, entity.StepStatusCompensated, nil); err != nil {
		return fmt.Errorf("failed to update saga step: %w", err)
	}

	return tx.Commit().Error
}

func (w *Worker) handleUpdateInventoryCompensation(ctx context.Context, msg StepMessage) error {
	// Start a transaction
	tx := w.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Release inventory logic
	if err := tx.Model(struct{ OrderID uuid.UUID }{}).
		Where("order_id = ?", msg.OrderID).
		Update("status", "RELEASED").Error; err != nil {
		return fmt.Errorf("failed to release inventory: %w", err)
	}

	// Update saga step status
	if err := w.sagaOrchestrator.ProcessStepResult(ctx, msg.SagaID, msg.Step, entity.StepStatusCompensated, nil); err != nil {
		return fmt.Errorf("failed to update saga step: %w", err)
	}

	return tx.Commit().Error
}

func initWorker() (*Worker, error) {
	// Initialize message broker
	rabbitmqURI := os.Getenv("RABBITMQ_URI")
	if rabbitmqURI == "" {
		rabbitmqURI = "amqp://guest:guest@localhost:5672/"
	}
	messageBroker, err := messaging.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize message broker: %w", err)
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize saga orchestrator
	sagaRepo := repository.NewPostgresSagaRepository(db)
	sagaOrchestrator := saga.NewSagaOrchestrator(
		sagaRepo,
		messageBroker,
		5*time.Minute, // Step timeout
		3,             // Max retries
	)

	return &Worker{
		messageBroker:    messageBroker,
		sagaOrchestrator: sagaOrchestrator,
		db:               db,
	}, nil
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

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
