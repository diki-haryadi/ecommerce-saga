package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/delivery/grpc/proto"
)

// SagaClient represents the gRPC client for saga service
type SagaClient struct {
	client pb.SagaServiceClient
	conn   *grpc.ClientConn
}

// NewSagaClient creates a new saga gRPC client
func NewSagaClient(address string) (*SagaClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := pb.NewSagaServiceClient(conn)
	return &SagaClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *SagaClient) Close() error {
	return c.conn.Close()
}

// StartOrderSaga starts a new order saga transaction
func (c *SagaClient) StartOrderSaga(ctx context.Context, orderID, userID string, amount float64, paymentMethod string, metadata map[string]string) (*pb.SagaTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.StartOrderSagaRequest{
		OrderId:       orderID,
		UserId:        userID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Metadata:      metadata,
	}

	resp, err := c.client.StartOrderSaga(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to start saga: %w", err)
	}

	return resp.Transaction, nil
}

// GetSagaStatus retrieves the status of a saga transaction
func (c *SagaClient) GetSagaStatus(ctx context.Context, sagaID string) (*pb.SagaTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.GetSagaStatusRequest{
		SagaId: sagaID,
	}

	resp, err := c.client.GetSagaStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get saga status: %w", err)
	}

	return resp.Transaction, nil
}

// CompensateTransaction initiates compensation for a saga transaction
func (c *SagaClient) CompensateTransaction(ctx context.Context, sagaID, stepID, reason string) (*pb.SagaTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.CompensateTransactionRequest{
		SagaId: sagaID,
		StepId: stepID,
		Reason: reason,
	}

	resp, err := c.client.CompensateTransaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to compensate transaction: %w", err)
	}

	return resp.Transaction, nil
}

// ListSagaTransactions retrieves a list of saga transactions
func (c *SagaClient) ListSagaTransactions(ctx context.Context, page, limit int32, status, transactionType string) (*pb.ListSagaTransactionsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.ListSagaTransactionsRequest{
		Page:   page,
		Limit:  limit,
		Status: status,
		Type:   transactionType,
	}

	return c.client.ListSagaTransactions(ctx, req)
}

// WithToken adds an authorization token to the context
func WithToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
}
