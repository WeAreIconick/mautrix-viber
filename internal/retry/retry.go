// Package retry provides exponential backoff retry logic with jitter
// for handling transient failures in external API calls.
package retry

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Config configures retry behavior.
type Config struct {
	MaxAttempts  int           // Maximum number of retry attempts
	InitialDelay time.Duration // Initial delay before first retry
	MaxDelay     time.Duration // Maximum delay between retries
	Multiplier   float64       // Exponential backoff multiplier
	Jitter       bool          // Add random jitter to delays
}

// DefaultConfig returns a default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
	}
}

// Do executes a function with retry logic. Returns the last error if all attempts fail.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if attempt > 0 {
			// Calculate delay with exponential backoff
			delay := time.Duration(float64(cfg.InitialDelay) * pow(cfg.Multiplier, float64(attempt-1)))
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
			if cfg.Jitter {
				// Add ±25% jitter
				jitter := time.Duration(float64(delay) * 0.25 * (rand.Float64()*2 - 1))
				delay += jitter
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		if err := fn(); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return fmt.Errorf("max attempts (%d) exceeded: %w", cfg.MaxAttempts, lastErr)
}

func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}

// IsRetryable determines if an error is worth retrying.
// Common retryable errors: network timeouts, 5xx errors, rate limits.
// Currently returns true for all errors; can be enhanced to check specific error types.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	// Could check for specific error types:
	// - Network timeouts: errors.Is(err, context.DeadlineExceeded)
	// - HTTP 5xx: check if err contains status code >= 500
	// - Rate limits: check if err contains 429 status code
	return true
}
