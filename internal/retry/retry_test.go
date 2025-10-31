// Package retry tests - unit tests for retry logic.
package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo_Success(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxAttempts = 3

	var attempts int
	err := Do(context.Background(), cfg, func() error {
		attempts++
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestDo_Retry(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxAttempts = 3
	cfg.InitialDelay = 10 * time.Millisecond

	var attempts int
	err := Do(context.Background(), cfg, func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error after retry, got %v", err)
	}

	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestDo_MaxAttempts(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxAttempts = 3
	cfg.InitialDelay = 10 * time.Millisecond

	var attempts int
	err := Do(context.Background(), cfg, func() error {
		attempts++
		return errors.New("persistent error")
	})

	if err == nil {
		t.Error("Expected error after max attempts")
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestDo_ContextCancel(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxAttempts = 10
	cfg.InitialDelay = 100 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := Do(ctx, cfg, func() error {
		return errors.New("error")
	})

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}
