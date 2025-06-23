package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/auth/delivery/grpc/proto"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/usecase"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/jwt"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type mockAuthUsecase struct {
	mock.Mock
}

func (m *mockAuthUsecase) Register(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}

func (m *mockAuthUsecase) Login(email, password string) (*usecase.TokenPair, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.TokenPair), args.Error(1)
}

func (m *mockAuthUsecase) RefreshToken(refreshToken string) (*usecase.TokenPair, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.TokenPair), args.Error(1)
}

func (m *mockAuthUsecase) UpdatePassword(userID uuid.UUID, currentPassword, newPassword string) error {
	args := m.Called(userID, currentPassword, newPassword)
	return args.Error(0)
}

func (m *mockAuthUsecase) GetJWKS() ([]jwt.JWK, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jwt.JWK), args.Error(1)
}

func (m *mockAuthUsecase) ValidateToken(token string) (*jwt.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Claims), args.Error(1)
}

func TestAuthServer_Register(t *testing.T) {
	mockUsecase := new(mockAuthUsecase)
	server := NewAuthServer(mockUsecase)

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, server)
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	tests := []struct {
		name          string
		email         string
		password      string
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful registration",
			email:         "test@example.com",
			password:      "password123",
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "user already exists",
			email:         "existing@example.com",
			password:      "password123",
			mockError:     usecase.ErrUserExists,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.On("Register", tt.email, tt.password).Return(tt.mockError).Once()

			req := &pb.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			}

			resp, err := client.Register(ctx, req)

			if tt.expectedError {
				assert.Error(t, err)
				mockUsecase.AssertExpectations(t)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.True(t, resp.Success)
			assert.Equal(t, "User registered successfully", resp.Message)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestAuthServer_Login(t *testing.T) {
	mockUsecase := new(mockAuthUsecase)
	server := NewAuthServer(mockUsecase)

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, server)
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	validTokens := &usecase.TokenPair{
		AccessToken:  "valid_access_token",
		RefreshToken: "valid_refresh_token",
	}

	tests := []struct {
		name          string
		email         string
		password      string
		mockTokens    *usecase.TokenPair
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful login",
			email:         "test@example.com",
			password:      "password123",
			mockTokens:    validTokens,
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "invalid credentials",
			email:         "wrong@example.com",
			password:      "wrongpass",
			mockTokens:    nil,
			mockError:     usecase.ErrInvalidCredentials,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.On("Login", tt.email, tt.password).Return(tt.mockTokens, tt.mockError).Once()

			req := &pb.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			}

			resp, err := client.Login(ctx, req)

			if tt.expectedError {
				assert.Error(t, err)
				mockUsecase.AssertExpectations(t)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.mockTokens.AccessToken, resp.AccessToken)
			assert.Equal(t, tt.mockTokens.RefreshToken, resp.RefreshToken)
			assert.Equal(t, "Bearer", resp.TokenType)
			assert.Equal(t, int32(3600), resp.ExpiresIn)
			mockUsecase.AssertExpectations(t)
		})
	}
}
