# Configuration Examples

Practical configuration examples for different deployment scenarios.

## Basic Setup

### Environment Variables

```bash
# Required
export VIBER_API_TOKEN="your-viber-bot-token"
export VIBER_WEBHOOK_URL="https://bridge.example.com/viber/webhook"

# Matrix (required for bridging)
export MATRIX_HOMESERVER_URL="https://matrix.example.com"
export MATRIX_ACCESS_TOKEN="syt_..."
export MATRIX_DEFAULT_ROOM_ID="!roomid:example.com"

# Optional
export LISTEN_ADDRESS=":8080"
export DATABASE_PATH="./data/bridge.db"
export LOG_LEVEL="info"
```

## Production Deployment

### Docker Compose

```yaml
version: '3.8'

services:
  mautrix-viber:
    image: mautrix-viber:latest
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - VIBER_API_TOKEN=${VIBER_API_TOKEN}
      - VIBER_WEBHOOK_URL=https://bridge.example.com/viber/webhook
      - MATRIX_HOMESERVER_URL=https://matrix.example.com
      - MATRIX_ACCESS_TOKEN=${MATRIX_ACCESS_TOKEN}
      - MATRIX_DEFAULT_ROOM_ID=${MATRIX_DEFAULT_ROOM_ID}
      - DATABASE_PATH=/data/bridge.db
      - LOG_LEVEL=info
    volumes:
      - ./data:/data
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/healthz"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### Kubernetes

```yaml
# See k8s/deployment.yaml for full example
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mautrix-viber
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: mautrix-viber
        image: mautrix-viber:latest
        env:
        - name: VIBER_API_TOKEN
          valueFrom:
            secretKeyRef:
              name: mautrix-viber-secrets
              key: viber-api-token
        # ... more config
```

## Development Setup

### Local Development with ngrok

```bash
# Terminal 1: Run bridge
export VIBER_API_TOKEN="dev-token"
export VIBER_WEBHOOK_URL="http://localhost:8080/viber/webhook"
./bin/mautrix-viber

# Terminal 2: Expose with ngrok
ngrok http 8080

# Update VIBER_WEBHOOK_URL to ngrok HTTPS URL
export VIBER_WEBHOOK_URL="https://abc123.ngrok.io/viber/webhook"
```

### Testing Configuration

```bash
# Minimal config for testing
export VIBER_API_TOKEN="test-token"
export VIBER_WEBHOOK_URL="https://test.example.com/webhook"
export LOG_LEVEL="debug"
```

## High Availability

### Load Balancer Configuration

```nginx
upstream mautrix_viber {
    least_conn;
    server bridge1:8080;
    server bridge2:8080;
}

server {
    listen 443 ssl;
    server_name bridge.example.com;

    location /viber/webhook {
        proxy_pass http://mautrix_viber;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        
        # Preserve body for signature verification
        proxy_pass_request_body on;
        proxy_pass_request_headers on;
        
        # Important: Don't buffer request body
        proxy_request_buffering off;
    }
}
```

### Shared Database (Multiple Instances)

```bash
# Use external database (e.g., PostgreSQL with SQLite compatibility)
# Or shared NFS/network storage for SQLite
export DATABASE_PATH="/shared/bridge.db"
```

## Performance Tuning

### High-Volume Configuration

```bash
# Increase connection pool
# (modify in database.go or via config)
export DATABASE_MAX_CONNS=50
export DATABASE_MAX_IDLE=10

# Enable Redis caching
export REDIS_URL="redis://localhost:6379"
export CACHE_TTL="5m"

# Adjust rate limits (in code)
# Increase base rate for high capacity
```

### Resource-Constrained Environment

```bash
# Reduce connection pool
# Limit memory usage
# Use smaller cache TTL
# Reduce rate limits
```

## Monitoring Setup

### Prometheus Scraping

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'mautrix-viber'
    static_configs:
      - targets: ['bridge1:8080', 'bridge2:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Grafana Dashboard

Import dashboard from `monitoring/grafana-dashboard.json` or create custom:

```json
{
  "dashboard": {
    "panels": [
      {
        "title": "Messages Bridged",
        "targets": [{
          "expr": "rate(viber_messages_forwarded_total[5m])"
        }]
      }
    ]
  }
}
```

## Troubleshooting Configurations

### Debug Mode

```bash
export LOG_LEVEL=debug
export VIBER_API_TOKEN="..."
export VIBER_WEBHOOK_URL="..."
```

### Verbose Logging

```bash
# Enable source location in logs
export LOG_LEVEL=debug
export LOG_SOURCE=true
```

## Security Configurations

### Production Security

```bash
# Enforce HTTPS
export VIBER_WEBHOOK_URL="https://bridge.example.com/viber/webhook"

# Use secrets management
# Don't set tokens in environment directly
# Use secret injection (Kubernetes Secrets, AWS Secrets Manager, etc.)
```

### Development Security

```bash
# Still use HTTPS in dev (ngrok provides this)
# Use separate dev tokens
# Don't commit real tokens
```

## Example YAML Config

```yaml
# config.yaml
viber:
  api_token: "${VIBER_API_TOKEN}"  # Use env var substitution
  webhook_url: "https://bridge.example.com/viber/webhook"

matrix:
  homeserver_url: "https://matrix.example.com"
  access_token: "${MATRIX_ACCESS_TOKEN}"
  default_room_id: "!roomid:example.com"

server:
  listen_address: ":8080"
  read_timeout: "10s"
  write_timeout: "15s"

database:
  path: "./data/bridge.db"
  max_open_conns: 25
  max_idle_conns: 5

logging:
  level: "info"
  format: "json"

cache:
  enabled: true
  redis_url: "redis://localhost:6379"
  ttl: "5m"

rate_limiting:
  enabled: true
  requests_per_second: 5
  burst_size: 10
```

