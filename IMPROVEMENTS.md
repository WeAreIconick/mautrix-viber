# Potential Improvements & Future Enhancements

This document outlines areas for future improvement to make the codebase even more robust and production-ready.

## 🔴 High Priority (Should Consider)

### 1. Database Context Support
**Issue**: Most database methods don't accept `context.Context`, preventing cancellation and timeout control.

**Impact**: 
- Long-running queries can't be cancelled
- No timeout control for database operations
- Can't propagate request context for tracing

**Example**:
```go
// Current
func (d *DB) UpsertViberUser(viberID, viberName string) error

// Improved
func (d *DB) UpsertViberUser(ctx context.Context, viberID, viberName string) error
```

**Effort**: Medium (requires updating all database methods and call sites)

### 2. Hardcoded URLs
**Issue**: Viber API endpoint and OpenStreetMap URL are hardcoded.

**Current**:
- `https://chatapi.viber.com/pa/set_webhook` (in `internal/viber/client.go`)
- `https://www.openstreetmap.org/?mlat=...` (in `internal/viber/location.go`)

**Recommendation**: Make these configurable via environment variables or config struct.

**Effort**: Low

### 3. Context Propagation
**Issue**: Some places use `context.Background()` when context should be propagated.

**Found in**:
- `internal/viber/client.go` - `EnsureWebhook()` doesn't accept context
- `internal/matrix/client.go` - Message forwarding uses `context.Background()`
- `internal/matrix/events.go` - Event handlers create new contexts

**Recommendation**: Propagate context from HTTP requests through the call chain.

**Effort**: Medium

## 🟡 Medium Priority (Nice to Have)

### 4. Database Query Timeouts
**Issue**: No explicit query timeouts beyond connection-level settings.

**Recommendation**: Add context-based timeouts for all database operations.

**Effort**: Medium (part of #1)

### 5. Error Wrapping Consistency
**Issue**: Some errors use `fmt.Errorf` without `%w` verb, losing error chain.

**Recommendation**: Ensure all error wrapping uses `%w` to preserve error chains for `errors.Is()` and `errors.As()`.

**Effort**: Low (audit and fix)

### 6. HTTP Client Configuration
**Issue**: HTTP client timeout is hardcoded (15 seconds).

**Current**:
```go
httpClient: &http.Client{Timeout: 15 * time.Second}
```

**Recommendation**: Make timeout configurable via environment variable.

**Effort**: Low

### 7. Structured Logging Consistency
**Issue**: Some places have comments indicating structured logging should be used but don't actually log.

**Found in**:
- `internal/viber/client.go` - Error handling in webhook handler
- `internal/tracing/tracing.go` - Shutdown error handling

**Recommendation**: Replace all comment placeholders with actual structured logging calls.

**Effort**: Low

## 🟢 Low Priority (Future Enhancements)

### 8. Database Connection Retry
**Issue**: Database connection failures during initialization don't retry.

**Recommendation**: Add exponential backoff retry for database connection on startup.

**Effort**: Low

### 9. Configuration Validation Enhancements
**Issue**: Some validation could be more specific (e.g., Matrix room ID format validation).

**Recommendation**: Add regex-based validation for Matrix IDs, URLs, etc.

**Effort**: Low

### 10. Metrics Coverage
**Issue**: Some operations don't emit metrics.

**Recommendation**: Add metrics for:
- Database query duration
- Matrix API call duration
- Message processing latency

**Effort**: Medium

### 11. Test Coverage Gaps
**Issue**: Some functions have placeholder tests that are skipped.

**Recommendation**: Implement integration tests with mock servers.

**Effort**: High

### 12. Go Module Pinning
**Issue**: `go.mod` may not pin all dependency versions.

**Recommendation**: Ensure all dependencies are pinned to specific versions for reproducible builds.

**Effort**: Low

### 13. Documentation Examples
**Issue**: Some exported functions don't have usage examples in godoc.

**Recommendation**: Add godoc examples for commonly used functions.

**Effort**: Medium

## Code Quality Checks

### ✅ Already Good
- ✅ All SQL queries use parameterized queries (no SQL injection risk)
- ✅ Error handling is present throughout
- ✅ Structured logging is used (where implemented)
- ✅ Graceful shutdown is implemented
- ✅ Rate limiting is implemented
- ✅ Input validation exists
- ✅ Go vet passes with no issues
- ✅ All exported functions have godoc comments

### 🔍 Could Be Better
- Database methods should accept context
- Context propagation could be improved
- Some URLs should be configurable
- More comprehensive test coverage

## Immediate Action Items

If you want to prioritize improvements, I'd suggest:

1. **Make Viber API URL configurable** (Low effort, good practice)
2. **Add context to database methods** (Medium effort, significant improvement)
3. **Fix context propagation** (Medium effort, better observability)

Would you like me to implement any of these improvements?

