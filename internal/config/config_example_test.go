package config_test

import (
	"fmt"
	"os"

	"github.com/example/mautrix-viber/internal/config"
)

func ExampleFromEnv() {
	// Set required environment variables
	os.Setenv("VIBER_API_TOKEN", "test_token")
	os.Setenv("VIBER_WEBHOOK_URL", "https://example.com/webhook")

	// Load configuration from environment
	cfg := config.FromEnv()

	fmt.Printf("API Token configured: %v\n", cfg.APIToken != "")
	fmt.Printf("Webhook URL configured: %v\n", cfg.WebhookURL != "")
	fmt.Printf("Default listen address: %s\n", cfg.ListenAddress)

	// Cleanup
	os.Unsetenv("VIBER_API_TOKEN")
	os.Unsetenv("VIBER_WEBHOOK_URL")

	// Output:
	// API Token configured: true
	// Webhook URL configured: true
	// Default listen address: :8080
}

func ExampleConfig_Validate() {
	// Create a valid configuration
	cfg := &config.Config{
		APIToken:   "valid_token",
		WebhookURL: "https://example.com/webhook",
	}

	// Validate the configuration
	err := cfg.Validate()
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Configuration is valid")
	}

	// Test with invalid configuration
	cfg.WebhookURL = "not-a-url"
	err = cfg.Validate()
	if err != nil {
		fmt.Println("Invalid URL correctly rejected")
	}

	// Output:
	// Configuration is valid
	// Invalid URL correctly rejected
}

func ExampleConfig_Validate_withMatrix() {
	// Create a configuration with Matrix bridging
	cfg := &config.Config{
		APIToken:           "valid_token",
		WebhookURL:         "https://example.com/webhook",
		MatrixHomeserverURL: "https://matrix.example.com",
		MatrixAccessToken:   "matrix_token",
		MatrixDefaultRoomID: "!abc:matrix.example.com",
	}

	// Validate the configuration
	err := cfg.Validate()
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Matrix configuration is valid")
	}

	// Output:
	// Matrix configuration is valid
}

