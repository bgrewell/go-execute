#!/bin/bash

set -e

# Run unit tests
# echo "Running unit tests..."
# go test ./... -v

# Run unit tests with race detector
# echo "Running unit tests with race detector..."
# go test ./... -race -v

# Run integration tests
echo "Running integration tests..."
go test -tags=integration ./... -v

# Run end-to-end tests
echo "Running end-to-end tests..."
go test -tags=e2e ./... -v
