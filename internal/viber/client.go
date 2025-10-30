package viber

import (
	"net/http"
)

type Config struct {
	APIToken        string
	WebhookURL      string
	ListenAddress   string
}

type Client struct {
	config Config
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{config: cfg, httpClient: &http.Client{}}
}

// Placeholder for setting webhook with Viber API
func (c *Client) EnsureWebhook() error {
	// TODO: call Viber set_webhook with c.config.WebhookURL
	return nil
}

// HTTP handler for receiving Viber callbacks
func (c *Client) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
