// Package api provides REST API endpoints for bridge management and information.
package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// BridgeInfo represents bridge status and statistics.
type BridgeInfo struct {
	Version    string        `json:"version"`
	Status     string        `json:"status"`
	Uptime     string        `json:"uptime"`
	StartedAt  time.Time     `json:"started_at"`
	Matrix     ServiceStatus `json:"matrix"`
	Viber      ServiceStatus `json:"viber"`
	Statistics Statistics    `json:"statistics"`
}

// ServiceStatus represents the status of a service (Matrix or Viber).
type ServiceStatus struct {
	Connected bool   `json:"connected"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
}

// Statistics represents bridge statistics.
type Statistics struct {
	MessagesBridged int64 `json:"messages_bridged"`
	UsersLinked     int64 `json:"users_linked"`
	RoomsMapped     int64 `json:"rooms_mapped"`
	WebhookRequests int64 `json:"webhook_requests"`
	Errors          int64 `json:"errors"`
}

var (
	bridgeInfo BridgeInfo
	startTime  time.Time
)

func init() {
	startTime = time.Now()
	bridgeInfo = BridgeInfo{
		Version:   "0.1.0",
		Status:    "running",
		StartedAt: startTime,
		Matrix:    ServiceStatus{Connected: false, Status: "unknown"},
		Viber:     ServiceStatus{Connected: false, Status: "unknown"},
	}
}

// UpdateMatrixStatus updates Matrix service status.
func UpdateMatrixStatus(connected bool, status, errMsg string) {
	bridgeInfo.Matrix = ServiceStatus{
		Connected: connected,
		Status:    status,
		Error:     errMsg,
	}
}

// UpdateViberStatus updates Viber service status.
func UpdateViberStatus(connected bool, status, errMsg string) {
	bridgeInfo.Viber = ServiceStatus{
		Connected: connected,
		Status:    status,
		Error:     errMsg,
	}
}

// UpdateStatistics updates bridge statistics.
func UpdateStatistics(stats Statistics) {
	bridgeInfo.Statistics = stats
}

// InfoHandler returns bridge information and status.
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	bridgeInfo.Uptime = time.Since(startTime).String()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(bridgeInfo)
}

// HealthHandler returns health check status.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if bridgeInfo.Status != "running" {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
