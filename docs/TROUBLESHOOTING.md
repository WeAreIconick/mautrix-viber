# Troubleshooting Guide

Common issues and solutions for mautrix-viber bridge.

## Webhook Issues

### Webhook Not Receiving Events

**Symptoms:**
- No messages appearing in Matrix
- Webhook endpoint not logging requests

**Diagnosis:**
```bash
# Check if webhook is registered
curl -X GET "https://chatapi.viber.com/pa/get_webhook" \
  -H "X-Viber-Auth-Token: YOUR_TOKEN"

# Check webhook endpoint accessibility
curl -X POST https://your-domain.com/viber/webhook \
  -H "Content-Type: application/json" \
  -d '{"test":"data"}'
```

**Solutions:**
1. Verify webhook URL is publicly accessible (use `ngrok` for testing)
2. Ensure webhook URL uses HTTPS (Viber requires HTTPS in production)
3. Check firewall/security group rules allow incoming connections
4. Verify webhook registration on startup (check logs)
5. Test with a simple webhook receiver to verify Viber can reach your server

### Invalid Signature Errors

**Symptoms:**
- 401 Unauthorized errors in logs
- Webhook requests rejected

**Solutions:**
1. Verify `VIBER_API_TOKEN` matches the token used to register webhook
2. Check that signature header is `X-Viber-Content-Signature`
3. Ensure request body is not modified (by middleware, proxy, etc.)
4. Verify signature calculation matches Viber's algorithm (HMAC-SHA256)

## Database Issues

### Database Locked Errors

**Symptoms:**
- "database is locked" errors
- Slow database operations

**Solutions:**
1. Check for concurrent write operations
2. Verify connection pool settings are appropriate
3. Ensure WAL mode is enabled (default)
4. Check disk space and I/O performance
5. Review database busy timeout settings

### Migration Failures

**Symptoms:**
- Startup fails with migration error
- Tables missing

**Solutions:**
1. Check database file permissions
2. Verify SQLite version supports all features
3. Review migration logs for specific errors
4. Manually run migration SQL if needed

## Matrix Connection Issues

### Cannot Connect to Matrix Homeserver

**Symptoms:**
- "failed to initialize matrix client" errors
- Matrix events not received

**Diagnosis:**
```bash
# Test Matrix connection
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://matrix.example.com/_matrix/client/r0/account/whoami
```

**Solutions:**
1. Verify `MATRIX_HOMESERVER_URL` is correct (no trailing slash)
2. Check `MATRIX_ACCESS_TOKEN` is valid and not expired
3. Verify homeserver is accessible from bridge location
4. Check firewall rules allow outbound HTTPS
5. Review Matrix server logs for connection issues

### Messages Not Bridging to Matrix

**Symptoms:**
- Viber messages received but not appearing in Matrix
- No errors in logs

**Solutions:**
1. Verify `MATRIX_DEFAULT_ROOM_ID` is correct format (`!roomid:server.com`)
2. Check bridge has permission to send messages to room
3. Verify Matrix client initialized successfully (check startup logs)
4. Review message forwarding logic in logs (debug mode)
5. Check for rate limiting on Matrix server

## Performance Issues

### High Memory Usage

**Symptoms:**
- Memory usage growing over time
- OOM (Out of Memory) errors

**Diagnosis:**
```bash
# Check memory usage
go tool pprof http://localhost:8080/debug/pprof/heap
```

**Solutions:**
1. Review goroutine leaks (use pprof)
2. Check for unbounded slice/map growth
3. Verify database connections are properly closed
4. Review connection pool settings
5. Check for memory leaks in long-running operations

### Slow Response Times

**Symptoms:**
- Webhook responses taking >5 seconds
- High latency in message bridging

**Solutions:**
1. Review database query performance (use EXPLAIN QUERY PLAN)
2. Check for N+1 query problems
3. Verify indexes are being used
4. Review external API call timeouts
5. Check network latency to Viber/Matrix servers
6. Monitor CPU and memory usage

## Configuration Issues

### Configuration Validation Failures

**Symptoms:**
- "configuration validation failed" on startup
- Missing required configuration errors

**Solutions:**
1. Review validation error message for specific issues
2. Verify all required environment variables are set
3. Check URL formats are correct
4. Ensure HTTPS URLs in production
5. Validate Matrix room ID format

### Environment Variables Not Loading

**Symptoms:**
- Default values used instead of env vars
- Configuration not taking effect

**Solutions:**
1. Verify environment variables are exported
2. Check for typos in variable names
3. Ensure `.env` file is loaded if using one
4. Restart service after changing environment variables
5. Use `printenv` to verify variables are set

## Security Issues

### Signature Verification Failures

**Symptoms:**
- All webhook requests rejected
- Constant 401 errors

**Solutions:**
1. Verify API token matches between registration and verification
2. Check proxy/load balancer isn't modifying request body
3. Review middleware order (signature should verify before body parsing)
4. Test signature calculation manually
5. Check for encoding issues (UTF-8 vs other)

### Rate Limiting Too Aggressive

**Symptoms:**
- Legitimate requests being rate limited
- 429 errors for valid users

**Solutions:**
1. Review rate limit settings in code
2. Adjust token bucket parameters if needed
3. Consider per-user rate limits instead of per-IP
4. Whitelist trusted IPs if appropriate
5. Review rate limit logs for patterns

## Monitoring & Debugging

### Enable Debug Logging

```bash
export LOG_LEVEL=debug
./bin/mautrix-viber
```

### Check Metrics

```bash
# View Prometheus metrics
curl http://localhost:8080/metrics

# Check specific metrics
curl http://localhost:8080/metrics | grep viber_webhook_requests
```

### Health Check Status

```bash
# Basic health
curl http://localhost:8080/healthz

# Readiness (includes dependency checks)
curl http://localhost:8080/readyz
```

### View Recent Logs

```bash
# If using structured logging
tail -f logs/bridge.log | jq

# Search for errors
grep -i error logs/bridge.log
```

## Common Error Messages

### "database is locked"
- **Cause**: Concurrent access or long-running transaction
- **Solution**: Review connection pool settings, ensure proper transaction handling

### "invalid signature"
- **Cause**: Token mismatch or request body modification
- **Solution**: Verify API token, check middleware order

### "context deadline exceeded"
- **Cause**: Operation timeout
- **Solution**: Increase timeout or optimize slow operations

### "rate limit exceeded"
- **Cause**: Too many requests from same IP
- **Solution**: Adjust rate limits or implement per-user limits

### "connection refused"
- **Cause**: Matrix/Viber server unreachable
- **Solution**: Check network connectivity, verify URLs

## Getting Help

1. Check logs with `LOG_LEVEL=debug`
2. Review metrics at `/metrics`
3. Test health checks at `/healthz` and `/readyz`
4. Review configuration with validation
5. Check GitHub Issues for known problems
6. Enable tracing for distributed debugging

## Diagnostic Commands

```bash
# Full diagnostic check
./scripts/health-check.sh

# Database integrity check
sqlite3 data/bridge.db "PRAGMA integrity_check;"

# Check database schema
sqlite3 data/bridge.db ".schema"

# View recent webhook activity
sqlite3 data/bridge.db "SELECT * FROM message_mappings ORDER BY created_at DESC LIMIT 10;"
```

