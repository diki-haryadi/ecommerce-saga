package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/usecase"
)

// AuthInterceptor provides authentication for gRPC services
type AuthInterceptor struct {
	authUsecase usecase.AuthUsecase
	// List of methods that don't require authentication
	publicMethods map[string]bool
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(authUsecase usecase.AuthUsecase) *AuthInterceptor {
	// Initialize public methods
	publicMethods := map[string]bool{
		"/auth.AuthService/Register":     true,
		"/auth.AuthService/Login":        true,
		"/auth.AuthService/RefreshToken": true,
		"/auth.AuthService/GetJWKS":      true,
	}

	return &AuthInterceptor{
		authUsecase:   authUsecase,
		publicMethods: publicMethods,
	}
}

// Unary returns a unary server interceptor for authentication
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip authentication for public methods
		if i.publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Get token from metadata
		userID, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		// Add user ID to context
		newCtx := context.WithValue(ctx, "user_id", userID)
		return handler(newCtx, req)
	}
}

// Stream returns a stream server interceptor for authentication
func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Skip authentication for public methods
		if i.publicMethods[info.FullMethod] {
			return handler(srv, ss)
		}

		// Get token from metadata
		userID, err := i.authenticate(ss.Context())
		if err != nil {
			return err
		}

		// Wrap the stream with new context containing user ID
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          context.WithValue(ss.Context(), "user_id", userID),
		}

		return handler(srv, wrappedStream)
	}
}

// authenticate validates the token and returns the user ID
func (i *AuthInterceptor) authenticate(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "invalid authorization format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := i.authUsecase.ValidateToken(token)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}

	return claims.UserID, nil
}

// wrappedServerStream wraps grpc.ServerStream to modify context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
