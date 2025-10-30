// Package config file provides YAML configuration file support
// in addition to environment variable configuration.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/example/mautrix-viber/internal/utils"
)

// FileConfig represents configuration from a YAML file.
type FileConfig struct {
	Viber struct {
		APIToken   string `yaml:"api_token"`
		WebhookURL string `yaml:"webhook_url"`
	} `yaml:"viber"`
	
	Matrix struct {
		HomeserverURL string `yaml:"homeserver_url"`
		AccessToken   string `yaml:"access_token"`
		DefaultRoomID string `yaml:"default_room_id"`
	} `yaml:"matrix"`
	
	Server struct {
		ListenAddress string `yaml:"listen_address"`
	} `yaml:"server"`
	
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
}

// LoadFromFile loads configuration from a YAML file.
// Environment variables override file values.
func LoadFromFile(path string) (Config, error) {
	var fileConfig FileConfig
	
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}
	
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return Config{}, fmt.Errorf("parse config file: %w", err)
	}
	
	cfg := Config{
		APIToken:            fileConfig.Viber.APIToken,
		WebhookURL:          fileConfig.Viber.WebhookURL,
		ListenAddress:       fileConfig.Server.ListenAddress,
		MatrixHomeserverURL: fileConfig.Matrix.HomeserverURL,
		MatrixAccessToken:   fileConfig.Matrix.AccessToken,
		MatrixDefaultRoomID: fileConfig.Matrix.DefaultRoomID,
	}
	
	// Override with environment variables if present
	envCfg := FromEnv()
	if envCfg.APIToken != "" {
		cfg.APIToken = envCfg.APIToken
	}
	if envCfg.WebhookURL != "" {
		cfg.WebhookURL = envCfg.WebhookURL
	}
	if envCfg.ListenAddress != "" {
		cfg.ListenAddress = envCfg.ListenAddress
	}
	if envCfg.MatrixHomeserverURL != "" {
		cfg.MatrixHomeserverURL = envCfg.MatrixHomeserverURL
	}
	if envCfg.MatrixAccessToken != "" {
		cfg.MatrixAccessToken = envCfg.MatrixAccessToken
	}
	if envCfg.MatrixDefaultRoomID != "" {
		cfg.MatrixDefaultRoomID = envCfg.MatrixDefaultRoomID
	}
	
	return cfg, nil
}

// Validate checks configuration for required fields and valid values.
// This is an enhanced version that also validates URLs and enforces HTTPS for production.
func (c Config) Validate() error {
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
		// Validate Matrix homeserver URL format
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


