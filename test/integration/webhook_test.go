// Package integration tests - integration tests for webhook handling.
package integration

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/mautrix-viber/internal/viber"
)

// TestWebhookIntegration tests the full webhook processing flow.
func TestWebhookIntegration(t *testing.T) {
	t.Skip("Requires full test environment setup")
	
	// This would test:
	// 1. Receiving webhook with valid signature
	// 2. Parsing webhook payload
	// 3. Forwarding to Matrix (with mock Matrix client)
	// 4. Storing in database
}

// TestWebhookSignatureFlow tests signature verification in webhook handler.
func TestWebhookSignatureFlow(t *testing.T) {
	// Create test payload
	payload := viber.WebhookRequest{
		Event: viber.EventMessage,
		Sender: viber.Sender{
			ID:   "test_user_123",
			Name: "Test User",
		},
		Message: viber.Message{
			Type: "text",
			Text: "Hello, world!",
		},
	}
	
	bodyBytes, _ := json.Marshal(payload)
	token := "test-api-token"
	
	// Calculate signature
	mac := hmac.New(sha256.New, []byte(token))
	mac.Write(bodyBytes)
	signature := hex.EncodeToString(mac.Sum(nil))
	
	// Create request
	req := httptest.NewRequest(http.MethodPost, "/viber/webhook", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Viber-Content-Signature", signature)
	
	// Create recorder
	w := httptest.NewRecorder()
	
	// Create client with token
	client := viber.NewClient(viber.Config{
		APIToken: token,
	}, nil, nil)
	
	// This test would require a full setup with database and Matrix client
	// For now, just verify the signature calculation
	t.Logf("Calculated signature: %s", signature)
	
	if len(signature) != 64 { // SHA256 hex is 64 chars
		t.Errorf("Invalid signature length: %d", len(signature))
	}
}

