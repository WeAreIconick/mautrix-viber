// Package main advanced rate limiting: per-user, per-room, adaptive limits.
package main

import (
	"net/http"
	"sync"
	"time"
)

// AdvancedRateLimiter provides per-user, per-room rate limiting with adaptive limits.
type AdvancedRateLimiter struct {
	mu            sync.RWMutex
	userLimiters  map[string]*AdaptiveLimiter
	roomLimiters  map[string]*AdaptiveLimiter
	globalLimiter *AdaptiveLimiter
	baseRate      float64
	burstSize     float64
}

// AdaptiveLimiter is a token bucket with adaptive rate adjustment.
type AdaptiveLimiter struct {
	mu         sync.Mutex
	tokens     float64
	rate       float64
	burst      float64
	lastRefill time.Time
	requests   int
	errors     int
}

// NewAdvancedRateLimiter creates a new advanced rate limiter.
func NewAdvancedRateLimiter(baseRate, burstSize float64) *AdvancedRateLimiter {
	return &AdvancedRateLimiter{
		userLimiters:  make(map[string]*AdaptiveLimiter),
		roomLimiters:  make(map[string]*AdaptiveLimiter),
		globalLimiter: NewAdaptiveLimiter(baseRate, burstSize),
		baseRate:      baseRate,
		burstSize:     burstSize,
	}
}

// NewAdaptiveLimiter creates a new adaptive limiter.
func NewAdaptiveLimiter(rate, burst float64) *AdaptiveLimiter {
	return &AdaptiveLimiter{
		tokens:     burst,
		rate:       rate,
		burst:      burst,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed.
func (al *AdaptiveLimiter) Allow() bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(al.lastRefill).Seconds()

	// Refill tokens
	al.tokens = min(al.burst, al.tokens+elapsed*al.rate)
	al.lastRefill = now

	// Adaptive rate adjustment based on error rate
	if al.requests > 10 {
		errorRate := float64(al.errors) / float64(al.requests)
		if errorRate > 0.1 {
			// High error rate - reduce rate by 20%
			al.rate = al.rate * 0.8
		} else if errorRate < 0.01 {
			// Low error rate - increase rate by 10%
			al.rate = min(al.burst/10, al.rate*1.1)
		}
		al.requests = 0
		al.errors = 0
	}

	if al.tokens < 1 {
		return false
	}

	al.tokens--
	al.requests++
	return true
}

// RecordError records an error for adaptive rate adjustment.
func (al *AdaptiveLimiter) RecordError() {
	al.mu.Lock()
	defer al.mu.Unlock()
	al.errors++
}

// AllowUser checks if a request is allowed for a specific user.
func (arl *AdvancedRateLimiter) AllowUser(userID string) bool {
	arl.mu.RLock()
	limiter, exists := arl.userLimiters[userID]
	arl.mu.RUnlock()

	if !exists {
		arl.mu.Lock()
		limiter = NewAdaptiveLimiter(arl.baseRate, arl.burstSize)
		arl.userLimiters[userID] = limiter
		arl.mu.Unlock()
	}

	// Check both global and user limits
	return arl.globalLimiter.Allow() && limiter.Allow()
}

// AllowRoom checks if a request is allowed for a specific room.
func (arl *AdvancedRateLimiter) AllowRoom(roomID string) bool {
	arl.mu.RLock()
	limiter, exists := arl.roomLimiters[roomID]
	arl.mu.RUnlock()

	if !exists {
		arl.mu.Lock()
		limiter = NewAdaptiveLimiter(arl.baseRate*2, arl.burstSize*2) // Rooms get higher limits
		arl.roomLimiters[roomID] = limiter
		arl.mu.Unlock()
	}

	return arl.globalLimiter.Allow() && limiter.Allow()
}

// withAdvancedRateLimit applies advanced rate limiting middleware.
// Not currently used - reserved for future per-user/per-room rate limiting
//
//nolint:deadcode,unused
func withAdvancedRateLimit(next http.Handler) http.Handler {
	limiter := NewAdvancedRateLimiter(10, 20) // 10 req/sec, burst 20

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID and room ID from request if available
		userID := r.Header.Get("X-User-ID")
		roomID := r.URL.Query().Get("room_id")

		allowed := false
		if userID != "" {
			allowed = limiter.AllowUser(userID)
		} else if roomID != "" {
			allowed = limiter.AllowRoom(roomID)
		} else {
			allowed = limiter.globalLimiter.Allow()
		}

		if !allowed {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
