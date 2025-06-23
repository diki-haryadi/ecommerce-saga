#!/bin/bash

# Generate protobuf code for cart service
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/features/cart/delivery/grpc/proto/cart.proto 