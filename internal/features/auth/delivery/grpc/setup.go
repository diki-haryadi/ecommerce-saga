package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/auth/delivery/grpc/proto"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/usecase"
)

// Server represents the gRPC server
type Server struct {
	server      *grpc.Server
	authServer  *AuthServer
	interceptor *AuthInterceptor
	authUsecase usecase.AuthUsecase
}

// NewServer creates a new gRPC server
func NewServer(authUsecase usecase.AuthUsecase) *Server {
	// Create auth interceptor
	interceptor := NewAuthInterceptor(authUsecase)

	// Create gRPC server with interceptors
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	// Create auth server
	authServer := NewAuthServer(authUsecase)

	// Register services
	pb.RegisterAuthServiceServer(server, authServer)

	// Register reflection service on gRPC server
	reflection.Register(server)

	return &Server{
		server:      server,
		authServer:  authServer,
		interceptor: interceptor,
		authUsecase: authUsecase,
	}
}

// Start starts the gRPC server
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	fmt.Printf("gRPC server is listening on port %d\n", port)
	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Stop stops the gRPC server
func (s *Server) Stop() {
	s.server.GracefulStop()
}
