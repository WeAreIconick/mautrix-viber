package config

import (
	"os"
)

type Config struct {
	APIToken      string
	WebhookURL    string
	ListenAddress string
}

func FromEnv() Config {
	cfg := Config{}
	cfg.APIToken = os.Getenv("VIBER_API_TOKEN")
	cfg.WebhookURL = os.Getenv("VIBER_WEBHOOK_URL")
	cfg.ListenAddress = os.Getenv("LISTEN_ADDRESS")
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = ":8080"
	}
	return cfg
}
