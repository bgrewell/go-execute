#!/bin/bash

set -e

# Run unit tests
echo "Running unit tests..."
go test ./... -v -cover -coverprofile=coverage.out

# Display coverage in the terminal
echo "Displaying coverage in the terminal..."
go tool cover -func=coverage.out

# Run unit tests with race detector
echo "Running unit tests with race detector..."
go test ./... -race -v

# Run integration tests
echo "Running integration tests..."
go test -tags=integration ./... -v

# Run end-to-end tests
echo "Running end-to-end tests..."
go test -tags=e2e ./... -v
