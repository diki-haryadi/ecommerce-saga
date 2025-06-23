.PHONY: build test clean run-api run-worker run-orchestrator migrate-up migrate-down lint docker-build docker-up

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

.DEFAULT_GOAL := build 