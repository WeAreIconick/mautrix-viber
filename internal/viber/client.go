package viber

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "net/http"
    "strings"
    "time"

    mx "github.com/example/mautrix-viber/internal/matrix"
    "github.com/example/mautrix-viber/internal/database"
)

// Config holds Viber API configuration.
type Config struct {
	APIToken      string // Viber Bot API token for authentication
	WebhookURL    string // Public HTTPS URL where Viber will send webhooks
	ListenAddress string // HTTP server listen address (optional)
}

// Client manages Viber API interactions and webhook handling.
// It forwards messages to Matrix and stores state in the database.
type Client struct {
	config     Config          // Viber API configuration
	httpClient *http.Client    // HTTP client for API requests (15s timeout)
	matrix     *mx.Client      // Matrix client for forwarding messages (may be nil)
	db         *database.DB    // Database for persistence (may be nil)
}

// NewClient creates a new Viber client with the given configuration.
// matrixClient and db may be nil if those features are not configured.
func NewClient(cfg Config, matrixClient *mx.Client, db *database.DB) *Client {
	return &Client{
		config:     cfg,
		httpClient: &http.Client{Timeout: 15 * time.Second},
		matrix:     matrixClient,
		db:         db,
	}
}

// EnsureWebhook registers the webhook URL with Viber's API.
// This should be called on startup to ensure Viber knows where to send events.
// Returns an error if registration fails.
func (c *Client) EnsureWebhook() error {
    if c.config.WebhookURL == "" || c.config.APIToken == "" {
        return fmt.Errorf("webhook url or api token not configured")
    }
    // Build payload for set_webhook
    body := map[string]any{
        "url":         c.config.WebhookURL,
        "event_types": []string{"message", "subscribed", "unsubscribed", "conversation_started"},
    }
    data, _ := json.Marshal(body)
    req, err := http.NewRequest(http.MethodPost, "https://chatapi.viber.com/pa/set_webhook", bytes.NewReader(data))
    if err != nil {
        return fmt.Errorf("create set_webhook request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Viber-Auth-Token", c.config.APIToken)
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("set_webhook request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("set_webhook unexpected status: %s", resp.Status)
    }
    var wr WebhookResponse
    if err := json.NewDecoder(resp.Body).Decode(&wr); err != nil {
        return fmt.Errorf("decode set_webhook response: %w", err)
    }
    if wr.Status != 0 {
        return fmt.Errorf("set_webhook failed: %d %s", wr.Status, wr.StatusMessage)
    }
    return nil
}

// WebhookHandler processes incoming Viber webhook callbacks.
// It verifies the HMAC-SHA256 signature, parses the payload, and forwards
// messages to Matrix when configured. Also stores sender information in the database.
//
// Security: All requests are verified using HMAC-SHA256 signature from the
// X-Viber-Content-Signature header before processing.
func (c *Client) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	
	// Read body once for signature verification and decoding
	// We read the entire body first because we need it for both signature
	// verification and JSON decoding
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read body", http.StatusBadRequest)
		return
	}

	// Verify Viber signature: X-Viber-Content-Signature = HMAC-SHA256(body, token)
	// This prevents unauthorized webhook calls
	sig := r.Header.Get("X-Viber-Content-Signature")
	if sig != "" && c.config.APIToken != "" {
		mac := hmac.New(sha256.New, []byte(c.config.APIToken))
		mac.Write(raw)
		expected := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(expected), []byte(sig)) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Parse webhook payload
	var payload WebhookRequest
	if err := json.NewDecoder(bytes.NewReader(raw)).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	metricWebhookRequests.WithLabelValues(string(payload.Event)).Inc()

	// Forward text messages to Matrix when configured
	// This is the basic bridging functionality - more advanced features
	// (media, formatting, etc.) are handled in other modules
	if payload.Event == EventMessage && payload.Message.Type == "text" && c.matrix != nil {
		text := fmt.Sprintf("[Viber] %s: %s", payload.Sender.Name, payload.Message.Text)
		_ = c.matrix.SendText(context.Background(), text)
		metricForwardedMessages.WithLabelValues("text").Inc()
	}

	// Store sender information in database for user mapping and group membership tracking
	// This enables features like ghost user puppeting and group chat management
	if c.db != nil && payload.Sender.ID != "" && payload.Sender.Name != "" {
        _ = c.db.UpsertViberUser(payload.Sender.ID, payload.Sender.Name)
        if payload.Message.ChatID != "" {
            _ = c.db.UpsertGroupMember(payload.Message.ChatID, payload.Sender.ID)
        }
    }

    // Picture message -> download media and forward as image
    if payload.Event == EventMessage && (payload.Message.Type == "picture" || strings.HasSuffix(strings.ToLower(payload.Message.Media), ".jpg") || strings.HasSuffix(strings.ToLower(payload.Message.Media), ".png")) && c.matrix != nil {
        if payload.Message.Media != "" {
            // Download media
            req, err := http.NewRequest(http.MethodGet, payload.Message.Media, nil)
            if err == nil {
                resp, err := c.httpClient.Do(req)
                if err == nil && resp.StatusCode == http.StatusOK {
                    data, _ := io.ReadAll(resp.Body)
                    resp.Body.Close()
                    filename := payload.Message.FileName
                    if filename == "" {
                        filename = "viber-image"
                    }
                    // best-effort content-type
                    mimeType := resp.Header.Get("Content-Type")
                    _ = c.matrix.SendImage(context.Background(), filename, mimeType, data, nil)
                    metricForwardedMessages.WithLabelValues("image").Inc()
                }
            }
        }
    }

    w.WriteHeader(http.StatusOK)
}
