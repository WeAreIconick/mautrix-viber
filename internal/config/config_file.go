// Package config file provides YAML configuration file support
// in addition to environment variable configuration.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
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
func (c Config) Validate() error {
	if c.APIToken == "" {
		return fmt.Errorf("viber api token is required")
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("viber webhook url is required")
	}
	if c.ListenAddress == "" {
		return fmt.Errorf("listen address is required")
	}
	// Matrix config is optional but if provided, should be complete
	if c.MatrixHomeserverURL != "" || c.MatrixAccessToken != "" || c.MatrixDefaultRoomID != "" {
		if c.MatrixHomeserverURL == "" {
			return fmt.Errorf("matrix homeserver url is required if matrix is configured")
		}
		if c.MatrixAccessToken == "" {
			return fmt.Errorf("matrix access token is required if matrix is configured")
		}
		if c.MatrixDefaultRoomID == "" {
			return fmt.Errorf("matrix default room id is required if matrix is configured")
		}
	}
	return nil
}

