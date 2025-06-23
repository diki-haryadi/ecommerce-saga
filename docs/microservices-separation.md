# Separating Features into Standalone Microservices

This guide explains how to separate the features from this monolithic codebase into standalone microservices.

## Current Structure

The current codebase is organized by features under `internal/features/`:
- Cart Service
- Order Service
- Payment Service
- Saga Orchestration

Each feature is self-contained with its own:
- Proto definitions
- gRPC server/client implementations
- Business logic (usecase)
- Domain models
- Clear interfaces

## Separation Process

### 1. Creating a New Microservice

Using the cart service as an example:

```bash
# Create new project
mkdir cart-service
cd cart-service

# Initialize Go module
go mod init github.com/your-org/cart-service

# Create directory structure
mkdir -p api/proto \
        internal/{delivery/grpc,usecase,repository,domain} \
        cmd/server \
        config \
        deployments
```

### 2. Project Structure

```
cart-service/
├── api/
│   └── proto/
│       └── cart.proto
├── internal/
│   ├── delivery/
│   │   └── grpc/
│   │       ├── server.go
│   │       └── client/
│   │           └── client.go
│   ├── usecase/
│   │   ├── cart_usecase.go
│   │   └── interface.go
│   ├── repository/
│   │   ├── postgres/
│   │   │   └── cart_repository.go
│   │   └── interface.go
│   └── domain/
│       ├── entity/
│       │   └── cart.go
│       └── service/
│           └── cart_service.go
├── cmd/
│   └── server/
│       └── main.go
├── config/
│   └── config.yaml
├── deployments/
│   ├── Dockerfile
│   └── docker-compose.yml
├── go.mod
└── Makefile
```

### 3. Implementation Steps

1. **Copy Required Components**
   - Proto definitions from `proto/cart`
   - Feature implementation from `internal/features/cart`
   - Required shared utilities

2. **Update Main Service**

```go
// cmd/server/main.go
package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
    "github.com/your-org/cart-service/internal/delivery/grpc"
    "github.com/your-org/cart-service/internal/usecase"
    "github.com/your-org/cart-service/internal/repository"
)

func main() {
    // Initialize dependencies
    repo := repository.NewCartRepository()
    usecase := usecase.NewCartUsecase(repo)
    server := grpc.NewCartServer(usecase)

    // Start gRPC server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterCartServiceServer(grpcServer, server)
    
    log.Printf("Starting cart service on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
```

3. **Configuration**

```yaml
# config/config.yaml
server:
  grpc:
    port: 50051
  http:
    port: 8080
database:
  host: localhost
  port: 5432
  name: cart_db
  user: postgres
  password: secret
```

4. **Docker Configuration**

```dockerfile
# deployments/Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /cart-service ./cmd/server

EXPOSE 50051

CMD ["/cart-service"]
```

```yaml
# deployments/docker-compose.yml
version: '3'

services:
  cart:
    build:
      context: ..
      dockerfile: deployments/Dockerfile
    ports:
      - "50051:50051"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=cart_db
      - DB_USER=postgres
      - DB_PASSWORD=secret
    depends_on:
      - postgres

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=cart_db
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=secret
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

5. **Makefile**

```makefile
.PHONY: build run test proto docker-build docker-up

build:
	go build -o bin/cart-service ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./...

proto:
	protoc --go_out=. \
		--go_opt=module=github.com/your-org/cart-service \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/your-org/cart-service \
		api/proto/cart.proto

docker-build:
	docker-compose -f deployments/docker-compose.yml build

docker-up:
	docker-compose -f deployments/docker-compose.yml up

docker-down:
	docker-compose -f deployments/docker-compose.yml down
```

### 4. Service Communication

Other services will communicate with the cart service using the gRPC client:

```go
// Example from order service
cartClient, err := cartclient.NewCartClient("cart-service:50051")
if err != nil {
    return err
}
defer cartClient.Close()

cart, err := cartClient.GetCart(ctx)
if err != nil {
    return err
}
```

### 5. Service Discovery

Add service discovery using Consul:

```go
import "github.com/hashicorp/consul/api"

func registerService() error {
    config := api.DefaultConfig()
    client, err := api.NewClient(config)
    if err != nil {
        return err
    }

    registration := &api.AgentServiceRegistration{
        ID:      "cart-service-1",
        Name:    "cart-service",
        Port:    50051,
        Address: "localhost",
        Check: &api.AgentServiceCheck{
            GRPC:     "localhost:50051",
            Interval: "10s",
        },
    }

    return client.Agent().ServiceRegister(registration)
}
```

## Benefits of Separation

1. **Independent Deployment**: Each service can be deployed independently
2. **Scalability**: Services can be scaled based on their specific needs
3. **Technology Freedom**: Each service can use different technologies if needed
4. **Team Autonomy**: Different teams can work on different services
5. **Fault Isolation**: Issues in one service don't directly affect others

## Considerations

1. **Database Per Service**: Each service should have its own database
2. **Service Discovery**: Implement proper service discovery mechanism
3. **Monitoring**: Add monitoring and tracing for distributed systems
4. **Authentication**: Implement service-to-service authentication
5. **Circuit Breaking**: Add circuit breakers for resilience

## Migration Strategy

1. **Gradual Migration**:
   - Start with non-critical services
   - Migrate one service at a time
   - Run old and new services in parallel during transition

2. **Data Migration**:
   - Plan data migration strategy
   - Consider dual-write period
   - Validate data consistency

3. **Testing**:
   - Test service in isolation
   - Test service integration
   - Perform load testing

4. **Rollback Plan**:
   - Maintain ability to rollback to monolith
   - Keep old service running until new service is stable
   - Monitor for issues during transition 