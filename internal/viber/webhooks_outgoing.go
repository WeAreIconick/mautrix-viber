// Package viber webhooks_outgoing provides outgoing webhooks for Matrix events (for external integrations).
package viber

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OutgoingWebhook represents an outgoing webhook configuration.
type OutgoingWebhook struct {
	URL     string
	Secret  string
	Events  []string // Event types to forward
	Enabled bool
}

// OutgoingWebhookManager manages outgoing webhooks.
type OutgoingWebhookManager struct {
	webhooks []OutgoingWebhook
	client   *http.Client
}

// NewOutgoingWebhookManager creates a new outgoing webhook manager.
func NewOutgoingWebhookManager() *OutgoingWebhookManager {
	return &OutgoingWebhookManager{
		webhooks: []OutgoingWebhook{},
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

// RegisterWebhook registers an outgoing webhook.
func (owm *OutgoingWebhookManager) RegisterWebhook(webhook OutgoingWebhook) {
	owm.webhooks = append(owm.webhooks, webhook)
}

// SendWebhook sends a webhook event to registered webhooks.
func (owm *OutgoingWebhookManager) SendWebhook(ctx context.Context, eventType string, payload interface{}) error {
	for _, webhook := range owm.webhooks {
		if !webhook.Enabled {
			continue
		}
		
		// Check if this webhook should receive this event type
		shouldSend := false
		for _, evt := range webhook.Events {
			if evt == eventType || evt == "*" {
				shouldSend = true
				break
			}
		}
		
		if !shouldSend {
			continue
		}
		
		// Send webhook
		if err := owm.sendWebhook(ctx, webhook, eventType, payload); err != nil {
			// Log error but continue with other webhooks
			// In production, use structured logging:
			// logger.Warn("failed to send outgoing webhook", "url", webhook.URL, "error", err)
		}
	}
	
	return nil
}

// sendWebhook sends a single webhook.
func (owm *OutgoingWebhookManager) sendWebhook(ctx context.Context, webhook OutgoingWebhook, eventType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook.URL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Event", eventType)
	
	// Add signature if secret is configured
	if webhook.Secret != "" {
		// TODO: Add HMAC signature
	}
	
	resp, err := owm.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	
	return nil
}

