#!/bin/bash
# Setup script for running tests

set -e

echo "Setting up test environment..."

# Download dependencies
echo "Downloading Go dependencies..."
go mod download
go mod tidy

# Verify dependencies
echo "Verifying dependencies..."
go mod verify

# Check for test tools
echo "Checking for test tools..."
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

echo "Setup complete!"
echo ""
echo "Run tests with: go test ./..."
echo "Run linter with: golangci-lint run"

