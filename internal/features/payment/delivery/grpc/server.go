package grpc

import (
	"context"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/usecase"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/delivery/grpc/proto"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
	paymentUsecase usecase.Usecase
}

func NewPaymentServer(paymentUsecase usecase.Usecase) *PaymentServer {
	return &PaymentServer{
		paymentUsecase: paymentUsecase,
	}
}

func (s *PaymentServer) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	paymentResp, err := s.paymentUsecase.CreatePayment(ctx, orderID, req.Amount, req.Currency, req.PaymentMethod)
	if err != nil {
		var errStatus error
		switch err {
		case usecase.ErrOrderNotFound:
			errStatus = status.Error(codes.NotFound, err.Error())
		case usecase.ErrInvalidProvider:
			errStatus = status.Error(codes.InvalidArgument, err.Error())
		default:
			errStatus = status.Error(codes.Internal, "failed to create payment")
		}
		return nil, errStatus
	}

	return &pb.CreatePaymentResponse{
		Success: true,
		Message: "Payment created successfully",
		Payment: convertPaymentToPb(paymentResp),
	}, nil
}

func (s *PaymentServer) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	paymentID, err := uuid.Parse(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment ID")
	}

	paymentResp, err := s.paymentUsecase.GetPayment(ctx, paymentID)
	if err != nil {
		var errStatus error
		switch err {
		case usecase.ErrNotFound:
			errStatus = status.Error(codes.NotFound, err.Error())
		default:
			errStatus = status.Error(codes.Internal, "failed to get payment")
		}
		return nil, errStatus
	}

	return &pb.GetPaymentResponse{
		Payment: convertPaymentToPb(paymentResp),
	}, nil
}

func (s *PaymentServer) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user ID format")
	}

	payments, total, err := s.paymentUsecase.ListPayments(ctx, uid, req.Page, req.Limit, req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list payments")
	}

	pbPayments := make([]*pb.Payment, len(payments))
	for i, p := range payments {
		pbPayments[i] = convertPaymentToPb(p)
	}

	return &pb.ListPaymentsResponse{
		Payments: pbPayments,
		Total:    int32(total),
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

func (s *PaymentServer) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	paymentID, err := uuid.Parse(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment ID")
	}

	details := &usecase.PaymentDetails{
		CardNumber:  req.PaymentDetails.CardNumber,
		ExpiryMonth: req.PaymentDetails.ExpiryMonth,
		ExpiryYear:  req.PaymentDetails.ExpiryYear,
		CVV:         req.PaymentDetails.Cvv,
		HolderName:  req.PaymentDetails.HolderName,
	}

	paymentResp, err := s.paymentUsecase.ProcessPayment(ctx, paymentID, details)
	if err != nil {
		var errStatus error
		switch err {
		case usecase.ErrNotFound:
			errStatus = status.Error(codes.NotFound, err.Error())
		case usecase.ErrCompleted:
			errStatus = status.Error(codes.FailedPrecondition, err.Error())
		case usecase.ErrProviderUnavailable:
			errStatus = status.Error(codes.Unavailable, err.Error())
		default:
			errStatus = status.Error(codes.Internal, "failed to process payment")
		}
		return nil, errStatus
	}

	return &pb.ProcessPaymentResponse{
		Success:       true,
		Message:       "Payment processed successfully",
		TransactionId: paymentResp.ProviderTransactionID,
		Payment:       convertPaymentToPb(paymentResp),
	}, nil
}

func (s *PaymentServer) RefundPayment(ctx context.Context, req *pb.RefundPaymentRequest) (*pb.RefundPaymentResponse, error) {
	paymentID, err := uuid.Parse(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment ID")
	}

	paymentResp, refundID, err := s.paymentUsecase.RefundPayment(ctx, paymentID, req.Amount, req.Reason)
	if err != nil {
		var errStatus error
		switch err {
		case usecase.ErrNotFound:
			errStatus = status.Error(codes.NotFound, err.Error())
		case usecase.ErrInvalidStatus:
			errStatus = status.Error(codes.FailedPrecondition, err.Error())
		case usecase.ErrProviderUnavailable:
			errStatus = status.Error(codes.Unavailable, err.Error())
		default:
			errStatus = status.Error(codes.Internal, "failed to refund payment")
		}
		return nil, errStatus
	}

	return &pb.RefundPaymentResponse{
		Success:  true,
		Message:  "Payment refunded successfully",
		RefundId: refundID,
		Payment:  convertPaymentToPb(paymentResp),
	}, nil
}

func convertPaymentToPb(p *usecase.PaymentResponse) *pb.Payment {
	return &pb.Payment{
		Id:            p.ID.String(),
		OrderId:       p.OrderID.String(),
		UserId:        p.UserID.String(),
		Amount:        p.Amount,
		Currency:      p.Currency,
		Status:        string(p.Status),
		PaymentMethod: p.PaymentMethod,
		TransactionId: p.ProviderTransactionID,
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}
}
