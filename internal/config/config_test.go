// Package config tests - unit tests for configuration.
package config

import (
	"os"
	"testing"
)

func TestFromEnv(t *testing.T) {
	// Set test environment variables
	_ = os.Setenv("VIBER_API_TOKEN", "test-token")
	_ = os.Setenv("VIBER_WEBHOOK_URL", "https://test.com/webhook")
	_ = os.Setenv("LISTEN_ADDRESS", ":9090")
	_ = os.Setenv("MATRIX_HOMESERVER_URL", "https://matrix.test.com")
	_ = os.Setenv("MATRIX_ACCESS_TOKEN", "test-access-token")
	_ = os.Setenv("MATRIX_DEFAULT_ROOM_ID", "!test:test.com")

	defer func() {
		_ = os.Unsetenv("VIBER_API_TOKEN")
		_ = os.Unsetenv("VIBER_WEBHOOK_URL")
		_ = os.Unsetenv("LISTEN_ADDRESS")
		_ = os.Unsetenv("MATRIX_HOMESERVER_URL")
		_ = os.Unsetenv("MATRIX_ACCESS_TOKEN")
		_ = os.Unsetenv("MATRIX_DEFAULT_ROOM_ID")
	}()

	cfg := FromEnv()

	if cfg.APIToken != "test-token" {
		t.Errorf("Expected APIToken 'test-token', got '%s'", cfg.APIToken)
	}

	if cfg.WebhookURL != "https://test.com/webhook" {
		t.Errorf("Expected WebhookURL 'https://test.com/webhook', got '%s'", cfg.WebhookURL)
	}

	if cfg.ListenAddress != ":9090" {
		t.Errorf("Expected ListenAddress ':9090', got '%s'", cfg.ListenAddress)
	}
}

func TestFromEnvDefaults(t *testing.T) {
	// Clear environment
	_ = os.Unsetenv("LISTEN_ADDRESS")

	cfg := FromEnv()

	// Should default to :8080
	if cfg.ListenAddress != ":8080" {
		t.Errorf("Expected default ListenAddress ':8080', got '%s'", cfg.ListenAddress)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				APIToken:      "token",
				WebhookURL:    "https://test.com/webhook",
				ListenAddress: ":8080",
			},
			wantErr: false,
		},
		{
			name: "missing api token",
			config: Config{
				WebhookURL:    "https://test.com/webhook",
				ListenAddress: ":8080",
			},
			wantErr: true,
		},
		{
			name: "missing webhook url",
			config: Config{
				APIToken:      "token",
				ListenAddress: ":8080",
			},
			wantErr: true,
		},
		{
			name: "partial matrix config",
			config: Config{
				APIToken:            "token",
				WebhookURL:          "https://test.com/webhook",
				ListenAddress:       ":8080",
				MatrixHomeserverURL: "https://matrix.test.com",
				// Missing access token and room ID
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
