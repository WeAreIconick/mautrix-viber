# API Documentation

## REST API Endpoints

### Bridge Information

#### GET /api/info
Returns bridge status and statistics.

**Response:**
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

### User Management

#### POST /api/v1/link
Link a Matrix user to a Viber user.

**Request:**
```json
{
  "matrix_user_id": "@user:example.com",
  "viber_user_id": "viber_user_123"
}
```

**Response:**
```json
{
  "status": "linked",
  "matrix_user_id": "@user:example.com",
  "viber_user_id": "viber_user_123"
}
```

#### POST /api/v1/unlink
Unlink a Matrix user from Viber.

**Request:**
```json
{
  "matrix_user_id": "@user:example.com"
}
```

### Room Management

#### GET /api/v1/rooms
List all mapped rooms.

**Response:**
```json
{
  "rooms": [
    {
      "matrix_room_id": "!room:example.com",
      "viber_chat_id": "viber_chat_123",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Health Checks

#### GET /healthz
Health check endpoint. Returns 200 if bridge is healthy.

#### GET /readyz
Readiness check endpoint. Returns 200 if bridge is ready to serve requests.

### Metrics

#### GET /metrics
Prometheus metrics endpoint. Returns metrics in Prometheus format.

## Matrix Admin Commands

Commands can be run in Matrix rooms:

- `!bridge help` - Show available commands
- `!bridge link <viber-user-id>` - Link Viber account
- `!bridge unlink` - Unlink Viber account
- `!bridge status` - Show bridge status
- `!bridge ping` - Test bridge responsiveness

## Webhook Endpoints

### POST /viber/webhook
Receives Viber webhook callbacks.

**Headers:**
- `X-Viber-Content-Signature`: HMAC-SHA256 signature

**Request Body:**
```json
{
  "event": "message",
  "sender": {
    "id": "viber_user_123",
    "name": "Alice"
  },
  "message": {
    "type": "text",
    "text": "Hello, world!"
  }
}
```

## Error Responses

All endpoints may return standard HTTP error codes:

- `400 Bad Request` - Invalid request format
- `401 Unauthorized` - Authentication failed
- `404 Not Found` - Resource not found
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service temporarily unavailable

Error responses include a JSON body:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

