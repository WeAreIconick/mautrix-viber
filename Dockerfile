# Multi-stage build for mautrix-viber bridge
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o mautrix-viber ./cmd/mautrix-viber

# Final stage
FROM alpine:latest

# Install SQLite and CA certificates
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/mautrix-viber .

# Create data directory for SQLite
RUN mkdir -p /data && chmod 755 /data

# Expose default port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Run the bridge
CMD ["./mautrix-viber"]

