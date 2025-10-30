// Package circuitbreaker implements circuit breaker pattern for external API calls.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
	ErrTimeout     = errors.New("operation timeout")
)

// State represents circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements circuit breaker pattern.
type CircuitBreaker struct {
	mu              sync.RWMutex
	state           State
	failureCount    int
	successCount    int
	maxFailures     int
	maxSuccesses    int
	timeout         time.Duration
	lastFailureTime time.Time
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(maxFailures, maxSuccesses int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:        StateClosed,
		maxFailures:  maxFailures,
		maxSuccesses: maxSuccesses,
		timeout:      timeout,
	}
}

// Execute executes a function with circuit breaker protection.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	// Check if circuit should transition from open to half-open
	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = StateHalfOpen
			cb.successCount = 0
		} else {
			return ErrCircuitOpen
		}
	}
	
	// Execute function
	err := fn()
	
	if err != nil {
		cb.onFailure()
		return err
	}
	
	cb.onSuccess()
	return nil
}

// onFailure handles a failure.
func (cb *CircuitBreaker) onFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()
	
	if cb.state == StateHalfOpen {
		// Half-open -> Open on any failure
		cb.state = StateOpen
		cb.failureCount = 0
	} else if cb.failureCount >= cb.maxFailures {
		// Closed -> Open on max failures
		cb.state = StateOpen
		cb.failureCount = 0
	}
}

// onSuccess handles a success.
func (cb *CircuitBreaker) onSuccess() {
	cb.failureCount = 0
	
	if cb.state == StateHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.maxSuccesses {
			// Half-open -> Closed on max successes
			cb.state = StateClosed
			cb.successCount = 0
		}
	}
}

// GetState returns the current state.
func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failureCount = 0
	cb.successCount = 0
}

