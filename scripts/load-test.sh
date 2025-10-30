#!/bin/bash
# Load testing script for mautrix-viber bridge

set -e

BRIDGE_URL="${BRIDGE_URL:-http://localhost:8080}"
CONCURRENT_USERS="${CONCURRENT_USERS:-10}"
REQUESTS_PER_USER="${REQUESTS_PER_USER:-100}"

echo "Load Testing mautrix-viber Bridge"
echo "URL: $BRIDGE_URL"
echo "Concurrent Users: $CONCURRENT_USERS"
echo "Requests per User: $REQUESTS_PER_USER"
echo ""

# Check if bridge is healthy
if ! curl -sf "$BRIDGE_URL/healthz" > /dev/null; then
    echo "Error: Bridge is not healthy"
    exit 1
fi

# Load test health endpoint
echo "Testing /healthz endpoint..."
ab -n $((CONCURRENT_USERS * REQUESTS_PER_USER)) -c $CONCURRENT_USERS "$BRIDGE_URL/healthz"

# Load test ready endpoint
echo ""
echo "Testing /readyz endpoint..."
ab -n $((CONCURRENT_USERS * REQUESTS_PER_USER)) -c $CONCURRENT_USERS "$BRIDGE_URL/readyz"

# Load test info endpoint
echo ""
echo "Testing /api/info endpoint..."
ab -n $((CONCURRENT_USERS * REQUESTS_PER_USER)) -c $CONCURRENT_USERS "$BRIDGE_URL/api/info"

echo ""
echo "Load test complete!"

# Note: Requires Apache Bench (ab) or similar tool
# Install: apt-get install apache2-utils

