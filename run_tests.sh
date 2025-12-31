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