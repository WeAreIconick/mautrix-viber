# Frequently Asked Questions (FAQ)

Common questions and answers about mautrix-viber.

## General

### What is mautrix-viber?

mautrix-viber is a bidirectional bridge between Matrix and Viber messaging platforms. It allows users to send and receive messages between Matrix rooms and Viber chats.

### Is it production-ready?

Yes! The bridge is production-ready with:
- 49+ implemented features
- Comprehensive testing (42+ tests)
- Production hardening (panic recovery, health checks, etc.)
- Full documentation
- Security best practices

### What are the requirements?

- Go 1.22+ (for building)
- Viber Bot API token
- Matrix homeserver access
- Public HTTPS endpoint for webhooks
- SQLite database (included)

## Configuration

### Do I need Matrix to use the bridge?

Matrix is optional. The bridge can run in Viber-only mode, but full bidirectional bridging requires Matrix configuration.

### Can I use HTTP instead of HTTPS?

HTTP works for development and testing, but **HTTPS is required for production** webhook URLs. Viber enforces HTTPS for production webhooks.

### How do I get a Viber Bot API token?

1. Go to [Viber Partners Portal](https://partners.viber.com/)
2. Create a bot account
3. Get your API token from the dashboard

### How do I get a Matrix access token?

1. Login to your Matrix client (Element, etc.)
2. Go to Settings → Help & About → Advanced
3. Click "Access Token" to copy it

Or use:
```bash
curl -X POST "https://matrix.example.com/_matrix/client/r0/login" \
  -d '{"type":"m.login.password","user":"@user:example.com","password":"password"}'
```

## Deployment

### Can I run multiple instances?

Yes, with a shared database. However, SQLite doesn't handle concurrent writes well, so consider:
- Using PostgreSQL adapter (future feature)
- Single instance with multiple workers
- Sticky sessions with load balancer

### How do I scale the bridge?

1. **Vertical scaling**: Increase CPU/memory resources
2. **Horizontal scaling**: Multiple instances with shared database (careful with SQLite)
3. **Optimization**: Enable caching, optimize queries
4. **Future**: PostgreSQL support for better concurrency

### What ports need to be open?

- **Inbound**: Port 8080 (or your configured `LISTEN_ADDRESS`)
- **Outbound**: HTTPS (443) to Viber API and Matrix homeserver

## Troubleshooting

### Messages not bridging

1. Check logs with `LOG_LEVEL=debug`
2. Verify Matrix credentials
3. Check room ID format
4. Review `/api/info` endpoint
5. Test Matrix connection manually

### High memory usage

1. Check for goroutine leaks
2. Review connection pool settings
3. Enable profiling: `go tool pprof http://localhost:8080/debug/pprof/heap`
4. Monitor metrics

### Database errors

1. Check file permissions
2. Verify disk space
3. Check for concurrent access issues
4. Review connection pool settings
5. Use WAL mode (default)

## Security

### Are my tokens secure?

The bridge:
- Never logs tokens
- Doesn't expose tokens in errors
- Uses environment variables or secrets management
- Validates all inputs

**You should**:
- Use secret management (Kubernetes Secrets, Vault, etc.)
- Rotate tokens regularly
- Never commit tokens to version control

### How secure is the webhook endpoint?

The bridge:
- Verifies HMAC-SHA256 signatures on all requests
- Rate limits per IP
- Limits request body size
- Validates all inputs
- Recovers from panics securely

### Can I add authentication to the API?

Yes! You can add authentication middleware. The bridge provides the foundation, you can add:
- API key authentication
- OAuth2
- Basic authentication
- Custom authentication

## Performance

### How many messages per second can it handle?

The bridge can handle:
- ~100-500 messages/second (depending on hardware)
- Higher with optimization (caching, connection pooling)
- Bottleneck is usually external APIs (Viber/Matrix)

### How much memory does it use?

Typical usage:
- **Idle**: ~50-100MB
- **Normal load**: ~100-256MB
- **High load**: ~256-512MB

Monitor with metrics and adjust limits accordingly.

### Can I optimize for my use case?

Yes! See [PERFORMANCE.md](PERFORMANCE.md) for:
- Database optimization
- Memory tuning
- Network optimization
- Caching strategies

## Features

### Does it support group chats?

Yes! Group chat support includes:
- Viber group → Matrix room mapping
- Member synchronization
- Group metadata sync

### Does it support E2EE?

Matrix E2EE support is implemented, but full end-to-end encryption bridging has limitations. See documentation for details.

### Can I bridge multiple Viber bots?

Currently, one bridge instance handles one Viber bot. For multiple bots, run multiple bridge instances.

## Contributing

### How do I contribute?

1. Fork the repository
2. Create a feature branch
3. Make changes following `.cursorrules`
4. Add tests
5. Update documentation
6. Submit pull request

See [CONTRIBUTING.md](../CONTRIBUTING.md) for details.

### What coding standards should I follow?

See [.cursorrules](../.cursorrules) for comprehensive coding standards and best practices.

## Support

### Where can I get help?

- **Documentation**: Check the `docs/` directory
- **Issues**: GitHub Issues
- **Discussions**: GitHub Discussions (if enabled)
- **Troubleshooting**: See [TROUBLESHOOTING.md](TROUBLESHOOTING.md)

### How do I report bugs?

Open a GitHub issue with:
- Description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Logs (with sensitive data redacted)
- Environment details

