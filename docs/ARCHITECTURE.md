# Architecture Documentation

## Overview

mautrix-viber is a bidirectional bridge between Matrix and Viber messaging platforms. It provides real-time message synchronization, user management, and advanced features like media forwarding and encryption support.

## System Architecture

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Viber     │────────▶│   Bridge      │────────▶│   Matrix    │
│   Platform  │◀────────│   (Go)        │◀────────│  Homeserver │
└─────────────┘         └──────────────┘         └─────────────┘
                              │
                              ▼
                        ┌──────────────┐
                        │   SQLite DB  │
                        └──────────────┘
```

## Components

### Core Components

1. **Viber Client** (`internal/viber/`)
   - Webhook handler
   - Message sending API
   - Media forwarding
   - Event processing

2. **Matrix Client** (`internal/matrix/`)
   - Event listeners
   - Message sending
   - Room management
   - Ghost user puppeting

3. **Database Layer** (`internal/database/`)
   - User mappings
   - Room mappings
   - Message deduplication
   - Group membership

4. **Bridge Core** (`cmd/mautrix-viber/`)
   - HTTP server
   - Request routing
   - Middleware (rate limiting, logging)
   - Graceful shutdown

### Supporting Components

- **Configuration** (`internal/config/`) - Config management with hot reload
- **Metrics** (`internal/metrics/`) - Prometheus metrics
- **Tracing** (`internal/tracing/`) - OpenTelemetry tracing
- **Cache** (`internal/cache/`) - Redis caching layer
- **Retry Logic** (`internal/retry/`) - Exponential backoff
- **Circuit Breaker** (`internal/circuitbreaker/`) - Fault tolerance

## Message Flow

### Viber → Matrix

1. Viber sends webhook to `/viber/webhook`
2. Signature verification (HMAC-SHA256)
3. Parse webhook payload
4. Store sender in database
5. Forward message to Matrix room
6. Store message mapping for deduplication

### Matrix → Viber

1. Matrix client sync receives message event
2. Parse message content
3. Format for Viber (handle rich content)
4. Send via Viber API
5. Store message mapping

## Security

- **Signature Verification**: All Viber webhooks verified with HMAC-SHA256
- **Rate Limiting**: Per-IP, per-user, per-room adaptive limits
- **Input Validation**: All user input sanitized
- **HTTPS Required**: Production deployments must use HTTPS

## Scalability

- Stateless design (state in database)
- Horizontal scaling support
- Message queue for async processing
- Redis cache for hot data
- Connection pooling

## Reliability

- Circuit breaker pattern for external APIs
- Exponential backoff retry logic
- Message deduplication
- Health checks and monitoring
- Graceful degradation

## Performance

- Structured logging for observability
- Prometheus metrics for monitoring
- OpenTelemetry tracing for debugging
- Database indexing for fast queries
- Caching layer for frequently accessed data

## Deployment

See [DEPLOYMENT.md](./DEPLOYMENT.md) for deployment instructions.

## Extension Points

The bridge is designed for extensibility:

- Plugin system for custom message handlers
- Webhook system for external integrations
- Admin API for management operations
- Configurable rate limits and timeouts

