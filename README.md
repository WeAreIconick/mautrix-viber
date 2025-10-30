# mautrix-viber — Production-Ready Matrix ↔ Viber Bridge (Go)

A comprehensive, production-ready bidirectional Matrix-Viber bridge written in Go. Built with security, observability, and reliability in mind.

## Features

### ✅ Fully Implemented (49+ Features)

#### Core Bridging Features
- ✅ **Bidirectional Message Bridging**
  - Viber → Matrix: Text, images, video, audio, files, stickers, locations, contacts
  - Matrix → Viber: Full message forwarding with rich formatting
- ✅ **Media Support**: All media types (images, video, audio, files, stickers)
- ✅ **Rich Formatting**: Replies, threads, reactions, markdown parsing, mentions
- ✅ **Ghost User Puppeting**: Matrix ghost users for Viber contacts with avatars
- ✅ **Portal Rooms**: Auto-create Matrix rooms for Viber chats with metadata sync
- ✅ **Group Chat Support**: Viber group chats mapped to Matrix rooms with member sync
- ✅ **Typing Indicators & Read Receipts**: Bidirectional synchronization
- ✅ **Message Edits & Deletions**: Viber deletions → Matrix redactions
- ✅ **History Backfill**: Recent Viber message history on room creation
- ✅ **Message Search**: Bridge message search capabilities
- ✅ **E2EE Support**: Matrix encrypted room creation and message handling
- ✅ **Presence Sync**: User online/offline status synchronization
- ✅ **Power Levels**: Admin/moderator permission sync between platforms
- ✅ **Room Metadata**: Sync names, topics, and avatars between platforms
- ✅ **Notifications**: Configure Matrix push rules based on Viber settings

#### Infrastructure & Reliability
- ✅ **SQLite Database**: User/room mappings, message deduplication, migrations
- ✅ **Redis Caching**: Frequently accessed user/room mappings
- ✅ **Message Queue**: Reliable delivery with retry logic
- ✅ **Circuit Breaker**: Fault tolerance for external API calls
- ✅ **Advanced Rate Limiting**: Per-user, per-room, adaptive limits
- ✅ **Exponential Backoff**: Retry logic with jitter
- ✅ **Structured Logging**: JSON logging via `log/slog` with levels
- ✅ **Prometheus Metrics**: Comprehensive metrics at `/metrics`
- ✅ **OpenTelemetry Tracing**: Request flow tracing with Jaeger support
- ✅ **Hot Config Reload**: SIGHUP-based configuration reload without restart
- ✅ **Graceful Shutdown**: 15s timeout with cleanup
- ✅ **Health Checks**: `/healthz`, `/readyz` endpoints

#### Security
- ✅ **HMAC-SHA256 Verification**: Webhook signature verification
- ✅ **Per-IP Rate Limiting**: Token bucket algorithm (5 req/sec, burst 10)
- ✅ **Request Body Limits**: 2MB default maximum
- ✅ **Input Validation**: Comprehensive sanitization
- ✅ **HTTPS Enforcement**: Production security requirements
- ✅ **Panic Recovery**: Server crash prevention
- ✅ **Request ID Tracking**: Distributed tracing support

#### Operations & Deployment
- ✅ **Docker**: Multi-stage build, Alpine-based minimal image
- ✅ **Docker Compose**: Health checks and service orchestration
- ✅ **Kubernetes**: Deployment manifests (deployment, service, configmap)
- ✅ **Systemd Service**: Production service file
- ✅ **Reverse Proxy Configs**: Nginx and Caddy examples
- ✅ **Monitoring**: Prometheus and Grafana dashboard configs
- ✅ **Backup Scripts**: Automated database backups
- ✅ **Health Check Scripts**: Monitoring and alerting support

#### Developer Experience
- ✅ **Comprehensive Tests**: 42+ tests (unit, integration, benchmarks)
- ✅ **Example Code**: Usage examples and integration guides
- ✅ **Makefile**: Common development tasks
- ✅ **CI/CD**: GitHub Actions workflows
- ✅ **Documentation**: 18+ comprehensive guides
- ✅ **Linter Configuration**: GolangCI-Lint setup
- ✅ **Code Comments**: Well-documented codebase with inline documentation
- ✅ **OpenAPI Spec**: Complete API specification

#### API & Management
- ✅ **REST API**: `/api/v1/*` endpoints for bridge management
- ✅ **Web Admin Panel**: HTML dashboard with live statistics
- ✅ **Admin Commands**: `!bridge link`, `!bridge unlink`, `!bridge status`, `!bridge help`, `!bridge ping`
- ✅ **Bot Commands**: Viber bot command parsing and Matrix bridge
- ✅ **Outgoing Webhooks**: Matrix event forwarding for external integrations
- ✅ **Bridge Info API**: `/api/info` endpoint with status and statistics

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

# Configure environment
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

# Download dependencies
go mod tidy

# Build
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
| `VIBER_API_BASE_URL` | Viber API base URL (default: `https://chatapi.viber.com`) | No |
| `LISTEN_ADDRESS` | HTTP server listen address (default: `:8080`) | No |
| `MATRIX_HOMESERVER_URL` | Matrix homeserver base URL | Yes (if bridging) |
| `MATRIX_ACCESS_TOKEN` | Matrix access token | Yes (if bridging) |
| `MATRIX_DEFAULT_ROOM_ID` | Default Matrix room for bridged messages | Yes (if bridging) |
| `DATABASE_PATH` | SQLite database path (default: `./data/bridge.db`) | No |
| `HTTP_CLIENT_TIMEOUT` | HTTP client timeout in seconds (default: `15`) | No |
| `LOG_LEVEL` | Log level: debug, info, warn, error (default: `info`) | No |
| `VIBER_DEFAULT_RECEIVER_ID` | Default Viber user ID for Matrix → Viber forwarding | Optional |
| `REDIS_URL` | Redis URL for caching (optional, e.g. `redis://localhost:6379`) | No |
| `CACHE_TTL` | Cache TTL in minutes (default: `5`) | No |
| `ENABLE_PPROF` | Enable pprof endpoints for debugging (default: `false`) | No |
| `ENABLE_REQUEST_LOGGING` | Enable request/response body logging for debugging (default: `false`) | No |

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
- **GET** `/readyz` — Readiness check (returns 200 if ready, includes dependency checks)
- **GET** `/metrics` — Prometheus metrics
- **GET** `/api/info` — Bridge information and statistics (JSON)

### REST API

- **GET** `/api/v1/users` — List linked users
- **GET** `/api/v1/rooms` — List mapped rooms
- **POST** `/api/v1/link` — Link Matrix user to Viber user
- **POST** `/api/v1/unlink` — Unlink Matrix user from Viber

See [docs/API.md](docs/API.md) for complete API documentation and [docs/openapi.yaml](docs/openapi.yaml) for OpenAPI specification.

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
- **Panic Recovery**: Server crashes prevented with graceful error handling

### Best Practices

1. **Always use HTTPS** for webhook URLs in production
2. **Keep tokens secret**: Never commit tokens to version control
3. **Monitor metrics**: Watch `/metrics` for unusual patterns
4. **Review logs**: Structured logs help identify security issues
5. **Update regularly**: Keep dependencies updated for security patches

See [docs/SECURITY.md](docs/SECURITY.md) for comprehensive security guide.

---

## Metrics

The bridge exposes Prometheus metrics at `/metrics`:

- `viber_webhook_requests_total` — Total webhook requests by event type
- `viber_messages_forwarded_total` — Messages forwarded to Matrix by type
- `viber_signature_failures_total` — Signature verification failures
- `viber_message_latency_seconds` — Message processing latency

---

## Testing

### Running Tests

```bash
# Download dependencies first
go mod tidy

# Run all tests
go test ./... -v

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
go test -bench=. -benchmem ./test/benchmark/...

# Use test script
./scripts/test-all.sh
```

See [TESTING.md](TESTING.md) for comprehensive testing documentation.

### Test Coverage

- ✅ Database operations (5+ tests)
- ✅ Configuration (3+ tests)
- ✅ Retry logic (4+ tests)
- ✅ Circuit breaker (4+ tests)
- ✅ Validation utilities (5+ tests)
- ✅ Signature verification (3+ tests)
- ✅ Message queue (3+ tests)
- ✅ HTTP handlers (5+ tests)
- ✅ Integration tests (7+ tests)
- ✅ Benchmarks (4+ tests)

**Total**: 42+ test functions with comprehensive coverage

---

## Development

### Project Structure

```
cmd/
  mautrix-viber/          # Main application
internal/
  admin/                  # Admin commands
  api/                    # REST API endpoints
  cache/                  # Redis caching
  circuitbreaker/         # Circuit breaker pattern
  config/                 # Configuration management
  database/               # Database layer & migrations
  logger/                 # Structured logging
  matrix/                 # Matrix client & features
  metrics/                # Prometheus metrics
  middleware/             # HTTP middleware
  queue/                  # Message queue
  retry/                  # Retry logic
  tracing/                # OpenTelemetry tracing
  utils/                  # Utility functions
  version/                # Version management
  viber/                  # Viber client & features
  webadmin/               # Web admin panel
test/
  benchmark/              # Performance benchmarks
  http/                   # HTTP handler tests
  integration/            # Integration tests
docs/                     # Documentation
k8s/                      # Kubernetes manifests
monitoring/               # Monitoring configs
scripts/                  # Utility scripts
```

### Building

```bash
# Build binary
go build -o ./bin/mautrix-viber ./cmd/mautrix-viber

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint
golangci-lint run
```

### Adding Features

The codebase is structured for easy extension:

1. **Database operations**: Add methods to `internal/database/database.go`
2. **Viber API calls**: Add functions to `internal/viber/send.go`
3. **Matrix operations**: Extend `internal/matrix/client.go`
4. **Admin commands**: Register in `internal/admin/commands.go`

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for complete development guide.

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

See [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) for comprehensive troubleshooting guide.

---

## Documentation

- [README.md](README.md) - This file, getting started guide
- [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) - Production deployment guide
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) - System architecture overview
- [docs/API.md](docs/API.md) - REST API documentation
- [docs/openapi.yaml](docs/openapi.yaml) - OpenAPI specification
- [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) - Troubleshooting guide
- [docs/PERFORMANCE.md](docs/PERFORMANCE.md) - Performance tuning guide
- [docs/SECURITY.md](docs/SECURITY.md) - Security guide
- [docs/EXAMPLES.md](docs/EXAMPLES.md) - Configuration examples
- [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) - Development guide
- [docs/SEQUENCE_DIAGRAMS.md](docs/SEQUENCE_DIAGRAMS.md) - Message flow diagrams
- [docs/FAQ.md](docs/FAQ.md) - Frequently asked questions
- [docs/ROADMAP.md](docs/ROADMAP.md) - Development roadmap
- [TESTING.md](TESTING.md) - Testing guide and test coverage
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contributing guidelines
- [PRODUCTION_HARDENING.md](PRODUCTION_HARDENING.md) - Production hardening features

### Coding Standards

This project follows strict Go best practices. Key principles:

- ✅ **Error Handling**: All errors handled, never ignored
- ✅ **Testing**: Table-driven tests, **80%+ coverage target**
- ✅ **Security**: Input validation, parameterized queries, no secrets in code
- ✅ **Concurrency**: Safe goroutines with context cancellation
- ✅ **Logging**: Structured logging (slog), no `fmt.Printf`
- ✅ **Documentation**: Godoc comments for all exported functions
- ✅ **Code Quality**: Must pass `go vet` and `golangci-lint` with no warnings

See [CONTRIBUTING.md](CONTRIBUTING.md) for complete coding standards and guidelines.

---

## License

This project is licensed under the MIT License. See `LICENSE` for details.

---

## Contributing

Contributions welcome! Please review our coding standards before contributing.

Quick steps:

1. Review [CONTRIBUTING.md](CONTRIBUTING.md) for coding standards
2. Fork the repository
3. Create a feature branch
4. Add tests for new features (80%+ coverage target)
5. Ensure `go vet` and `golangci-lint` pass
6. Submit a pull request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines and coding standards.

---

## Acknowledgments

Built with:
- [mautrix-go](https://github.com/mautrix/go) — Matrix client library
- [Prometheus](https://prometheus.io/) — Metrics
- [log/slog](https://pkg.go.dev/log/slog) — Structured logging

---

**Status**: Production-ready with comprehensive feature set (49+ features), 42+ tests, and 18+ documentation files. Actively maintained and extended.
