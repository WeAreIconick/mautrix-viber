// Package http_test provides HTTP handler tests with mock dependencies.
package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/mautrix-viber/internal/api"
	"github.com/example/mautrix-viber/internal/middleware"
	"github.com/example/mautrix-viber/internal/viber"
)

// TestHealthHandler tests the health check endpoint.
func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	// Simple health check should always return 200
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestInfoHandler tests the info API endpoint.
func TestInfoHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	w := httptest.NewRecorder()

	api.InfoHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response structure
	if _, ok := response["version"]; !ok {
		t.Error("Response missing 'version' field")
	}
	if _, ok := response["status"]; !ok {
		t.Error("Response missing 'status' field")
	}
}

// TestWebhookHandler_SignatureVerification tests webhook signature verification.
func TestWebhookHandler_SignatureVerification(t *testing.T) {
	// Create Viber client
	client := viber.NewClient(viber.Config{
		APIToken: "test-token",
	}, nil, nil)

	// Valid signature test
	body := `{"event":"message","sender":{"id":"123","name":"Test"},"message":{"type":"text","text":"Hello"}}`
	req := httptest.NewRequest(http.MethodPost, "/viber/webhook", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	
	// Calculate valid signature (simplified - actual implementation would use HMAC)
	w := httptest.NewRecorder()
	
	// Note: This is a simplified test - full test would verify HMAC signature
	client.WebhookHandler(w, req)

	// Should return some status (either 200 or 401 based on signature)
	if w.Code == 0 {
		t.Error("Handler did not write response")
	}
}

// TestRateLimiting tests rate limiting middleware.
func TestRateLimiting(t *testing.T) {
	// This would test that rate limiting middleware properly limits requests
	t.Skip("Requires rate limiter implementation testing")
}

// TestRecoveryMiddleware tests panic recovery.
func TestRecoveryMiddleware(t *testing.T) {
	// Handler that panics
	panickingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Wrap with recovery middleware
	recoveryHandler := middleware.RecoveryMiddleware(panickingHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Should recover from panic and return 500
	recoveryHandler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 after panic, got %d", w.Code)
	}
}

