package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cartHttp "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/delivery/http"
	cartRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/repository/postgres"
	cartUsecase "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/usecase"
	orderRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/order/repository/postgres"
	sagaHandler "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/delivery/http"
	sagaRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/repository/postgres"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/usecase"
	"github.com/diki-haryadi/ecommerce-saga/test/integration/testutil"
)

// MockProductService implements cartUsecase.ProductService for testing
type MockProductService struct{}

func (m *MockProductService) GetProduct(ctx context.Context, id uuid.UUID) (*cartUsecase.Product, error) {
	return &cartUsecase.Product{
		ID:    id,
		Name:  "Test Product",
		Price: 10.0,
		Stock: 100,
	}, nil
}

func setupTestApp(t *testing.T) (*fiber.App, *testutil.TestDB) {
	// Setup database
	tdb := testutil.NewTestDB(t)

	// Initialize repositories
	cartRepository := cartRepo.NewCartRepository(tdb.DB)
	orderRepository := orderRepo.NewOrderRepository(tdb.DB)
	sagaRepository := sagaRepo.NewSagaRepository(tdb.DB)

	// Setup Fiber app
	app := fiber.New()

	// Setup routes
	api := app.Group("/api")

	// Initialize cart usecase and handler
	cartUsecase := cartUsecase.NewCartUsecase(cartRepository, &MockProductService{}, 24*time.Hour)
	cartHandler := cartHttp.NewCartHandler(cartUsecase)
	cartGroup := api.Group("/cart")
	cartGroup.Get("/:user_id", cartHandler.GetCart)
	cartGroup.Post("/items", cartHandler.AddItem)
	cartGroup.Delete("/:user_id", cartHandler.ClearCart)

	// Initialize saga usecase and handler
	sagaUsecase := usecase.NewSagaUsecase(sagaRepository, orderRepository, nil, nil, nil, nil)
	sagaHandler := sagaHandler.NewSagaHandler(sagaUsecase)
	sagaGroup := api.Group("/saga")
	sagaGroup.Post("/order-payment", sagaHandler.StartOrderPaymentSaga)
	sagaGroup.Get("/:id", sagaHandler.GetSagaStatus)

	return app, tdb
}

func TestCartAPI(t *testing.T) {
	app, tdb := setupTestApp(t)
	defer tdb.Cleanup()

	t.Run("create and get cart", func(t *testing.T) {
		userID := uuid.New()

		// Get cart (will create if not exists)
		resp, err := app.Test(httptest.NewRequest("GET", fmt.Sprintf("/api/cart/%s", userID), nil))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var cartResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&cartResp)
		require.NoError(t, err)
		assert.NotNil(t, cartResp["data"])
	})

	t.Run("add item to cart", func(t *testing.T) {
		// Add item to cart
		addItemReq := fiber.Map{
			"product_id": uuid.New().String(),
			"quantity":   2,
		}
		req, _ := json.Marshal(addItemReq)
		resp, err := app.Test(httptest.NewRequest("POST", "/api/cart/items", bytes.NewReader(req)))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestSagaAPI(t *testing.T) {
	app, tdb := setupTestApp(t)
	defer tdb.Cleanup()

	t.Run("start and monitor saga", func(t *testing.T) {
		orderID := uuid.New()

		// Start saga
		req := fiber.Map{
			"order_id": orderID.String(),
		}
		reqBody, _ := json.Marshal(req)
		resp, err := app.Test(httptest.NewRequest("POST", "/api/saga/order-payment", bytes.NewReader(reqBody)))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check saga status
		resp, err = app.Test(httptest.NewRequest("GET", fmt.Sprintf("/api/saga/%s", orderID), nil))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var sagaStatus map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&sagaStatus)
		require.NoError(t, err)
		assert.Contains(t, []string{"PENDING", "COMPLETED", "FAILED"}, sagaStatus["status"])
	})
}
