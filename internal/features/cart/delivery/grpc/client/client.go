package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/delivery/grpc/proto"
)

// CartClient represents the gRPC client for cart service
type CartClient struct {
	client pb.CartServiceClient
	conn   *grpc.ClientConn
}

// NewCartClient creates a new cart gRPC client
func NewCartClient(address string) (*CartClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := pb.NewCartServiceClient(conn)
	return &CartClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *CartClient) Close() error {
	return c.conn.Close()
}

// AddItem adds an item to the cart
func (c *CartClient) AddItem(ctx context.Context, userID string, productID string, quantity int32) (*pb.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.AddItemRequest{
		UserId:    userID,
		ProductId: productID,
		Quantity:  quantity,
	}

	resp, err := c.client.AddItem(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to add item: %w", err)
	}

	return resp.Cart, nil
}

// RemoveItem removes an item from the cart
func (c *CartClient) RemoveItem(ctx context.Context, cartItemID string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.RemoveItemRequest{
		CartItemId: cartItemID,
	}

	resp, err := c.client.RemoveItem(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to remove item: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to remove item: %s", resp.Message)
	}

	return nil
}

// UpdateItem updates an item in the cart
func (c *CartClient) UpdateItem(ctx context.Context, userID string, cartItemID string, quantity int32) (*pb.Cart, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.UpdateItemRequest{
		UserId:     userID,
		CartItemId: cartItemID,
		Quantity:   quantity,
	}

	resp, err := c.client.UpdateItem(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return resp.Cart, nil
}

// GetCart retrieves the current cart
func (c *CartClient) GetCart(ctx context.Context) (*pb.GetCartResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.GetCartRequest{}
	return c.client.GetCart(ctx, req)
}

// ClearCart removes all items from the cart
func (c *CartClient) ClearCart(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.ClearCartRequest{}
	resp, err := c.client.ClearCart(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to clear cart: %s", resp.Message)
	}

	return nil
}

// WithToken adds an authorization token to the context
func WithToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
}
