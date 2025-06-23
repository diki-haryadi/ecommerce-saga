.PHONY: build test clean run-api run-worker run-orchestrator migrate-up migrate-down lint docker-build docker-up proto proto-clean proto-all run-grpc

# Build commands
build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker
	go build -o bin/orchestrator ./cmd/saga-orchestrator

# Run commands
run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

run-orchestrator:
	go run ./cmd/saga-orchestrator

# Test commands
test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Database migrations
migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

# Linting
lint:
	golangci-lint run

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up

docker-down:
	docker-compose down

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate swagger docs
swagger:
	swag init -g cmd/api/main.go -o api/docs

# Protobuf commands
PROTO_DIR=proto
GO_OUT_DIR=internal/features

.PHONY: proto
proto: proto-auth proto-cart proto-order proto-payment proto-saga

.PHONY: proto-auth
proto-auth:
	@echo "Generating auth proto..."
	protoc --go_out=. \
		--go_opt=module=github.com/diki-haryadi/ecommerce-saga \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/diki-haryadi/ecommerce-saga \
		$(PROTO_DIR)/auth/auth.proto

.PHONY: proto-cart
proto-cart:
	@echo "Generating cart proto..."
	protoc --go_out=. \
		--go_opt=module=github.com/diki-haryadi/ecommerce-saga \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/diki-haryadi/ecommerce-saga \
		$(PROTO_DIR)/cart/cart.proto

.PHONY: proto-order
proto-order:
	@echo "Generating order proto..."
	protoc --go_out=. \
		--go_opt=module=github.com/diki-haryadi/ecommerce-saga \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/diki-haryadi/ecommerce-saga \
		$(PROTO_DIR)/order/order.proto

.PHONY: proto-payment
proto-payment:
	@echo "Generating payment proto..."
	protoc --go_out=. \
		--go_opt=module=github.com/diki-haryadi/ecommerce-saga \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/diki-haryadi/ecommerce-saga \
		$(PROTO_DIR)/payment/payment.proto

.PHONY: proto-saga
proto-saga:
	@echo "Generating saga proto..."
	protoc --go_out=. \
		--go_opt=module=github.com/diki-haryadi/ecommerce-saga \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/diki-haryadi/ecommerce-saga \
		$(PROTO_DIR)/saga/saga.proto

proto-clean: ## Clean generated protobuf code
	@echo "Cleaning generated protobuf code..."
	@find . -name "*.pb.go" -type f -delete

proto-all: proto-clean proto ## Clean and regenerate all protobuf code

# New target
run-grpc: ## Run the gRPC server
	@echo "Starting gRPC server..."
	@go run cmd/grpc/main.go

.PHONY: clean-proto
clean-proto:
	@echo "Cleaning generated proto files..."
	rm -rf $(GO_OUT_DIR)/*/delivery/grpc/proto

.DEFAULT_GOAL := build 