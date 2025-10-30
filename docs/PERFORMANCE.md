# Performance Tuning Guide

Guide for optimizing mautrix-viber bridge performance.

## Database Performance

### Connection Pooling

The bridge uses connection pooling with these defaults:
- MaxOpenConns: 25
- MaxIdleConns: 5
- ConnMaxLifetime: 5 minutes
- ConnMaxIdleTime: 10 minutes

**Tuning:**
- Increase `MaxOpenConns` for high concurrency
- Adjust `MaxIdleConns` based on average load
- Set `ConnMaxLifetime` based on database stability requirements

### Query Optimization

**Use Indexes:**
```sql
-- Indexes are automatically created, but verify:
.schema

-- If needed, add custom indexes:
CREATE INDEX idx_custom ON table(column);
```

**Prepared Statements:**
All queries use prepared statements for performance.

**Transaction Batching:**
Group multiple operations into transactions when possible.

### SQLite-Specific Optimizations

```go
// WAL mode (already enabled)
?_journal_mode=WAL

// Busy timeout (already configured)
?_busy_timeout=5000

// Additional options:
?_synchronous=NORMAL    // Faster writes, slight risk
?_cache_size=10000       // Increase cache for larger DBs
```

## Memory Optimization

### Goroutine Management

Monitor goroutine count:
```bash
curl http://localhost:8080/debug/pprof/goroutine
```

**Best Practices:**
- Use context cancellation for all background operations
- Limit goroutine pool sizes
- Monitor for goroutine leaks

### Memory Profiling

```bash
# Enable pprof endpoint (add to main.go)
go tool pprof http://localhost:8080/debug/pprof/heap

# Check memory usage
go tool pprof http://localhost:8080/debug/pprof/heap
(pprof) top
(pprof) list functionName
```

### Resource Limits

Set appropriate limits:
```yaml
# Kubernetes
resources:
  limits:
    memory: "512Mi"
    cpu: "500m"
  requests:
    memory: "128Mi"
    cpu: "100m"
```

## Network Performance

### Timeout Configuration

Optimize timeouts for your network:
```go
// HTTP Client timeouts
httpClient := &http.Client{
    Timeout: 15 * time.Second, // Adjust based on network latency
}

// Server timeouts
ReadTimeout:       10 * time.Second,
WriteTimeout:      15 * time.Second,
IdleTimeout:       60 * time.Second,
```

### Connection Reuse

HTTP clients reuse connections automatically via connection pooling.

### Keep-Alive

Ensure keep-alive is enabled (default in Go HTTP clients).

## Rate Limiting

### Adaptive Rate Limiting

The bridge supports adaptive rate limiting that adjusts based on error rates.

**Configuration:**
- Base rate: 5 req/sec
- Burst: 10 requests
- Adjustment: ±10-20% based on error rate

**Tuning:**
- Increase rates for high-capacity deployments
- Decrease rates if hitting upstream limits
- Monitor error rates in metrics

## Caching

### Redis Caching (Optional)

Enable Redis for frequently accessed data:
```bash
export REDIS_URL="redis://localhost:6379"
```

**Cache TTL:**
- User mappings: 5 minutes
- Room mappings: 10 minutes
- Message mappings: 1 hour

## Message Processing

### Async Processing

Use message queue for non-blocking processing:
- Enables parallel message handling
- Provides retry logic
- Prevents blocking webhook responses

### Batch Operations

Batch database writes when possible:
```go
// Instead of:
for _, msg := range messages {
    db.StoreMessageMapping(...)
}

// Use:
tx, _ := db.Begin()
for _, msg := range messages {
    tx.StoreMessageMapping(...)
}
tx.Commit()
```

## Monitoring Performance

### Key Metrics

Monitor these Prometheus metrics:
- `viber_message_latency_seconds` - Message processing time
- `viber_webhook_requests_total` - Request volume
- `database_connections_active` - Connection pool usage
- `goroutine_count` - Goroutine leaks

### Alerting Thresholds

Set alerts for:
- Message latency > 2 seconds
- Error rate > 5%
- Connection pool usage > 80%
- Memory usage > 80%
- CPU usage > 80% sustained

## Benchmarking

### Run Benchmarks

```bash
# Benchmark database operations
go test -bench=BenchmarkDatabase -benchmem ./test/benchmark/...

# Benchmark HTTP handlers
go test -bench=BenchmarkHandler -benchmem ./test/benchmark/...

# Compare before/after changes
go test -bench=. -benchmem -count=5 ./...
```

### Performance Targets

- Database query: < 10ms for indexed lookups
- Webhook processing: < 100ms end-to-end
- Message forwarding: < 500ms to Matrix
- Memory usage: < 256MB under normal load
- CPU usage: < 50% under normal load

## Optimization Checklist

- [ ] Database indexes verified
- [ ] Connection pool sized appropriately
- [ ] Timeouts configured for network conditions
- [ ] Rate limits tuned for expected load
- [ ] Caching enabled for hot data
- [ ] Async processing enabled
- [ ] Goroutine leaks monitored
- [ ] Memory profiling done
- [ ] Benchmarks run regularly
- [ ] Performance metrics monitored

## Troubleshooting Performance

1. **Slow Database Queries**
   - Run EXPLAIN QUERY PLAN
   - Verify indexes are used
   - Check for table scans

2. **High Memory Usage**
   - Profile with pprof
   - Check for unbounded growth
   - Review connection pool settings

3. **Slow Message Processing**
   - Check external API response times
   - Review retry logic impact
   - Verify async processing is working

4. **High CPU Usage**
   - Profile CPU usage
   - Check for tight loops
   - Review cryptographic operations

