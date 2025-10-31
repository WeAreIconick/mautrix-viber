// Package circuitbreaker tests - unit tests for circuit breaker.
package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_ClosedState(t *testing.T) {
	cb := NewCircuitBreaker(3, 2, 1*time.Second)

	// Should allow execution when closed
	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if cb.GetState() != StateClosed {
		t.Errorf("Expected state Closed, got %v", cb.GetState())
	}
}

func TestCircuitBreaker_OpenAfterFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, 2, 100*time.Millisecond)

	// Fail 3 times
	for i := 0; i < 3; i++ {
		_ = cb.Execute(func() error {
			return errors.New("test error")
		})
	}

	// Should be open now
	if cb.GetState() != StateOpen {
		t.Errorf("Expected state Open, got %v", cb.GetState())
	}

	// Should reject execution
	err := cb.Execute(func() error {
		return nil
	})

	if err != ErrCircuitOpen {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_HalfOpenRecovery(t *testing.T) {
	cb := NewCircuitBreaker(2, 2, 50*time.Millisecond)

	// Cause it to open
	for i := 0; i < 2; i++ {
		_ = cb.Execute(func() error {
			return errors.New("test error")
		})
	}

	// Wait for timeout
	time.Sleep(100 * time.Millisecond)

	// Execute should transition to half-open state when called
	err := cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// First success in half-open should keep it in half-open (need 2 successes)
	if cb.GetState() != StateHalfOpen {
		t.Errorf("Expected state HalfOpen after first success, got %v", cb.GetState())
	}

	// Second success should close it
	err = cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// After second success, should transition to Closed
	if cb.GetState() != StateClosed {
		t.Errorf("Expected state Closed, got %v", cb.GetState())
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	cb := NewCircuitBreaker(2, 2, 100*time.Millisecond)

	// Open it
	for i := 0; i < 2; i++ {
		_ = cb.Execute(func() error {
			return errors.New("test error")
		})
	}

	// Reset
	cb.Reset()

	if cb.GetState() != StateClosed {
		t.Errorf("Expected state Closed after reset, got %v", cb.GetState())
	}
}
