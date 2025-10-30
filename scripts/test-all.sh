#!/bin/bash
# Run all tests with proper setup

set -e

echo "🧪 Running mautrix-viber test suite..."
echo ""

# Setup
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

echo ""
echo "✅ Running tests..."
echo ""

# Run tests with coverage
go test -v -cover ./...

echo ""
echo "📊 Generating coverage report..."
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "✨ Tests complete! Coverage report saved to coverage.html"

