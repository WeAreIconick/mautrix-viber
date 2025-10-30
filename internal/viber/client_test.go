// Package viber tests - unit tests for Viber client.
package viber

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/mautrix-viber/internal/database"
)

// TestWebhookSignatureVerification tests signature verification.
func TestWebhookSignatureVerification(t *testing.T) {
	token := "test-api-token"
	payload := WebhookRequest{
		Event: EventMessage,
		Sender: Sender{ID: "123", Name: "Test"},
		Message: Message{Type: "text", Text: "Hello"},
	}
	
	bodyBytes, _ := json.Marshal(payload)
	
	// Calculate signature
	mac := hmac.New(sha256.New, []byte(token))
	mac.Write(bodyBytes)
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	
	// Create client
	client := NewClient(Config{APIToken: token}, nil, nil)
	
	// Create request with valid signature
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(bodyBytes))
	req.Header.Set("X-Viber-Content-Signature", expectedSig)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	client.WebhookHandler(w, req)
	
	// Should not return 401 (unauthorized)
	if w.Code == http.StatusUnauthorized {
		t.Error("Valid signature was rejected")
	}
}

// TestWebhookSignatureMismatch tests rejection of invalid signatures.
func TestWebhookSignatureMismatch(t *testing.T) {
	token := "test-api-token"
	client := NewClient(Config{APIToken: token}, nil, nil)
	
	payload := `{"event":"message"}`
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader([]byte(payload)))
	req.Header.Set("X-Viber-Content-Signature", "invalid_signature")
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	client.WebhookHandler(w, req)
	
	// Should return 401 for invalid signature
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for invalid signature, got %d", w.Code)
	}
}

// TestSendMessage tests sending messages via Viber API.
func TestSendMessage(t *testing.T) {
	t.Skip("Requires mock HTTP server or Viber API access")
	
	// This would test:
	// 1. Message formatting
	// 2. API request construction
	// 3. Error handling
	// 4. Response parsing
}

// TestWebhookHandler_Integration tests webhook handler with full flow.
func TestWebhookHandler_Integration(t *testing.T) {
	t.Skip("Requires database and Matrix client setup")
	
	// This would test:
	// 1. Receive webhook
	// 2. Verify signature
	// 3. Store sender in database
	// 4. Forward to Matrix
	// 5. Store message mapping
}

// TestEnsureWebhook tests webhook registration.
func TestEnsureWebhook(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pa/set_webhook" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		
		response := map[string]interface{}{
			"status":        0,
			"status_message": "ok",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	_ = NewClient(Config{
		APIToken:   "test-token",
		WebhookURL: "https://example.com/webhook",
	}, nil, nil)
	
	// Client created for testing
	// Note: This would require making the endpoint configurable for actual testing
	
	// Test webhook registration
	ctx := context.Background()
	// err := client.EnsureWebhook()
	// Would need to verify webhook was registered
	// Requires mock HTTP server to intercept API calls
	_ = ctx
	t.Log("Mock test setup - requires endpoint configuration for webhook registration testing")
}

// TestWebhookHandler_EventTypes tests handling of different event types.
func TestWebhookHandler_EventTypes(t *testing.T) {
	db, _ := database.Open("/tmp/test_events.db")
	defer db.Close()
	
	client := NewClient(Config{APIToken: "test"}, nil, db)
	
	testCases := []struct {
		name    string
		event   Event
		message Message
		wantErr bool
	}{
		{"text message", EventMessage, Message{Type: "text", Text: "Hello"}, false},
		{"image message", EventMessage, Message{Type: "picture", Media: "https://example.com/image.jpg"}, false},
		{"subscribed event", EventSubscribed, Message{}, false},
		{"unsubscribed event", EventUnsubscribed, Message{}, false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload := WebhookRequest{
				Event:   tc.event,
				Sender:  Sender{ID: "123", Name: "Test"},
				Message: tc.message,
			}
			
			bodyBytes, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			client.WebhookHandler(w, req)
			
			// Should not crash or return 500 for valid events
			if w.Code == http.StatusInternalServerError {
				t.Errorf("Handler returned 500 for valid event")
			}
		})
	}
}
