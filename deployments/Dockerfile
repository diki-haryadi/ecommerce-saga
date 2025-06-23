# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker
RUN CGO_ENABLED=0 GOOS=linux go build -o saga-orchestrator ./cmd/saga-orchestrator

# Development stage
FROM golang:1.21-alpine AS development

# Install development tools
RUN apk add --no-cache git make curl

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy binaries from builder
COPY --from=builder /app/api /app/worker /app/saga-orchestrator ./

# Expose ports
EXPOSE 8080

# Production stage
FROM alpine:latest AS production

# Install CA certificates
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/api /app/worker /app/saga-orchestrator ./

# Copy config files
COPY config ./config

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080

# Set default command (will be overridden by docker-compose)
CMD ["./api"] 