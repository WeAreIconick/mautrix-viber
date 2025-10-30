## mautrix-viber — Matrix ↔ Viber bridge (Go)

Robust, minimal scaffolding for a Matrix ↔ Viber bridge written in Go. It exposes an HTTP webhook endpoint for Viber callbacks and lays the groundwork to forward events to Matrix using `maunium.net/go/mautrix`.

This repository is currently a foundation: it starts a web server, defines Viber types, and contains a stubbed webhook registration call. It’s intended for developers who want to build a full-featured bridge while starting from a clean, idiomatic codebase.

### Status
- Early-stage scaffold. Safe for local development and experimentation.
- Webhook registration is a placeholder (no outbound API call yet).
- No message relaying to Matrix implemented yet.

---

## Features
- HTTP server with a Viber webhook handler at `/viber/webhook` returning `200 OK`.
- Minimal Viber types for common events (message, subscribe, unsubscribe, conversation started).
- Environment-based configuration helper (`internal/config`).
- Clean package structure ready to extend with Matrix bridging logic.

---

## How it works (high level)
1. The process starts an HTTP server (defaults to `:8080`).
2. On startup, it attempts to ensure the Viber webhook is registered (currently a no-op stub).
3. Viber sends callbacks to your public webhook URL which should route to `/viber/webhook`.
4. The webhook handler currently acknowledges with `200 OK`; you can extend it to relay to Matrix.

---

## Requirements
- Go 1.22+
- A Viber Bot and its API token
- A publicly reachable HTTPS URL for the webhook (e.g., via a reverse proxy or a tunneling tool like `ngrok`)
- Optional: Matrix homeserver URL, access token, and a target room ID to relay incoming Viber messages

---

## Configuration
There are two configuration structs in the codebase:

- `internal/config.Config` supports environment variables:
  - `VIBER_API_TOKEN`: Viber bot token
  - `VIBER_WEBHOOK_URL`: Public HTTPS URL for Viber to call (e.g. `https://your.domain.tld/viber/webhook`)
  - `LISTEN_ADDRESS`: Address for the HTTP server (default `:8080`)
  - `MATRIX_HOMESERVER_URL`: Matrix homeserver base URL (e.g. `https://matrix.example.com`)
  - `MATRIX_ACCESS_TOKEN`: Matrix access token of the bot/user to send with
  - `MATRIX_DEFAULT_ROOM_ID`: Matrix room ID to receive relayed messages (e.g. `!roomid:example.com`)

- `internal/viber.Config` is the runtime config used by the Viber client. In `cmd/mautrix-viber/main.go` the values are currently hard-coded as empty placeholders. You can either:
  - Replace the hard-coded config with `internal/config.FromEnv()` to load from environment variables, or
  - Manually set the values in `main.go`.

Example environment configuration (recommended):

```bash
export VIBER_API_TOKEN="viber-xxxxxxxxxxxxxxxx"
export VIBER_WEBHOOK_URL="https://your.public.url/viber/webhook"
export LISTEN_ADDRESS=":8080"
export MATRIX_HOMESERVER_URL="https://matrix.example.com"
export MATRIX_ACCESS_TOKEN="syt_xxxxxxxxxxxxxxxxxxxxxxxxx"
export MATRIX_DEFAULT_ROOM_ID="!roomid:example.com"
```

---

## Quickstart

### 1) Build
```bash
go build -o ./bin/mautrix-viber ./cmd/mautrix-viber
```

### 2) Run (local)
```bash
./bin/mautrix-viber
```

You should see a log line similar to:

```text
listening on :8080
```

### 3) Expose a public URL (for Viber)
Use your preferred method (reverse proxy, tunnel, or production ingress). For quick testing with `ngrok`:

```bash
ngrok http 8080
```

Copy the HTTPS URL (e.g., `https://abcd1234.ngrok.io`) and set `VIBER_WEBHOOK_URL` to `https://abcd1234.ngrok.io/viber/webhook`.

### 4) Register the Viber webhook

Webhook registration is currently a stub in code (`EnsureWebhook` does nothing). Until that’s implemented, set the webhook manually with Viber’s API.

Replace `VIBER_API_TOKEN` and `PUBLIC_URL` and run:

```bash
curl -X POST \
  -H "X-Viber-Auth-Token: ${VIBER_API_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
        "url": "'"${PUBLIC_URL}/viber/webhook"'",
        "event_types": ["message", "subscribed", "unsubscribed", "conversation_started"]
      }' \
  https://chatapi.viber.com/pa/set_webhook
```

If successful, Viber will start sending events to your endpoint and return a JSON payload with `status` and `status_message`.

### 5) Verify Matrix relay (optional)
If Matrix credentials are configured, send a text message to the Viber bot. The bridge will post a message like `"[Viber] Alice: Hello"` into `MATRIX_DEFAULT_ROOM_ID`.

---

## Endpoints
- `POST /viber/webhook` — receives Viber callbacks. Currently responds with `200 OK` without additional processing.

Health endpoints are not yet implemented; consider adding one (e.g., `GET /healthz`) for production deployments.

---

## Project layout

```
cmd/
  mautrix-viber/
    main.go           # Entry point (HTTP server, webhook setup)
internal/
  config/
    config.go         # Env-driven config loader
  matrix/
    bridge.go         # Bridge interface (placeholder)
  viber/
    client.go         # Viber client + webhook handler
    types.go          # Minimal Viber types (events, payloads)
go.mod
```

---

## Development

### Run tests
There are currently no tests. As you add functionality, prefer table-driven unit tests for handlers and Matrix/Viber adapters.

### Lint/format
Use your preferred Go linters/formatters. A typical setup:

```bash
go fmt ./...
go vet ./...
golangci-lint run
```

### Hot reload (optional)
For iterative development:

```bash
go install github.com/cosmtrek/air@latest
air
```

---

## Extending to a real bridge

Suggested next steps:
- Implement Viber webhook registration in `viber.Client.EnsureWebhook()` using `chatapi.viber.com/pa/set_webhook`.
- Parse incoming `WebhookRequest` bodies in `WebhookHandler` and validate signatures (if applicable).
- Map Viber events to Matrix events using `maunium.net/go/mautrix`.
- Maintain user/room mappings (Matrix user ↔ Viber user, Matrix room ↔ Viber chat).
- Add persistence for state (user links, message IDs) using a lightweight DB.
- Implement backfill/history, media handling, typing/read receipts as needed.

Security considerations:
- Validate webhook origin and/or signature headers from Viber.
- Use HTTPS everywhere; terminate TLS at the edge if needed.
- Avoid logging sensitive content (tokens, message bodies) in production.

---

## FAQ

**Why doesn’t the webhook register automatically?**
The `EnsureWebhook` method is a placeholder. Use the curl snippet above for now, or implement the call to Viber’s API.

**Does it forward messages to Matrix already?**
Not yet. The repo sets up the plumbing so you can add that next.

**Docker? Helm?**
Not provided yet. Contributions welcome.

---

## License
This project is licensed under the MIT License. See `LICENSE` for details.

# mautrix-viber

Matrix bridge for Viber using mautrix-go bridgev2.
