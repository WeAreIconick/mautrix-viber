package main

import (
    "net"
    "net/http"
    "sync"
    "time"
)

type tokenBucket struct {
    tokens       float64
    lastRefill   time.Time
}

type ipRateLimiter struct {
    mu      sync.Mutex
    rate    float64
    burst   float64
    buckets map[string]*tokenBucket
}

func newIPRateLimiter(ratePerSec float64, burst float64) *ipRateLimiter {
    return &ipRateLimiter{rate: ratePerSec, burst: burst, buckets: make(map[string]*tokenBucket)}
}

func (l *ipRateLimiter) Allow(ip string) bool {
    now := time.Now()
    l.mu.Lock()
    defer l.mu.Unlock()
    b, ok := l.buckets[ip]
    if !ok {
        b = &tokenBucket{tokens: l.burst, lastRefill: now}
        l.buckets[ip] = b
    }
    // refill
    elapsed := now.Sub(b.lastRefill).Seconds()
    b.tokens = min(l.burst, b.tokens+elapsed*l.rate)
    b.lastRefill = now
    if b.tokens < 1 {
        return false
    }
    b.tokens -= 1
    return true
}

func min(a, b float64) float64 { if a < b { return a }; return b }

func clientIP(r *http.Request) string {
    // If behind a reverse proxy, consider parsing X-Forwarded-For carefully.
    host, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return r.RemoteAddr
    }
    return host
}


