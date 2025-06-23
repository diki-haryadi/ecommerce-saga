# E-Commerce Saga System

A feature-based microservices system implementing the Saga pattern for distributed transactions.

## Features

- Authentication with JWK (JSON Web Key)
- Cart Management
- Order Processing
- Payment Processing
- Saga Orchestration

## Architecture

This project follows a Feature-Based Structure with Clean Architecture principles:

- API Server
- Worker Service
- Saga Orchestrator
- Multiple Message Brokers (RabbitMQ, Kafka, NSQ, NATS)
- Dual Database System (PostgreSQL & MongoDB)

## Prerequisites

- Go 1.21+
- PostgreSQL 14+
- MongoDB 5+
- RabbitMQ 3.9+
- Apache Kafka 3.0+
- NSQ
- NATS Server

## Project Structure

```
ecommerce-saga/
├── cmd/                    # Application entry points
│   ├── api/               # API server
│   ├── worker/            # Background worker
│   └── saga-orchestrator/ # Saga orchestration service
├── internal/              # Private application code
│   ├── features/          # Business features
│   │   ├── auth/         # Authentication feature
│   │   ├── cart/         # Shopping cart feature
│   │   ├── order/        # Order management
│   │   ├── payment/      # Payment processing
│   │   └── saga/         # Saga coordination
│   ├── shared/           # Shared code
│   └── infrastructure/   # Infrastructure code
├── pkg/                  # Public libraries
├── config/              # Configuration files
├── migrations/          # Database migrations
└── deployments/        # Deployment configurations
```

## Getting Started

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up environment variables (copy from .env.example)
4. Run the migrations:
   ```bash
   make migrate-up
   ```
5. Start the services:
   ```bash
   make run-api
   make run-worker
   make run-orchestrator
   ```

## Development

### Running Tests
```bash
make test
```

### Running Linter
```bash
make lint
```

### Building
```bash
make build
```

## Docker

Build and run with Docker Compose:
```bash
docker-compose up --build
```

## API Documentation

API documentation is available at `/swagger/index.html` when running in development mode.

## License

MIT License 