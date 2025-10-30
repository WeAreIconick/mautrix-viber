# Production Hardening Checklist

This document tracks all production hardening improvements made to ensure the bridge is rock solid.

## ✅ Implemented Hardening Features

### 1. Panic Recovery
- **Location**: `internal/middleware/recovery.go`
- **What**: Global panic recovery middleware that catches panics, logs them with stack traces, and returns 500 errors instead of crashing the server
- **Benefit**: Prevents server crashes from unexpected panics, improves availability

### 2. Request ID Tracking
- **Location**: `internal/middleware/request_id.go`
- **What**: Adds unique request ID to every request for distributed tracing
- **Benefit**: Enables tracking requests across services, improves debugging

### 3. Configuration Validation
- **Location**: `internal/config/validation.go`
- **What**: Validates all configuration on startup before the server starts
- **Benefit**: Fails fast with clear error messages if configuration is invalid

### 4. Database Connection Pooling
- **Location**: `internal/database/database.go` (Open function)
- **What**: Configures connection pool limits and timeouts:
  - MaxOpenConns: 25 connections
  - MaxIdleConns: 5 connections
  - ConnMaxLifetime: 5 minutes
  - ConnMaxIdleTime: 10 minutes
  - BusyTimeout: 5 seconds
- **Benefit**: Prevents connection exhaustion, improves stability under load

### 5. Enhanced Health Checks
- **Location**: `cmd/mautrix-viber/health.go`
- **What**: 
  - `/healthz`: Basic health check (always returns 200 if process is running)
  - `/readyz`: Readiness check that verifies database connectivity with timeout
- **Benefit**: Kubernetes/orchestration can properly route traffic based on actual readiness

### 6. Middleware Chain
- **Location**: `cmd/mautrix-viber/main.go`
- **What**: Proper middleware ordering:
  1. Recovery (catch panics)
  2. Request ID (add tracking)
  3. Logging (structured request logs)
  4. Rate limiting (prevent abuse)
  5. Body size limits (prevent DoS)
- **Benefit**: Defense in depth, proper request handling pipeline

### 7. Server Configuration
- **Location**: `cmd/mautrix-viber/main.go` (http.Server)
- **What**: Comprehensive timeouts and limits:
  - ReadTimeout: 10 seconds
  - ReadHeaderTimeout: 5 seconds
  - WriteTimeout: 15 seconds
  - IdleTimeout: 60 seconds
  - MaxHeaderBytes: 1MB
- **Benefit**: Prevents resource exhaustion from slow clients

## 🔄 Already Implemented (From Before)

- ✅ Graceful shutdown (15s timeout)
- ✅ Rate limiting (per-IP token bucket)
- ✅ Request body size limits (2MB)
- ✅ HMAC-SHA256 signature verification
- ✅ Structured logging (log/slog)
- ✅ Prometheus metrics
- ✅ Circuit breaker pattern
- ✅ Retry logic with exponential backoff
- ✅ Input validation and sanitization

## 📋 Additional Recommendations

### High Priority

1. **Goroutine Leak Prevention**
   - Add context cancellation to all background goroutines
   - Use sync.WaitGroup for cleanup
   - Consider adding goroutine tracking metrics

2. **Memory Limits**
   - Add memory profiling endpoint (`/debug/pprof/heap`)
   - Set GOMEMLIMIT if using Go 1.19+
   - Monitor memory usage in metrics

3. **Context Propagation**
   - Ensure all external calls accept and respect context
   - Add request timeout middleware
   - Use context.WithTimeout for all I/O operations

4. **Error Rate Monitoring**
   - Track error rates by type
   - Alert on high error rates
   - Circuit breaker integration with metrics

### Medium Priority

5. **Database Query Timeouts**
   - Ensure all database queries use context with timeout
   - Add query timeout configuration
   - Track slow queries

6. **External API Timeouts**
   - Verify all HTTP clients have timeouts
   - Use circuit breaker for external APIs
   - Track API response times

7. **Request Validation**
   - Validate all webhook payloads
   - Check message sizes before processing
   - Reject malformed requests early

8. **Resource Limits**
   - Monitor goroutine count
   - Track connection pool usage
   - Alert on resource exhaustion

### Low Priority

9. **Performance Testing**
   - Load testing scripts
   - Benchmark critical paths
   - Performance regression tests

10. **Security Audit**
    - Dependency vulnerability scanning
    - Security headers (CSP, HSTS, etc.)
    - Penetration testing

## Testing Hardening

Run these to verify hardening:

```bash
# Test panic recovery
curl -X POST http://localhost:8080/viber/webhook -d '{"invalid":"json"}'

# Test rate limiting
for i in {1..20}; do curl http://localhost:8080/healthz; done

# Test configuration validation
VIBER_API_TOKEN="" ./bin/mautrix-viber

# Test health checks
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

## Monitoring

Monitor these metrics for production health:

- Panic recovery count (should be near zero)
- Request ID coverage (should be 100%)
- Database connection pool usage
- Health check response times
- Rate limit rejections
- Error rates by type

## Summary

The bridge now has:
- ✅ Panic recovery to prevent crashes
- ✅ Request tracking for debugging
- ✅ Configuration validation on startup
- ✅ Database connection pooling
- ✅ Enhanced health checks
- ✅ Comprehensive middleware chain
- ✅ Proper server timeouts

**Status**: Production-hardened with defense in depth principles applied.

