package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/auth/delivery/grpc/proto"
)

// AuthClient represents the gRPC client for auth service
type AuthClient struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

// NewAuthClient creates a new auth gRPC client
func NewAuthClient(address string) (*AuthClient, error) {
	// Set up a connection to the server with insecure transport credentials
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *AuthClient) Close() error {
	return c.conn.Close()
}

// Register registers a new user
func (c *AuthClient) Register(ctx context.Context, email, password string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.RegisterRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.client.Register(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("registration failed: %s", resp.Message)
	}

	return nil
}

// Login authenticates a user and returns tokens
func (c *AuthClient) Login(ctx context.Context, email, password string) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.LoginRequest{
		Email:    email,
		Password: password,
	}

	return c.client.Login(ctx, req)
}

// RefreshToken refreshes the access token
func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*pb.RefreshTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	return c.client.RefreshToken(ctx, req)
}

// UpdatePassword updates the user's password
func (c *AuthClient) UpdatePassword(ctx context.Context, currentPassword, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.UpdatePasswordRequest{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}

	resp, err := c.client.UpdatePassword(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("password update failed: %s", resp.Message)
	}

	return nil
}

// GetJWKS retrieves the JSON Web Key Set
func (c *AuthClient) GetJWKS(ctx context.Context) ([]*pb.JWK, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req := &pb.GetJWKSRequest{}
	resp, err := c.client.GetJWKS(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %w", err)
	}

	return resp.Keys, nil
}

// WithToken adds an authorization token to the context
func WithToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
}
