# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the applications
RUN make build

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binaries from builder
COPY --from=builder /app/bin/api /app/bin/api
COPY --from=builder /app/bin/worker /app/bin/worker
COPY --from=builder /app/bin/orchestrator /app/bin/orchestrator

# Copy config files
COPY --from=builder /app/config /app/config

# Create non-root user
RUN adduser -D appuser
USER appuser

# Set environment variables
ENV GIN_MODE=release

# Expose ports
EXPOSE 8080

# Default command (can be overridden)
CMD ["/app/bin/api"] 