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

type Config struct {
	APIToken        string
	WebhookURL      string
	ListenAddress   string
}

type Client struct {
    config     Config
    httpClient *http.Client
    matrix     *mx.Client
    db         *database.DB
}

func NewClient(cfg Config, matrixClient *mx.Client, db *database.DB) *Client {
    return &Client{config: cfg, httpClient: &http.Client{Timeout: 15 * time.Second}, matrix: matrixClient, db: db}
}

// Placeholder for setting webhook with Viber API
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

// HTTP handler for receiving Viber callbacks
func (c *Client) WebhookHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    // Read body once for signature verification and decoding
    raw, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "unable to read body", http.StatusBadRequest)
        return
    }

    // Verify Viber signature: X-Viber-Content-Signature = HMAC-SHA256(body, token)
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

    var payload WebhookRequest
    if err := json.NewDecoder(bytes.NewReader(raw)).Decode(&payload); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    metricWebhookRequests.WithLabelValues(string(payload.Event)).Inc()

    // Minimal behavior: on message event, forward text to Matrix
    if payload.Event == EventMessage && payload.Message.Type == "text" && c.matrix != nil {
        text := fmt.Sprintf("[Viber] %s: %s", payload.Sender.Name, payload.Message.Text)
        _ = c.matrix.SendText(context.Background(), text)
        metricForwardedMessages.WithLabelValues("text").Inc()
    }

    // Upsert sender into DB for mapping and group membership when chat ID is present
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
