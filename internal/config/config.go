// Package config provides configuration loading from environment variables
// with sensible defaults and validation.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/example/mautrix-viber/internal/utils"
)

// Config holds all bridge configuration settings.
// Fields are loaded from environment variables via FromEnv().
// All configuration is documented in the README.
type Config struct {
	// Viber API configuration
	APIToken        string // Viber Bot API token (required)
	WebhookURL      string // Public HTTPS URL for Viber webhooks (required)
	ViberAPIBaseURL string // Viber API base URL (default: "https://chatapi.viber.com")
	ListenAddress   string // HTTP server listen address (default: ":8080")

	// Matrix client configuration
	MatrixHomeserverURL string // Matrix homeserver base URL (required if bridging)
	MatrixAccessToken   string // Matrix access token (required if bridging)
	MatrixDefaultRoomID string // Default Matrix room for bridged messages (required if bridging)

	// Optional features
	ViberDefaultReceiverID string        // Default Viber user ID for Matrix → Viber forwarding (optional)
	DatabasePath           string        // SQLite database path (default: "./data/bridge.db")
	HTTPClientTimeout      time.Duration // HTTP client timeout for API calls (default: 15s)
	RedisURL               string        // Redis URL for caching (optional)
	CacheTTL               time.Duration // Cache TTL duration (default: 5 minutes)
	EnableRequestLogging   bool          // Enable request/response body logging (default: false, debug only)
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
	cfg.RedisURL = os.Getenv("REDIS_URL")

	// Cache TTL (in minutes)
	if ttlStr := os.Getenv("CACHE_TTL"); ttlStr != "" {
		if ttlMin, err := strconv.Atoi(ttlStr); err == nil && ttlMin > 0 {
			cfg.CacheTTL = time.Duration(ttlMin) * time.Minute
		} else {
			cfg.CacheTTL = 5 * time.Minute // Default on parse error
		}
	} else {
		cfg.CacheTTL = 5 * time.Minute // Default
	}

	// Enable request logging (for debugging only - should be disabled in production)
	cfg.EnableRequestLogging = os.Getenv("ENABLE_REQUEST_LOGGING") == "true"

	return cfg
}

// Validate checks if the configuration is valid.
// It returns an error if any required fields are missing or invalid.
func (c *Config) Validate() error {
	var errors []string

	// Viber configuration validation
	if c.APIToken == "" {
		errors = append(errors, "VIBER_API_TOKEN is required")
	}
	if c.WebhookURL == "" {
		errors = append(errors, "VIBER_WEBHOOK_URL is required")
	} else {
		if _, err := url.Parse(c.WebhookURL); err != nil {
			errors = append(errors, fmt.Sprintf("VIBER_WEBHOOK_URL is invalid: %v", err))
		}
		// In production, webhook URL should be HTTPS
		if !strings.HasPrefix(c.WebhookURL, "https://") {
			errors = append(errors, "VIBER_WEBHOOK_URL should use HTTPS in production")
		}
	}

	if c.ListenAddress == "" {
		c.ListenAddress = ":8080" // Use default if not set
	}

	// Matrix configuration validation (required if bridging)
	hasMatrixConfig := c.MatrixHomeserverURL != "" || c.MatrixAccessToken != "" || c.MatrixDefaultRoomID != ""
	if hasMatrixConfig {
		if c.MatrixHomeserverURL == "" {
			errors = append(errors, "MATRIX_HOMESERVER_URL is required when Matrix bridging is enabled")
		} else {
			if _, err := url.Parse(c.MatrixHomeserverURL); err != nil {
				errors = append(errors, fmt.Sprintf("MATRIX_HOMESERVER_URL is invalid: %v", err))
			}
		}

		if c.MatrixAccessToken == "" {
			errors = append(errors, "MATRIX_ACCESS_TOKEN is required when Matrix bridging is enabled")
		}

		if c.MatrixDefaultRoomID == "" {
			errors = append(errors, "MATRIX_DEFAULT_ROOM_ID is required when Matrix bridging is enabled")
		} else {
			// Validate Matrix room ID format using regex
			if err := utils.ValidateMatrixRoomID(c.MatrixDefaultRoomID); err != nil {
				errors = append(errors, fmt.Sprintf("MATRIX_DEFAULT_ROOM_ID has invalid format: %v", err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed:\n  %s", strings.Join(errors, "\n  "))
	}

	return nil
}
