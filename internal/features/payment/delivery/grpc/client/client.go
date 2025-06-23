package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/delivery/grpc/proto"
)

// PaymentClient represents the gRPC client for payment service
type PaymentClient struct {
	client pb.PaymentServiceClient
	conn   *grpc.ClientConn
}

// NewPaymentClient creates a new payment gRPC client
func NewPaymentClient(address string) (*PaymentClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := pb.NewPaymentServiceClient(conn)
	return &PaymentClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *PaymentClient) Close() error {
	return c.conn.Close()
}

// CreatePayment creates a new payment
func (c *PaymentClient) CreatePayment(ctx context.Context, orderID string, amount float64, currency, paymentMethod string) (*pb.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.CreatePaymentRequest{
		OrderId:       orderID,
		Amount:        amount,
		Currency:      currency,
		PaymentMethod: paymentMethod,
	}

	resp, err := c.client.CreatePayment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return resp.Payment, nil
}

// GetPayment retrieves a payment by ID
func (c *PaymentClient) GetPayment(ctx context.Context, paymentID string) (*pb.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.GetPaymentRequest{
		PaymentId: paymentID,
	}

	resp, err := c.client.GetPayment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return resp.Payment, nil
}

// ListPayments retrieves a list of payments
func (c *PaymentClient) ListPayments(ctx context.Context, page, limit int32, status string) (*pb.ListPaymentsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.ListPaymentsRequest{
		Page:   page,
		Limit:  limit,
		Status: status,
	}

	return c.client.ListPayments(ctx, req)
}

// ProcessPayment processes a payment
func (c *PaymentClient) ProcessPayment(ctx context.Context, paymentID string, details *pb.PaymentDetails) (*pb.ProcessPaymentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.ProcessPaymentRequest{
		PaymentId:      paymentID,
		PaymentDetails: details,
	}

	return c.client.ProcessPayment(ctx, req)
}

// RefundPayment processes a refund
func (c *PaymentClient) RefundPayment(ctx context.Context, paymentID string, amount float64, reason string) (*pb.RefundPaymentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.RefundPaymentRequest{
		PaymentId: paymentID,
		Amount:    amount,
		Reason:    reason,
	}

	return c.client.RefundPayment(ctx, req)
}

// WithToken adds an authorization token to the context
func WithToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
}
