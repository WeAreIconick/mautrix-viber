// Package config tests - unit tests for configuration.
package config

import (
	"os"
	"testing"
)

func TestFromEnv(t *testing.T) {
	// Set test environment variables
	os.Setenv("VIBER_API_TOKEN", "test-token")
	os.Setenv("VIBER_WEBHOOK_URL", "https://test.com/webhook")
	os.Setenv("LISTEN_ADDRESS", ":9090")
	os.Setenv("MATRIX_HOMESERVER_URL", "https://matrix.test.com")
	os.Setenv("MATRIX_ACCESS_TOKEN", "test-access-token")
	os.Setenv("MATRIX_DEFAULT_ROOM_ID", "!test:test.com")

	defer func() {
		os.Unsetenv("VIBER_API_TOKEN")
		os.Unsetenv("VIBER_WEBHOOK_URL")
		os.Unsetenv("LISTEN_ADDRESS")
		os.Unsetenv("MATRIX_HOMESERVER_URL")
		os.Unsetenv("MATRIX_ACCESS_TOKEN")
		os.Unsetenv("MATRIX_DEFAULT_ROOM_ID")
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
	os.Unsetenv("LISTEN_ADDRESS")

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
