# Saga Pattern Implementation Guide

## Overview
The Saga pattern is implemented in this service to manage distributed transactions across multiple microservices. It provides a robust way to maintain data consistency and handle failures in distributed systems.

## Core Components

### Status Types
```go
const (
    StatusPending     Status = "PENDING"
    StatusInProgress  Status = "IN_PROGRESS"
    StatusCompleted   Status = "COMPLETED"
    StatusFailed      Status = "FAILED"
    StatusCompensated Status = "COMPENSATED"
)
```

These statuses track the lifecycle of a saga transaction:
- `PENDING`: Initial state when saga is created
- `IN_PROGRESS`: Saga steps are being executed
- `COMPLETED`: All steps completed successfully
- `FAILED`: A step has failed and compensation may be needed
- `COMPENSATED`: All compensation actions completed

### Main Interface
```go
type Usecase interface {
    StartOrderSaga(ctx context.Context, orderID, userID uuid.UUID, amount float64, paymentMethod string, metadata map[string]string) (*SagaResponse, error)
    GetSagaStatus(ctx context.Context, sagaID uuid.UUID) (*SagaResponse, error)
    CompensateTransaction(ctx context.Context, sagaID, stepID uuid.UUID, reason string) (*SagaResponse, error)
    ListSagaTransactions(ctx context.Context, page, limit int32, status string, sagaType Type) ([]*SagaResponse, int64, error)
}
```

## How to Use

### 1. Starting a Saga Transaction

```go
sagaResponse, err := sagaService.StartOrderSaga(
    ctx,
    orderID,
    userID,
    100.00,              // amount
    "credit_card",       // payment method
    map[string]string{   // metadata
        "order_type": "regular",
        "currency": "USD",
    },
)
```

### 2. Saga Steps Structure
```go
type SagaStep struct {
    ID                 uuid.UUID              
    Name               string                 
    Service            string                 
    Action             string                 
    CompensationAction string                 
    Payload            map[string]interface{} 
    ErrorMessage       string                 
    ExecutedAt         time.Time              
}
```

### 3. Typical Order Flow Steps

1. **Order Creation**
   - Action: Create order
   - Compensation: Cancel order
   - Service: Order Service

2. **Payment Processing**
   - Action: Process payment
   - Compensation: Refund payment
   - Service: Payment Service

3. **Inventory Update**
   - Action: Reserve inventory
   - Compensation: Release inventory
   - Service: Inventory Service

### 4. Error Handling and Compensation

```go
// Trigger compensation
compensatedSaga, err := sagaService.CompensateTransaction(
    ctx,
    sagaID,
    failedStepID,
    "Payment processing failed"
)
```

### 5. Monitoring Saga Status

```go
// Check specific saga status
status, err := sagaService.GetSagaStatus(ctx, sagaID)

// List transactions
transactions, total, err := sagaService.ListSagaTransactions(
    ctx,
    1,                  // page
    10,                 // limit
    "IN_PROGRESS",      // status
    saga.TypeOrder,     // type
)
```

## Best Practices

### 1. Idempotency
- Ensure all operations are idempotent
- Use unique transaction IDs
- Handle duplicate requests gracefully
- Implement proper retry mechanisms

### 2. Monitoring
- Track saga progress
- Log compensation actions
- Monitor step execution times
- Set up alerts for failed sagas

### 3. Error Handling
```go
var (
    ErrNotFound      = NewError("saga not found")
    ErrAlreadyExists = NewError("saga already exists")
    ErrInvalidStep   = NewError("invalid step order")
)
```

## Complete Integration Example

```go
func ProcessOrder(ctx context.Context, orderDetails OrderDetails) error {
    // 1. Start saga
    sagaResp, err := sagaService.StartOrderSaga(
        ctx,
        orderDetails.OrderID,
        orderDetails.UserID,
        orderDetails.Amount,
        orderDetails.PaymentMethod,
        orderDetails.Metadata,
    )
    if err != nil {
        return err
    }

    // 2. Monitor saga progress
    for {
        status, err := sagaService.GetSagaStatus(ctx, sagaResp.ID)
        if err != nil {
            return err
        }

        switch status.Status {
        case saga.StatusCompleted:
            return nil
        case saga.StatusFailed, saga.StatusCompensated:
            return fmt.Errorf("saga failed: %s", status.Steps[len(status.Steps)-1].ErrorMessage)
        case saga.StatusPending, saga.StatusInProgress:
            time.Sleep(time.Second)
            continue
        }
    }
}
```

## Benefits

1. **Distributed Transaction Management**
   - Maintains data consistency across services
   - Handles long-running transactions
   - Provides clear transaction boundaries

2. **Automatic Compensation**
   - Rolls back failed transactions
   - Maintains system consistency
   - Handles partial failures gracefully

3. **Monitoring and Observability**
   - Tracks transaction progress
   - Provides audit trail
   - Enables debugging and troubleshooting

4. **Scalability**
   - Supports microservices architecture
   - Handles concurrent transactions
   - Enables service independence

## Common Pitfalls to Avoid

1. **Non-idempotent Operations**
   - Ensure all steps can be safely retried
   - Implement proper deduplication
   - Handle edge cases in compensation

2. **Insufficient Monitoring**
   - Set up comprehensive logging
   - Implement proper alerting
   - Track performance metrics

3. **Incomplete Error Handling**
   - Handle all possible error scenarios
   - Implement proper timeout handling
   - Consider network failures

4. **Poor Compensation Design**
   - Ensure compensation actions are complete
   - Test compensation scenarios
   - Handle compensation failures

## Testing Recommendations

1. **Unit Tests**
   - Test individual saga steps
   - Verify compensation logic
   - Test error scenarios

2. **Integration Tests**
   - Test complete saga flows
   - Verify cross-service interactions
   - Test timeout scenarios

3. **Chaos Testing**
   - Test network failures
   - Test service unavailability
   - Test partial failures

## Conclusion

The Saga pattern implementation in this service provides a robust solution for managing distributed transactions. By following these guidelines and best practices, you can ensure reliable and consistent operation of your distributed systems. 