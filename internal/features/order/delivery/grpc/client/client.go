package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/order/delivery/grpc/proto"
)

// OrderClient represents the gRPC client for order service
type OrderClient struct {
	client pb.OrderServiceClient
	conn   *grpc.ClientConn
}

// NewOrderClient creates a new order gRPC client
func NewOrderClient(address string) (*OrderClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := pb.NewOrderServiceClient(conn)
	return &OrderClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *OrderClient) Close() error {
	return c.conn.Close()
}

// CreateOrder creates a new order
func (c *OrderClient) CreateOrder(ctx context.Context, shippingAddress, paymentMethod string) (*pb.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.CreateOrderRequest{
		ShippingAddress: shippingAddress,
		PaymentMethod:   paymentMethod,
	}

	resp, err := c.client.CreateOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return resp.Order, nil
}

// GetOrder retrieves an order by ID
func (c *OrderClient) GetOrder(ctx context.Context, orderID string) (*pb.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.GetOrderRequest{
		OrderId: orderID,
	}

	resp, err := c.client.GetOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return resp.Order, nil
}

// ListOrders retrieves a list of orders
func (c *OrderClient) ListOrders(ctx context.Context, page, limit int32, status string) (*pb.ListOrdersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.ListOrdersRequest{
		Page:   page,
		Limit:  limit,
		Status: status,
	}

	return c.client.ListOrders(ctx, req)
}

// CancelOrder cancels an order
func (c *OrderClient) CancelOrder(ctx context.Context, orderID, reason string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.CancelOrderRequest{
		OrderId: orderID,
		Reason:  reason,
	}

	resp, err := c.client.CancelOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to cancel order: %s", resp.Message)
	}

	return nil
}

// UpdateOrderStatus updates the status of an order
func (c *OrderClient) UpdateOrderStatus(ctx context.Context, orderID, status string) (*pb.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.UpdateOrderStatusRequest{
		OrderId: orderID,
		Status:  status,
	}

	resp, err := c.client.UpdateOrderStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return resp.Order, nil
}

// WithToken adds an authorization token to the context
func WithToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
}
