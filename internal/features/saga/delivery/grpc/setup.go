package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga"
	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/delivery/grpc/proto"
)

type Server struct {
	server      *grpc.Server
	sagaServer  *SagaServer
	sagaUsecase saga.Usecase
}

func NewServer(sagaUsecase saga.Usecase) *Server {
	// Create gRPC server
	server := grpc.NewServer()

	// Create saga server
	sagaServer := NewSagaServer(sagaUsecase)

	// Register services
	pb.RegisterSagaServiceServer(server, sagaServer)

	// Register reflection service on gRPC server
	reflection.Register(server)

	return &Server{
		server:      server,
		sagaServer:  sagaServer,
		sagaUsecase: sagaUsecase,
	}
}

func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return s.server.Serve(listener)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
