# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2024-01-XX

### Added - Core Features
- Bidirectional message bridging between Matrix and Viber
- Full media support (text, images, video, audio, files, stickers, locations, contacts)
- Rich formatting (replies, threads, reactions, markdown parsing, mentions)
- Ghost user puppeting with avatar sync
- Portal room auto-creation for Viber chats
- Group chat support with member synchronization
- Typing indicators and read receipts synchronization
- Message edits and deletions (redactions)
- Message history backfill support

### Added - Infrastructure
- SQLite database for persistence
- Database migration tool
- Redis caching layer
- Message queue with retry logic
- Circuit breaker pattern for fault tolerance
- Advanced rate limiting (per-user, per-room, adaptive)

### Added - Observability
- Prometheus metrics integration
- OpenTelemetry tracing support
- Structured JSON logging (log/slog)
- Health check endpoints (/healthz, /readyz)
- Bridge info API endpoint (/api/info)

### Added - Security
- HMAC-SHA256 webhook signature verification
- Request body size limits
- Per-IP rate limiting
- Input validation and sanitization
- HTTPS enforcement for production

### Added - Deployment
- Docker containerization
- Docker Compose configuration
- Kubernetes manifests
- Systemd service file
- Reverse proxy configurations (Nginx, Caddy)

### Added - Operations
- Hot configuration reload (SIGHUP)
- Graceful shutdown
- Automatic webhook registration
- Backup and restore scripts
- Health check scripts

### Added - Developer Experience
- Comprehensive documentation
- Example code
- Unit tests
- Integration tests
- CI/CD workflows (GitHub Actions)
- Makefile for common tasks
- Contributing guidelines

### Added - Advanced Features
- Matrix E2EE support
- Admin commands (!bridge commands)
- Bot command parsing
- Power level synchronization
- Room metadata sync (name, topic, avatar)
- Notification configuration
- Message search capabilities
- Delivery receipt tracking
- Outgoing webhooks for integrations
- REST API for bridge management
- Web admin panel

### Added - Additional Enhancements
- HTTP request logging middleware
- Message queue for reliable delivery
- Input validation utilities
- Monitoring dashboards (Grafana)
- Deployment guides
- Architecture documentation

## Future Plans

- Enhanced E2EE support
- Improved transcoding for media
- Advanced search with full-text indexing
- WebSocket support for real-time updates
- Plugin system for custom handlers

