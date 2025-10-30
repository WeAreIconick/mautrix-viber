#!/bin/bash
# Health check script for mautrix-viber bridge

BRIDGE_URL="${BRIDGE_URL:-http://localhost:8080}"

# Check health endpoint
HEALTH=$(curl -s -o /dev/null -w "%{http_code}" "$BRIDGE_URL/healthz")

if [ "$HEALTH" = "200" ]; then
    echo "✓ Bridge is healthy"
    exit 0
else
    echo "✗ Bridge health check failed (HTTP $HEALTH)"
    exit 1
fi

