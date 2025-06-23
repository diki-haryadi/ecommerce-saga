package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga"
	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/delivery/grpc/proto"
)

type SagaServer struct {
	pb.UnimplementedSagaServiceServer
	sagaUsecase saga.Usecase
}

func NewSagaServer(sagaUsecase saga.Usecase) *SagaServer {
	return &SagaServer{
		sagaUsecase: sagaUsecase,
	}
}

func (s *SagaServer) StartOrderSaga(ctx context.Context, req *pb.StartOrderSagaRequest) (*pb.StartOrderSagaResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	sagaResp, err := s.sagaUsecase.StartOrderSaga(ctx, orderID, userID, req.Amount, req.PaymentMethod, req.Metadata)
	if err != nil {
		var errStatus error
		switch err {
		case saga.ErrAlreadyExists:
			errStatus = status.Error(codes.AlreadyExists, err.Error())
		default:
			errStatus = status.Error(codes.Internal, "failed to start saga")
		}
		return nil, errStatus
	}

	return &pb.StartOrderSagaResponse{
		Success:     true,
		Message:     "Saga started successfully",
		Transaction: convertSagaTransactionToPb(sagaResp),
	}, nil
}

func (s *SagaServer) GetSagaStatus(ctx context.Context, req *pb.GetSagaStatusRequest) (*pb.GetSagaStatusResponse, error) {
	sagaID, err := uuid.Parse(req.SagaId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid saga ID")
	}

	sagaResp, err := s.sagaUsecase.GetSagaStatus(ctx, sagaID)
	if err != nil {
		var errStatus error
		switch err {
		case saga.ErrNotFound:
			errStatus = status.Error(codes.NotFound, err.Error())
		default:
			errStatus = status.Error(codes.Internal, "failed to get saga status")
		}
		return nil, errStatus
	}

	return &pb.GetSagaStatusResponse{
		Transaction: convertSagaTransactionToPb(sagaResp),
	}, nil
}

func (s *SagaServer) CompensateTransaction(ctx context.Context, req *pb.CompensateTransactionRequest) (*pb.CompensateTransactionResponse, error) {
	sagaID, err := uuid.Parse(req.SagaId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid saga ID")
	}

	stepID, err := uuid.Parse(req.StepId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid step ID")
	}

	sagaResp, err := s.sagaUsecase.CompensateTransaction(ctx, sagaID, stepID, req.Reason)
	if err != nil {
		var errStatus error
		switch err {
		case saga.ErrNotFound:
			errStatus = status.Error(codes.NotFound, err.Error())
		case saga.ErrInvalidStep:
			errStatus = status.Error(codes.FailedPrecondition, err.Error())
		default:
			errStatus = status.Error(codes.Internal, "failed to compensate transaction")
		}
		return nil, errStatus
	}

	return &pb.CompensateTransactionResponse{
		Success:     true,
		Message:     "Transaction compensated successfully",
		Transaction: convertSagaTransactionToPb(sagaResp),
	}, nil
}

func (s *SagaServer) ListSagaTransactions(ctx context.Context, req *pb.ListSagaTransactionsRequest) (*pb.ListSagaTransactionsResponse, error) {
	transactions, total, err := s.sagaUsecase.ListSagaTransactions(ctx, req.Page, req.Limit, req.Status, saga.Type(req.Type))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list transactions")
	}

	pbTransactions := make([]*pb.SagaTransaction, len(transactions))
	for i, t := range transactions {
		pbTransactions[i] = convertSagaTransactionToPb(t)
	}

	return &pb.ListSagaTransactionsResponse{
		Transactions: pbTransactions,
		Total:        int32(total),
		Page:         req.Page,
		Limit:        req.Limit,
	}, nil
}

func convertSagaTransactionToPb(s *saga.SagaResponse) *pb.SagaTransaction {
	steps := make([]*pb.SagaStep, len(s.Steps))
	for i, step := range s.Steps {
		// Convert map[string]interface{} to map[string]string
		payload := make(map[string]string)
		for k, v := range step.Payload {
			if str, ok := v.(string); ok {
				payload[k] = str
			}
		}

		steps[i] = &pb.SagaStep{
			Id:                 step.ID.String(),
			Name:               step.Name,
			Status:             string(step.Status),
			Service:            step.Service,
			Action:             step.Action,
			CompensationAction: step.CompensationAction,
			Payload:            payload,
			ErrorMessage:       step.ErrorMessage,
			ExecutedAt:         timestamppb.New(step.ExecutedAt),
		}
	}

	return &pb.SagaTransaction{
		Id:        s.ID.String(),
		Type:      string(s.Type),
		Status:    string(s.Status),
		Steps:     steps,
		Metadata:  s.Metadata,
		CreatedAt: timestamppb.New(s.CreatedAt),
		UpdatedAt: timestamppb.New(s.UpdatedAt),
	}
}
