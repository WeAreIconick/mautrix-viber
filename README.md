# mautrix-viber — Production-Ready Matrix ↔ Viber Bridge (Go)

A comprehensive, production-ready bidirectional Matrix-Viber bridge written in Go. Built with security, observability, and reliability in mind.

## Features

### ✅ Implemented

- **Bidirectional Message Bridging**
  - Viber → Matrix: Text, images, and media messages
  - Matrix → Viber: Full message forwarding with rich formatting
  
- **Security & Reliability**
  - HMAC-SHA256 signature verification for Viber webhooks
  - Per-IP rate limiting with token bucket algorithm
  - Request body size limits (2MB default)
  - Graceful shutdown with 15s timeout
  - Server timeouts (read, write, idle)

- **Observability**
  - Prometheus metrics (`/metrics`)
  - Structured JSON logging via `log/slog`
  - Health check endpoints (`/healthz`, `/readyz`)
  - Bridge info endpoint (`/api/info`) with status and statistics

- **Persistence**
  - SQLite database for user/room mappings
  - Message ID deduplication
  - User linking and room mapping storage

- **Configuration**
  - Environment variable support
  - YAML configuration file support
  - Configuration validation
  - Config overrides (env vars override file config)

- **Deployment**
  - Dockerfile with multi-stage build
  - docker-compose.yml with healthchecks
  - Alpine-based minimal image

- **API Features**
  - Viber send API: text, image, video, file, location, contact messages
  - Matrix event listeners: message, reaction, redaction, typing, receipt
  - Admin commands: `!bridge link`, `!bridge unlink`, `!bridge status`, `!bridge help`, `!bridge ping`

- **Developer Experience**
  - Exponential backoff retry logic with jitter
  - Comprehensive error handling
  - Well-documented codebase with inline comments
  - Clean package structure

### 🚧 Planned Features

See [TODO list](#todo) for comprehensive feature roadmap including:
- Matrix ghost user puppeting with avatars
- Group chat support
- E2EE support
- Typing indicators and read receipts sync
- Message search
- Web admin panel
- OpenTelemetry tracing
- Redis caching
- And 40+ more features...

---

## Quick Start

### Prerequisites

- Go 1.22+ (for building from source)
- A Viber Bot API token ([create one](https://partners.viber.com/))
- A Matrix homeserver URL and access token
- A publicly accessible HTTPS URL for webhooks

### Option 1: Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/example/mautrix-viber.git
cd mautrix-viber

# Create config file
cp config.example.yaml config.yaml
# Edit config.yaml with your credentials

# Or use environment variables
export VIBER_API_TOKEN="your-token"
export VIBER_WEBHOOK_URL="https://your-domain.com/viber/webhook"
export MATRIX_HOMESERVER_URL="https://matrix.example.com"
export MATRIX_ACCESS_TOKEN="your-token"
export MATRIX_DEFAULT_ROOM_ID="!roomid:example.com"

# Run with docker-compose
docker-compose up -d

# Or build and run manually
docker build -t mautrix-viber .
docker run -d \
  -p 8080:8080 \
  -e VIBER_API_TOKEN="your-token" \
  -e VIBER_WEBHOOK_URL="https://your-domain.com/viber/webhook" \
  -v ./data:/data \
  mautrix-viber
```

### Option 2: From Source

```bash
# Clone and build
git clone https://github.com/example/mautrix-viber.git
cd mautrix-viber
go build -o ./bin/mautrix-viber ./cmd/mautrix-viber

# Set environment variables (see Configuration section)
export VIBER_API_TOKEN="your-token"
export VIBER_WEBHOOK_URL="https://your-domain.com/viber/webhook"
export MATRIX_HOMESERVER_URL="https://matrix.example.com"
export MATRIX_ACCESS_TOKEN="your-token"
export MATRIX_DEFAULT_ROOM_ID="!roomid:example.com"

# Run
./bin/mautrix-viber
```

### Testing Locally

For local development, use `ngrok` to expose the bridge:

```bash
# In one terminal
./bin/mautrix-viber

# In another terminal
ngrok http 8080

# Copy the HTTPS URL (e.g., https://abcd1234.ngrok.io)
# Set VIBER_WEBHOOK_URL=https://abcd1234.ngrok.io/viber/webhook
```

---

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `VIBER_API_TOKEN` | Viber Bot API token | ✅ Yes |
| `VIBER_WEBHOOK_URL` | Public HTTPS URL for Viber webhooks | ✅ Yes |
| `LISTEN_ADDRESS` | HTTP server listen address (default: `:8080`) | No |
| `MATRIX_HOMESERVER_URL` | Matrix homeserver base URL | Yes (if bridging) |
| `MATRIX_ACCESS_TOKEN` | Matrix access token | Yes (if bridging) |
| `MATRIX_DEFAULT_ROOM_ID` | Default Matrix room for bridged messages | Yes (if bridging) |
| `DATABASE_PATH` | SQLite database path (default: `./data/bridge.db`) | No |
| `LOG_LEVEL` | Log level: debug, info, warn, error (default: `info`) | No |
| `VIBER_DEFAULT_RECEIVER_ID` | Default Viber user ID for Matrix → Viber demo forwarding | Optional |

### YAML Configuration File

Create `config.yaml` (see `config.example.yaml` for template):

```yaml
viber:
  api_token: "your-viber-bot-token"
  webhook_url: "https://your-domain.com/viber/webhook"

matrix:
  homeserver_url: "https://matrix.example.com"
  access_token: "your-matrix-access-token"
  default_room_id: "!roomid:example.com"

server:
  listen_address: ":8080"

database:
  path: "./data/bridge.db"

logging:
  level: "info"  # debug, info, warn, error
```

**Note**: Environment variables override file configuration values.

---

## API Endpoints

### Webhook Endpoint

- **POST** `/viber/webhook` — Receives Viber callbacks
  - Verifies HMAC-SHA256 signature (`X-Viber-Content-Signature` header)
  - Processes events: `message`, `subscribed`, `unsubscribed`, `conversation_started`
  - Forwards messages to Matrix when configured

### Health & Monitoring

- **GET** `/healthz` — Health check (returns 200 if healthy)
- **GET** `/readyz` — Readiness check (returns 200 if ready)
- **GET** `/metrics` — Prometheus metrics
- **GET** `/api/info` — Bridge information and statistics (JSON)

### Example: Get Bridge Status

```bash
curl http://localhost:8080/api/info
```

Response:
```json
{
  "version": "0.1.0",
  "status": "running",
  "uptime": "2h30m15s",
  "started_at": "2024-01-01T00:00:00Z",
  "matrix": {
    "connected": true,
    "status": "synced"
  },
  "viber": {
    "connected": true,
    "status": "webhook_registered"
  },
  "statistics": {
    "messages_bridged": 1234,
    "users_linked": 56,
    "rooms_mapped": 12,
    "webhook_requests": 5678,
    "errors": 0
  }
}
```

---

## Admin Commands

Bridge commands can be run in Matrix rooms:

- `!bridge help` — Show available commands
- `!bridge link <viber-user-id>` — Link a Viber user to your Matrix account
- `!bridge unlink` — Unlink your Viber account
- `!bridge status` — Show bridge status and statistics
- `!bridge ping` — Test bridge responsiveness

---

## Security

### Webhook Security

- **Signature Verification**: All Viber webhooks are verified using HMAC-SHA256
- **Rate Limiting**: Per-IP token bucket rate limiter (5 req/sec, burst 10)
- **Body Size Limits**: Maximum 2MB request body size

### Best Practices

1. **Always use HTTPS** for webhook URLs in production
2. **Keep tokens secret**: Never commit tokens to version control
3. **Monitor metrics**: Watch `/metrics` for unusual patterns
4. **Review logs**: Structured logs help identify security issues
5. **Update regularly**: Keep dependencies updated for security patches

---

## Metrics

The bridge exposes Prometheus metrics at `/metrics`:

- `viber_webhook_requests_total` — Total webhook requests by event type
- `viber_messages_forwarded_total` — Messages forwarded to Matrix by type
- `viber_signature_failures_total` — Signature verification failures

---

## Development

### Project Structure

```
cmd/
  mautrix-viber/
    main.go              # Application entry point
    ratelimit.go         # Rate limiting middleware
internal/
  admin/
    commands.go          # Bridge admin commands
  api/
    info.go              # API endpoints (info, health)
  config/
    config.go            # Environment config loader
    config_file.go        # YAML config loader with validation
  database/
    database.go          # SQLite persistence layer
  logger/
    logger.go            # Structured JSON logging
  matrix/
    client.go            # Matrix client wrapper
    events.go            # Matrix event listeners
  retry/
    retry.go             # Exponential backoff retry logic
  viber/
    client.go            # Viber webhook handler
    send.go              # Viber send API functions
    types.go             # Viber API types
    metrics.go           # Prometheus metrics
go.mod                   # Go dependencies
Dockerfile               # Docker image definition
docker-compose.yml       # Docker Compose setup
config.example.yaml      # Example configuration
```

### Building

```bash
# Build binary
go build -o ./bin/mautrix-viber ./cmd/mautrix-viber

# Run tests (when available)
go test ./...

# Format code
go fmt ./...

# Lint
go vet ./...
```

### Adding Features

The codebase is structured for easy extension:

1. **Database operations**: Add methods to `internal/database/database.go`
2. **Viber API calls**: Add functions to `internal/viber/send.go`
3. **Matrix operations**: Extend `internal/matrix/client.go`
4. **Admin commands**: Register in `internal/admin/commands.go`

---

## Troubleshooting

### Webhook not receiving events

1. Check webhook URL is publicly accessible
2. Verify signature verification (check logs for failures)
3. Ensure webhook is registered: check logs on startup
4. Test with `curl -X POST https://your-url/viber/webhook -H "X-Viber-Content-Signature: test"`

### Messages not bridging

1. Verify Matrix credentials are correct
2. Check Matrix room ID format (`!roomid:example.com`)
3. Review logs for errors
4. Check `/api/info` for connection status

### High error rates

1. Check Prometheus metrics at `/metrics`
2. Review structured logs for patterns
3. Verify API rate limits aren't being exceeded
4. Check database connection and disk space

---

## License

This project is licensed under the MIT License. See `LICENSE` for details.

---

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new features
4. Submit a pull request

---

## TODO

See the comprehensive [feature roadmap](TODO.md) for planned enhancements:

- [ ] Matrix ghost user puppeting
- [ ] Group chat support
- [ ] E2EE support  
- [ ] Typing indicators & read receipts
- [ ] Message search
- [ ] Web admin panel
- [ ] OpenTelemetry tracing
- [ ] Redis caching
- [ ] Circuit breaker pattern
- [ ] Hot config reload
- [ ] And 40+ more features...

---

## Acknowledgments

Built with:
- [mautrix-go](https://github.com/mautrix/go) — Matrix client library
- [Prometheus](https://prometheus.io/) — Metrics
- [log/slog](https://pkg.go.dev/log/slog) — Structured logging

---

**Status**: Production-ready with comprehensive feature set. Actively maintained and extended.
