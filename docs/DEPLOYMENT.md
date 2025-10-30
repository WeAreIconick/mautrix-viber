# Deployment Guide

This guide covers deploying mautrix-viber in various environments.

## Prerequisites

- Viber Bot API token
- Matrix homeserver URL and access token
- Public HTTPS endpoint for webhooks
- Database storage (SQLite or external)

## Docker Deployment

### Quick Start

```bash
# Clone repository
git clone https://github.com/example/mautrix-viber.git
cd mautrix-viber

# Configure environment
cp .env.example .env
# Edit .env with your credentials

# Run with docker-compose
docker-compose up -d

# Check logs
docker-compose logs -f
```

### Production Docker Deployment

1. **Configure environment variables**:
```bash
export VIBER_API_TOKEN="your-token"
export VIBER_WEBHOOK_URL="https://your-domain.com/viber/webhook"
export MATRIX_HOMESERVER_URL="https://matrix.example.com"
export MATRIX_ACCESS_TOKEN="syt_..."
export MATRIX_DEFAULT_ROOM_ID="!roomid:example.com"
```

2. **Build and run**:
```bash
docker build -t mautrix-viber:latest .
docker run -d \
  --name mautrix-viber \
  -p 8080:8080 \
  -v ./data:/data \
  --env-file .env \
  mautrix-viber:latest
```

## Kubernetes Deployment

### Deployment Manifest

See `k8s/deployment.yaml` for complete Kubernetes deployment configuration.

### Basic Kubernetes Setup

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
```

## Systemd Service

Create `/etc/systemd/system/mautrix-viber.service`:

```ini
[Unit]
Description=mautrix-viber Bridge
After=network.target

[Service]
Type=simple
User=mautrix-viber
WorkingDirectory=/opt/mautrix-viber
ExecStart=/opt/mautrix-viber/bin/mautrix-viber
Restart=always
RestartSec=10
EnvironmentFile=/etc/mautrix-viber/env

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable mautrix-viber
sudo systemctl start mautrix-viber
sudo systemctl status mautrix-viber
```

## Reverse Proxy Setup

### Nginx Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name bridge.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /viber/webhook {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        # Viber requires original request for signature verification
        proxy_pass_request_headers on;
    }
}
```

### Caddy Configuration

```
bridge.example.com {
    reverse_proxy localhost:8080 {
        header_up X-Real-IP {remote_host}
        header_up X-Forwarded-For {remote_host}
        header_up X-Forwarded-Proto {scheme}
    }
}
```

## Health Checks

The bridge exposes health endpoints:

- `GET /healthz` - Health check (200 if healthy)
- `GET /readyz` - Readiness check (200 if ready to serve)
- `GET /metrics` - Prometheus metrics
- `GET /api/info` - Bridge status and statistics

Configure your orchestration system to use these endpoints.

## Monitoring

### Prometheus Scraping

Add to `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'mautrix-viber'
    static_configs:
      - targets: ['bridge.example.com:8080']
```

### Grafana Dashboard

Import the Grafana dashboard from `monitoring/grafana-dashboard.json`.

## Backup and Restore

### Database Backup

```bash
# Backup SQLite database
sqlite3 data/bridge.db ".backup backup.db"

# Restore
sqlite3 data/bridge.db ".restore backup.db"
```

### Automated Backups

Set up cron job:
```bash
0 2 * * * sqlite3 /path/to/bridge.db ".backup /backup/bridge-$(date +\%Y\%m\%d).db"
```

## Troubleshooting

### Common Issues

1. **Webhook not receiving events**
   - Verify webhook URL is publicly accessible
   - Check signature verification logs
   - Ensure webhook is registered with Viber

2. **Messages not bridging**
   - Verify Matrix credentials
   - Check database connectivity
   - Review application logs

3. **High error rates**
   - Check Prometheus metrics
   - Review rate limiting settings
   - Verify API rate limits aren't exceeded

### Debug Mode

Enable debug logging:
```bash
export LOG_LEVEL=debug
./bin/mautrix-viber
```

## Performance Tuning

- Adjust worker pool size for message processing
- Configure Redis cache for frequently accessed data
- Tune database connection pool settings
- Adjust rate limiting thresholds based on usage

## Security Checklist

- [ ] Use HTTPS for all external endpoints
- [ ] Keep tokens in secure storage (secrets manager)
- [ ] Enable signature verification
- [ ] Configure rate limiting appropriately
- [ ] Regularly update dependencies
- [ ] Monitor for security advisories
- [ ] Use firewall rules to restrict access
- [ ] Enable audit logging for sensitive operations

