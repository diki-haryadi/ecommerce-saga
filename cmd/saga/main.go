package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	cartClient "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/delivery/grpc/client"
	orderClient "github.com/diki-haryadi/ecommerce-saga/internal/features/order/delivery/grpc/client"
	orderPostgres "github.com/diki-haryadi/ecommerce-saga/internal/features/order/repository/postgres"
	paymentClient "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/delivery/grpc/client"
	paymentPostgres "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/repository/postgres"
	grpcServer "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/delivery/grpc"
	sagaRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/shared/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := gorm.Open(postgres.Open(cfg.GetPostgresDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize gRPC clients
	orderGrpcClient, err := orderClient.NewOrderClient(fmt.Sprintf("%s:%d", cfg.Services.Order.Host, cfg.Services.Order.Port))
	if err != nil {
		log.Fatalf("Failed to create order client: %v", err)
	}

	paymentGrpcClient, err := paymentClient.NewPaymentClient(fmt.Sprintf("%s:%d", cfg.Services.Payment.Host, cfg.Services.Payment.Port))
	if err != nil {
		log.Fatalf("Failed to create payment client: %v", err)
	}

	cartGrpcClient, err := cartClient.NewCartClient(fmt.Sprintf("%s:%d", cfg.Services.Cart.Host, cfg.Services.Cart.Port))
	if err != nil {
		log.Fatalf("Failed to create cart client: %v", err)
	}

	// Initialize repositories
	sagaRepository := sagaRepo.NewSagaRepository(db)
	orderRepository := orderPostgres.NewOrderRepository(db)
	paymentRepository := paymentPostgres.NewPaymentRepository(db)

	// Initialize usecase with all dependencies
	sagaUsecase := usecase.NewSagaUsecase(
		sagaRepository,
		orderRepository,
		paymentRepository,
		orderGrpcClient,
		paymentGrpcClient,
		cartGrpcClient,
	)

	// Create gRPC server
	grpcServer := grpcServer.NewServer(sagaUsecase)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down Saga gRPC server...")
		grpcServer.Stop()

		// Close gRPC client connections
		orderGrpcClient.Close()
		paymentGrpcClient.Close()
		cartGrpcClient.Close()
	}()

	// Start gRPC server
	fmt.Printf("Starting Saga gRPC server on %s:%d...\n", cfg.GRPC.Host, cfg.GRPC.Port)
	if err := grpcServer.Start(cfg.GRPC.Port); err != nil {
		log.Fatalf("Failed to start Saga gRPC server: %v", err)
	}
}
