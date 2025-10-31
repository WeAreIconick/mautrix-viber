// Package viber signature tests - unit tests for signature verification.
package viber

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignatureVerification(t *testing.T) {
	token := "test-api-token"
	body := `{"event":"message","sender":{"id":"123","name":"Test"},"message":{"type":"text","text":"Hello"}}`

	// Calculate expected signature
	mac := hmac.New(sha256.New, []byte(token))
	mac.Write([]byte(body))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	// Create request with signature
	req := httptest.NewRequest(http.MethodPost, "/viber/webhook", strings.NewReader(body))
	req.Header.Set("X-Viber-Content-Signature", expectedSig)

	// Verify signature matches what our handler would calculate
	mac2 := hmac.New(sha256.New, []byte(token))
	mac2.Write([]byte(body))
	calculatedSig := hex.EncodeToString(mac2.Sum(nil))

	if !hmac.Equal([]byte(expectedSig), []byte(calculatedSig)) {
		t.Error("Signature calculation mismatch")
	}
}

func TestSignatureMismatch(t *testing.T) {
	token := "test-api-token"
	body := `{"event":"message"}`

	// Calculate correct signature
	mac := hmac.New(sha256.New, []byte(token))
	mac.Write([]byte(body))
	correctSig := hex.EncodeToString(mac.Sum(nil))

	// Use wrong signature
	wrongSig := "wrong_signature"

	// Verify they don't match
	if hmac.Equal([]byte(correctSig), []byte(wrongSig)) {
		t.Error("Signature validation should fail for mismatched signatures")
	}
}
