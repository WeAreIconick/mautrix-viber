package config

import (
	"os"
)

type Config struct {
	APIToken      string
	WebhookURL    string
	ListenAddress string
    MatrixHomeserverURL string
    MatrixAccessToken   string
    MatrixDefaultRoomID string
    ViberDefaultReceiverID string
}

func FromEnv() Config {
	cfg := Config{}
	cfg.APIToken = os.Getenv("VIBER_API_TOKEN")
	cfg.WebhookURL = os.Getenv("VIBER_WEBHOOK_URL")
	cfg.ListenAddress = os.Getenv("LISTEN_ADDRESS")
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = ":8080"
	}
    cfg.MatrixHomeserverURL = os.Getenv("MATRIX_HOMESERVER_URL")
    cfg.MatrixAccessToken = os.Getenv("MATRIX_ACCESS_TOKEN")
    cfg.MatrixDefaultRoomID = os.Getenv("MATRIX_DEFAULT_ROOM_ID")
    cfg.ViberDefaultReceiverID = os.Getenv("VIBER_DEFAULT_RECEIVER_ID")
	return cfg
}
