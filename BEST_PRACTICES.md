# Best Practices Summary

This document summarizes the key best practices used in the mautrix-viber project. For detailed guidelines, see [.cursorrules](.cursorrules).

## Code Quality Standards

### ✅ Implemented Best Practices

1. **Error Handling**
   - All errors are explicitly handled
   - Errors are wrapped with context using `fmt.Errorf` with `%w`
   - Use `errors.Is()` and `errors.As()` for error checking
   - Never ignore errors silently

2. **Context Usage**
   - All I/O operations accept `context.Context` as first parameter
   - Context cancellation is checked before long operations
   - Timeouts are set for external API calls

3. **Documentation**
   - All exported functions/types have package comments
   - Complex logic is explained with inline comments
   - Package comments explain the purpose and usage

4. **Testing**
   - 26+ comprehensive unit tests
   - Table-driven tests for multiple test cases
   - Test helpers use `t.Helper()`
   - Tests clean up resources with `defer`

5. **Security**
   - HMAC-SHA256 signature verification for webhooks
   - Input validation and sanitization
   - Rate limiting to prevent abuse
   - No secrets in code or logs

6. **Code Organization**
   - Clear package structure
   - No circular dependencies
   - Interfaces for testability
   - Single responsibility principle

## Key Patterns

### Dependency Injection
- Clients accept dependencies via constructor
- Nil checks for optional dependencies
- Interface-based design for testability

### Resource Management
- Always use `defer` for cleanup
- Close resources (files, DB connections, HTTP bodies)
- Use timeouts for external calls

### Logging
- Structured logging with `log/slog`
- Include relevant context in log entries
- Appropriate log levels (debug, info, warn, error)

### Configuration
- Environment variables with defaults
- Configuration validation on startup
- Sensible defaults for optional settings

## Project-Specific Guidelines

### Bridge Operations
- Always check if clients are configured before use
- Store message mappings for deduplication
- Use retry logic for external API calls
- Handle all webhook event types

### Database
- Use transactions for multi-step operations
- Handle `sql.ErrNoRows` explicitly
- Use prepared statements
- Set connection timeouts

### HTTP Handlers
- Verify signatures before processing
- Parse request body once
- Return appropriate HTTP status codes
- Don't expose internal errors to clients

See [.cursorrules](.cursorrules) for the complete coding standards and guidelines.

