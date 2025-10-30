# Development Guide

Complete guide for developing and contributing to mautrix-viber.

## Development Setup

### Prerequisites

- Go 1.22 or later
- Git
- Make (optional, for Makefile targets)
- Docker (for testing deployments)

### Initial Setup

```bash
# Clone repository
git clone https://github.com/example/mautrix-viber.git
cd mautrix-viber

# Download dependencies
go mod download
go mod tidy

# Verify setup
go build ./cmd/mautrix-viber
go test ./...
```

## Project Structure

```
mautrix-viber/
├── cmd/mautrix-viber/    # Main application entry point
│   ├── main.go           # Server setup and routing
│   ├── health.go         # Health check handlers
│   └── ratelimit*.go    # Rate limiting implementations
├── internal/
│   ├── admin/            # Matrix admin commands
│   ├── api/              # REST API endpoints
│   ├── cache/            # Redis caching
│   ├── circuitbreaker/   # Circuit breaker pattern
│   ├── config/           # Configuration management
│   ├── database/         # Database layer
│   ├── logger/           # Structured logging
│   ├── matrix/           # Matrix client wrapper
│   ├── metrics/          # Prometheus metrics
│   ├── middleware/       # HTTP middleware
│   ├── queue/            # Message queue
│   ├── retry/            # Retry logic
│   ├── tracing/         # OpenTelemetry tracing
│   ├── utils/            # Utility functions
│   ├── version/          # Version management
│   └── viber/            # Viber client and handlers
├── test/                 # Integration tests
├── docs/                 # Documentation
├── scripts/              # Utility scripts
└── k8s/                  # Kubernetes manifests
```

## Development Workflow

### Making Changes

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make Changes**
   - Follow coding standards (see `.cursorrules`)
   - Write tests for new functionality
   - Update documentation

3. **Test Changes**
   ```bash
   # Run tests
   go test ./...
   
   # Run linters
   golangci-lint run
   
   # Build
   go build ./cmd/mautrix-viber
   ```

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

5. **Push and Create PR**
   ```bash
   git push origin feature/my-feature
   # Create pull request on GitHub
   ```

## Testing

### Unit Tests

```bash
# Run all unit tests
go test ./... -v

# Run specific package
go test ./internal/database/... -v

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Run integration tests (requires setup)
go test ./test/integration/... -v
```

### Benchmarks

```bash
go test -bench=. -benchmem ./test/benchmark/...
```

## Code Style

### Formatting

```bash
# Format code
go fmt ./...

# Or use goimports
goimports -w .
```

### Linting

```bash
# Run golangci-lint
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix
```

### Code Review Checklist

- [ ] Code follows Go conventions
- [ ] All exported functions documented
- [ ] Tests added for new functionality
- [ ] No hardcoded secrets
- [ ] Error handling appropriate
- [ ] Context used for cancellation
- [ ] No race conditions
- [ ] Performance considerations

## Adding New Features

### 1. Database Changes

If adding new tables or columns:

```go
// 1. Add migration in database.go migrate() function
// 2. Add corresponding methods
// 3. Add tests in database_test.go
// 4. Update schema documentation
```

### 2. New HTTP Endpoints

```go
// 1. Add handler function
func (h *Handler) HandleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// 2. Register in main.go
mux.HandleFunc("/api/new", handler.HandleNewEndpoint)

// 3. Add tests
func TestHandleNewEndpoint(t *testing.T) {
    // Test implementation
}
```

### 3. New Viber Features

```go
// 1. Add handler in internal/viber/
// 2. Integrate into WebhookHandler
// 3. Add forwarding logic if needed
// 4. Add tests
```

## Debugging

### Debug Logging

```bash
export LOG_LEVEL=debug
./bin/mautrix-viber
```

### Profiling

```bash
# Add to main.go:
import _ "net/http/pprof"

# Then access:
# http://localhost:8080/debug/pprof/
```

### Database Inspection

```bash
# Open SQLite database
sqlite3 data/bridge.db

# Check schema
.schema

# Query data
SELECT * FROM viber_users LIMIT 10;
```

## Common Tasks

### Adding Dependencies

```bash
# Add dependency
go get github.com/new/package

# Update go.mod and go.sum
go mod tidy
```

### Updating Dependencies

```bash
# Update all dependencies
go get -u ./...

# Update specific dependency
go get -u github.com/package/name

# Verify
go mod verify
```

### Running Locally

```bash
# Set environment variables
export VIBER_API_TOKEN="..."
export VIBER_WEBHOOK_URL="..."

# Run
go run ./cmd/mautrix-viber

# Or build and run
go build -o bin/mautrix-viber ./cmd/mautrix-viber
./bin/mautrix-viber
```

## CI/CD

### GitHub Actions

The project includes GitHub Actions workflows:
- `.github/workflows/ci.yml` - Runs on PR and push

### Local CI Checks

```bash
# Run all CI checks locally
make lint
make test
make build
```

## Documentation

### Updating Documentation

- README.md - Main project documentation
- docs/ - Detailed guides
- Code comments - Inline documentation

### Documentation Standards

- All exported functions must have comments
- Complex logic should be explained
- Examples should be provided
- Keep documentation up to date

## Getting Help

- Review existing code for patterns
- Check `.cursorrules` for standards
- Review test files for examples
- Ask questions in GitHub Issues or Discussions

