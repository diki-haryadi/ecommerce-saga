package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	grpcServer "github.com/diki-haryadi/ecommerce-saga/internal/features/auth/delivery/grpc"
	authRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/auth/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/jwt"
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

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg.Auth.JWT.Secret, cfg.Auth.JWT.TokenExpiry, cfg.Auth.JWT.RefreshExpiry)

	// Initialize repositories
	userRepo := authRepo.NewUserRepository(db)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService)

	// Create gRPC server
	server := grpcServer.NewServer(authUsecase)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down gRPC server...")
		server.Stop()
	}()

	// Start gRPC server
	fmt.Printf("Starting gRPC server on %s:%d...\n", cfg.GRPC.Host, cfg.GRPC.Port)
	if err := server.Start(cfg.GRPC.Port); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
