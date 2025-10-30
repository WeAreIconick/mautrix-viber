# Testing Guide

## Running Tests

### Prerequisites

1. **Download Dependencies**:
```bash
go mod tidy
```

This will download all required packages:
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/prometheus/client_golang` - Prometheus metrics
- `maunium.net/go/mautrix` - Matrix client library
- `gopkg.in/yaml.v3` - YAML parsing
- And other dependencies

### Running All Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Running Specific Test Packages

```bash
# Test database operations
go test ./internal/database/... -v

# Test configuration
go test ./internal/config/... -v

# Test retry logic
go test ./internal/retry/... -v

# Test circuit breaker
go test ./internal/circuitbreaker/... -v

# Test validation utilities
go test ./internal/utils/... -v

# Test signature verification
go test ./internal/viber/... -v -run TestSignature
```

### Test Coverage Goals

- **Core Components**: >80% coverage
- **Database Layer**: >90% coverage
- **Validation**: 100% coverage
- **Retry Logic**: >85% coverage

## Test Structure

### Unit Tests

Located in `*_test.go` files alongside source code:

- `internal/database/database_test.go` - Database operations
- `internal/config/config_test.go` - Configuration loading
- `internal/retry/retry_test.go` - Retry logic
- `internal/circuitbreaker/circuitbreaker_test.go` - Circuit breaker
- `internal/utils/validation_test.go` - Input validation
- `internal/viber/signature_test.go` - Signature verification

### Integration Tests

Integration tests should be added in `test/integration/` for:
- Full webhook processing flow
- Matrix ↔ Viber message bridging
- End-to-end scenarios

## Writing New Tests

### Example Test Structure

```go
func TestFeatureName(t *testing.T) {
    // Arrange - Set up test data
    testData := setupTestData(t)
    defer cleanup(testData)
    
    // Act - Execute the function
    result, err := functionUnderTest(testData)
    
    // Assert - Verify results
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if result != expectedValue {
        t.Errorf("Expected %v, got %v", expectedValue, result)
    }
}
```

### Test Helpers

Create test helpers for common setup:

```go
func setupTestDB(t *testing.T) *database.DB {
    dbPath := "/tmp/test_" + t.Name() + ".db"
    db, err := database.Open(dbPath)
    if err != nil {
        t.Fatalf("Failed to open test DB: %v", err)
    }
    t.Cleanup(func() {
        db.Close()
        os.Remove(dbPath)
    })
    return db
}
```

## Continuous Integration

Tests run automatically on:
- Push to main/develop branches
- Pull requests
- Via GitHub Actions (see `.github/workflows/ci.yml`)

## Known Issues

Some tests require external services:
- Matrix client tests require mock Matrix server
- Viber client tests require mock Viber API
- Redis cache tests require Redis instance

These can be mocked using test doubles or testcontainers.

