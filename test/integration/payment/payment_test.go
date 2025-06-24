package payment_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/test/integration/testutil"
)

func TestPaymentRepository(t *testing.T) {
	// Setup test database
	tdb := testutil.NewTestDB(t)
	defer tdb.Cleanup()

	// Clean up tables before test
	require.NoError(t, tdb.TruncateTables("payments"))

	// Initialize repository
	paymentRepo := postgres.NewPaymentRepository(tdb.DB)

	t.Run("create and retrieve payment", func(t *testing.T) {
		// Create test payment
		payment := &entity.Payment{
			ID:       uuid.New(),
			OrderID:  uuid.New(),
			Amount:   100.0,
			Currency: "USD",
			Status:   entity.PaymentStatusPending,
			Provider: entity.PaymentProviderStripe,
		}

		// Create payment
		err := paymentRepo.Create(context.Background(), payment)
		require.NoError(t, err)

		// Retrieve payment
		retrieved, err := paymentRepo.GetByID(context.Background(), payment.ID)
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, payment.ID, retrieved.ID)
		assert.Equal(t, payment.OrderID, retrieved.OrderID)
		assert.Equal(t, payment.Amount, retrieved.Amount)
		assert.Equal(t, payment.Status, retrieved.Status)
	})

	t.Run("update payment status", func(t *testing.T) {
		// Create test payment
		payment := &entity.Payment{
			ID:       uuid.New(),
			OrderID:  uuid.New(),
			Amount:   100.0,
			Currency: "USD",
			Status:   entity.PaymentStatusPending,
			Provider: entity.PaymentProviderStripe,
		}

		// Create payment
		err := paymentRepo.Create(context.Background(), payment)
		require.NoError(t, err)

		// Update status
		payment.Status = entity.PaymentStatusSuccess
		err = paymentRepo.Update(context.Background(), payment)
		require.NoError(t, err)

		// Verify status update
		updated, err := paymentRepo.GetByID(context.Background(), payment.ID)
		require.NoError(t, err)
		assert.Equal(t, entity.PaymentStatusSuccess, updated.Status)
	})

	t.Run("get payment by order ID", func(t *testing.T) {
		orderID := uuid.New()
		payment := &entity.Payment{
			ID:       uuid.New(),
			OrderID:  orderID,
			Amount:   100.0,
			Currency: "USD",
			Status:   entity.PaymentStatusPending,
			Provider: entity.PaymentProviderStripe,
		}

		// Create payment
		err := paymentRepo.Create(context.Background(), payment)
		require.NoError(t, err)

		// Retrieve by order ID
		retrieved, err := paymentRepo.GetByOrderID(context.Background(), orderID)
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, payment.ID, retrieved.ID)
		assert.Equal(t, orderID, retrieved.OrderID)
	})
}
