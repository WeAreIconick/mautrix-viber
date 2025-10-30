# Security Guide

Security best practices and audit checklist for mautrix-viber.

## Security Features

### ✅ Implemented

- **HMAC-SHA256 Signature Verification**: All webhooks verified
- **Rate Limiting**: Per-IP token bucket algorithm
- **Input Validation**: Comprehensive sanitization
- **HTTPS Enforcement**: Production requirement
- **Request Body Limits**: 2MB maximum
- **Panic Recovery**: Prevents information leakage
- **Error Sanitization**: Internal errors not exposed to clients

## Security Checklist

### Configuration Security

- [ ] **Secrets Management**
  - [ ] No secrets committed to version control
  - [ ] Use environment variables or secret management (Vault, AWS Secrets Manager)
  - [ ] Rotate API tokens regularly
  - [ ] Use different tokens for development/production

- [ ] **Network Security**
  - [ ] Webhook URL uses HTTPS
  - [ ] Firewall rules restrict access
  - [ ] VPN or private network for Matrix connection
  - [ ] Rate limiting enabled

- [ ] **Access Control**
  - [ ] Database file permissions restricted (chmod 600)
  - [ ] Service runs as non-root user
  - [ ] File system permissions reviewed
  - [ ] API endpoints protected if exposed publicly

### Runtime Security

- [ ] **Webhook Security**
  - [ ] Signature verification enabled
  - [ ] Invalid signatures logged
  - [ ] Signature failures monitored
  - [ ] Webhook URL not guessable

- [ ] **Input Validation**
  - [ ] All user inputs validated
  - [ ] SQL injection prevented (use prepared statements)
  - [ ] XSS prevention in admin panel
  - [ ] Path traversal prevention

- [ ] **Error Handling**
  - [ ] Internal errors not exposed
  - [ ] Stack traces only in debug mode
  - [ ] Error messages don't leak information
  - [ ] Failed requests logged appropriately

### Dependencies Security

- [ ] **Dependency Management**
  - [ ] Regular dependency updates
  - [ ] Security advisories monitored
  - [ ] `go mod tidy` run regularly
  - [ ] `go.sum` checked into version control

- [ ] **Vulnerability Scanning**
  ```bash
  # Use govulncheck
  go install golang.org/x/vuln/cmd/govulncheck@latest
  govulncheck ./...
  
  # Use Dependabot or similar for automated scanning
  ```

### Deployment Security

- [ ] **Container Security**
  - [ ] Base image regularly updated
  - [ ] Minimal base image (Alpine)
  - [ ] No unnecessary packages
  - [ ] Non-root user in container

- [ ] **Kubernetes Security**
  - [ ] Secrets in Kubernetes Secrets, not ConfigMaps
  - [ ] Pod security policies enabled
  - [ ] Network policies configured
  - [ ] RBAC properly configured

- [ ] **Network Security**
  - [ ] HTTPS/TLS for all external endpoints
  - [ ] TLS certificates valid and not expired
  - [ ] CORS configured appropriately
  - [ ] Security headers set (HSTS, CSP, etc.)

## Security Headers

Add these headers for web endpoints:

```go
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("X-XSS-Protection", "1; mode=block")
w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
```

## Audit Logging

### What to Log

- All webhook signature verification failures
- Authentication/authorization failures
- Rate limit violations
- Configuration changes
- Database access patterns (anomalies)

### Log Retention

- Structured logs retained for 30-90 days
- Security events retained longer
- Comply with data retention requirements

## Incident Response

### If Compromised

1. **Immediate Actions**
   - Rotate all API tokens
   - Review recent logs for suspicious activity
   - Check database for unauthorized changes
   - Revoke compromised access tokens

2. **Investigation**
   - Review audit logs
   - Identify attack vector
   - Assess data exposure
   - Document timeline

3. **Remediation**
   - Patch vulnerabilities
   - Update dependencies
   - Improve security controls
   - Notify affected parties if required

## Security Testing

### Regular Audits

1. **Dependency Audits**
   ```bash
   govulncheck ./...
   npm audit  # if using npm tools
   ```

2. **Code Review**
   - Security-focused code reviews
   - Static analysis tools
   - SAST (Static Application Security Testing)

3. **Penetration Testing**
   - Webhook endpoint testing
   - API security testing
   - Infrastructure security review

## Secure Development Practices

1. **Never commit secrets**
2. **Always validate input**
3. **Use prepared statements for SQL**
4. **Handle errors securely**
5. **Keep dependencies updated**
6. **Follow principle of least privilege**
7. **Enable security features by default**
8. **Security through obscurity is not security**

## Compliance

### GDPR Considerations

- User data stored securely
- Right to deletion implemented
- Data export capability
- Privacy policy available

### Data Protection

- Encrypt sensitive data at rest
- Encrypt data in transit (TLS)
- Secure key management
- Regular backups with encryption

## Reporting Security Issues

Report security vulnerabilities to: [security@example.com]

Please include:
- Description of vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if available)

