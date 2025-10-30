# Makefile for mautrix-viber bridge

.PHONY: build run test lint clean docker docker-run help

# Build the bridge
build:
	go build -o bin/mautrix-viber ./cmd/mautrix-viber

# Run the bridge
run: build
	./bin/mautrix-viber

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linters
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Build Docker image
docker:
	docker build -t mautrix-viber:latest .

# Run Docker container
docker-run:
	docker-compose up -d

# Stop Docker container
docker-stop:
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# Database migration
migrate:
	@echo "Run migrations manually using the migration tool"

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate documentation
docs:
	godoc -http=:6060

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the bridge binary"
	@echo "  run           - Build and run the bridge"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint          - Run linters"
	@echo "  fmt           - Format code"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker        - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-stop   - Stop Docker container"
	@echo "  logs          - View Docker logs"
	@echo "  deps          - Download dependencies"
	@echo "  docs          - Generate documentation"

