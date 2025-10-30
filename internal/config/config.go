// Package config provides configuration loading from environment variables
// with sensible defaults and validation.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all bridge configuration settings.
// Fields are loaded from environment variables via FromEnv().
type Config struct {
	// Viber API configuration
	APIToken      string // Viber Bot API token (required)
	WebhookURL    string // Public HTTPS URL for Viber webhooks (required)
	ViberAPIBaseURL string // Viber API base URL (default: "https://chatapi.viber.com")
	ListenAddress string // HTTP server listen address (default: ":8080")
	
	// Matrix client configuration
	MatrixHomeserverURL string // Matrix homeserver base URL (required if bridging)
	MatrixAccessToken   string // Matrix access token (required if bridging)
	MatrixDefaultRoomID string // Default Matrix room for bridged messages (required if bridging)
	
	// Optional features
	ViberDefaultReceiverID string        // Default Viber user ID for Matrix → Viber forwarding (optional)
	DatabasePath           string        // SQLite database path (default: "./data/bridge.db")
	HTTPClientTimeout      time.Duration // HTTP client timeout for API calls (default: 15s)
}

// FromEnv loads configuration from environment variables.
// Returns a Config with default values applied where appropriate.
// Environment variables override any default values.
func FromEnv() Config {
	cfg := Config{}
	cfg.APIToken = os.Getenv("VIBER_API_TOKEN")
	cfg.WebhookURL = os.Getenv("VIBER_WEBHOOK_URL")
	cfg.ViberAPIBaseURL = os.Getenv("VIBER_API_BASE_URL")
	if cfg.ViberAPIBaseURL == "" {
		cfg.ViberAPIBaseURL = "https://chatapi.viber.com"
	}
	cfg.ListenAddress = os.Getenv("LISTEN_ADDRESS")
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = ":8080"
	}
    cfg.MatrixHomeserverURL = os.Getenv("MATRIX_HOMESERVER_URL")
    cfg.MatrixAccessToken = os.Getenv("MATRIX_ACCESS_TOKEN")
    cfg.MatrixDefaultRoomID = os.Getenv("MATRIX_DEFAULT_ROOM_ID")
    cfg.ViberDefaultReceiverID = os.Getenv("VIBER_DEFAULT_RECEIVER_ID")
    cfg.DatabasePath = os.Getenv("DATABASE_PATH")
    if cfg.DatabasePath == "" {
        cfg.DatabasePath = "./data/bridge.db"
    }
	// HTTP client timeout (in seconds)
	if timeoutStr := os.Getenv("HTTP_CLIENT_TIMEOUT"); timeoutStr != "" {
		if timeoutSec, err := strconv.Atoi(timeoutStr); err == nil && timeoutSec > 0 {
			cfg.HTTPClientTimeout = time.Duration(timeoutSec) * time.Second
		} else {
			cfg.HTTPClientTimeout = 15 * time.Second // Default on parse error
		}
	} else {
		cfg.HTTPClientTimeout = 15 * time.Second // Default
	}
	return cfg
}
