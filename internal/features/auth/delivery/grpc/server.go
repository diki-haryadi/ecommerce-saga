package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/auth/delivery/grpc/proto"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/usecase"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authUsecase usecase.AuthUsecase
}

func NewAuthServer(authUsecase usecase.AuthUsecase) *AuthServer {
	return &AuthServer{
		authUsecase: authUsecase,
	}
}

// Register implements the Register RPC method
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := s.authUsecase.Register(req.Email, req.Password); err != nil {
		switch err {
		case usecase.ErrUserExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case usecase.ErrInvalidCredentials:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

// Login implements the Login RPC method
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	tokens, err := s.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case usecase.ErrInvalidCredentials:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &pb.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// RefreshToken implements the RefreshToken RPC method
func (s *AuthServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	tokens, err := s.authUsecase.RefreshToken(req.RefreshToken)
	if err != nil {
		switch err {
		case usecase.ErrInvalidToken:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// UpdatePassword implements the UpdatePassword RPC method
func (s *AuthServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	// Get user ID from context (assuming it's set by auth interceptor)
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user ID format")
	}

	if err := s.authUsecase.UpdatePassword(uid, req.CurrentPassword, req.NewPassword); err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case usecase.ErrInvalidCredentials:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &pb.UpdatePasswordResponse{
		Success: true,
		Message: "Password updated successfully",
	}, nil
}

// GetJWKS implements the GetJWKS RPC method
func (s *AuthServer) GetJWKS(ctx context.Context, req *pb.GetJWKSRequest) (*pb.GetJWKSResponse, error) {
	jwks, err := s.authUsecase.GetJWKS()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get JWKS")
	}

	pbJWKs := make([]*pb.JWK, len(jwks))
	for i, jwk := range jwks {
		pbJWKs[i] = &pb.JWK{
			Kid: jwk.KeyID,
			Kty: jwk.KeyType,
			Alg: jwk.Algorithm,
			Use: jwk.Use,
			N:   jwk.N,
			E:   jwk.E,
		}
	}

	return &pb.GetJWKSResponse{
		Keys: pbJWKs,
	}, nil
}
