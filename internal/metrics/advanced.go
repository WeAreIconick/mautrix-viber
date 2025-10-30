// Package metrics advanced provides detailed Prometheus metrics: message latency, queue depth, error rates.
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// MessageLatency tracks message processing latency.
	MessageLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "viber_message_latency_seconds",
			Help:    "Time taken to process and forward messages",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"direction", "type"}, // direction: viber_to_matrix, matrix_to_viber
	)

	// QueueDepth tracks the depth of message queues.
	QueueDepth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "viber_queue_depth",
			Help: "Current depth of message processing queues",
		},
		[]string{"queue_type"}, // queue_type: pending, failed, retrying
	)

	// ErrorRate tracks error rates by type.
	ErrorRate = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "viber_errors_total",
			Help: "Total number of errors by type",
		},
		[]string{"error_type", "source"}, // error_type: signature_failure, send_failure, decode_error, etc.
	)

	// OperationDuration tracks duration of various operations.
	OperationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "viber_operation_duration_seconds",
			Help:    "Duration of various bridge operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"}, // operation: webhook_process, matrix_send, viber_send, etc.
	)

	// ActiveConnections tracks active connections.
	ActiveConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "viber_active_connections",
			Help: "Number of active connections",
		},
		[]string{"connection_type"}, // connection_type: matrix_sync, viber_webhook
	)
)

func init() {
	prometheus.MustRegister(MessageLatency)
	prometheus.MustRegister(QueueDepth)
	prometheus.MustRegister(ErrorRate)
	prometheus.MustRegister(OperationDuration)
	prometheus.MustRegister(ActiveConnections)
}

// RecordMessageLatency records message processing latency.
func RecordMessageLatency(direction, msgType string, duration time.Duration) {
	MessageLatency.WithLabelValues(direction, msgType).Observe(duration.Seconds())
}

// RecordQueueDepth records current queue depth.
func RecordQueueDepth(queueType string, depth int) {
	QueueDepth.WithLabelValues(queueType).Set(float64(depth))
}

// RecordError records an error occurrence.
func RecordError(errorType, source string) {
	ErrorRate.WithLabelValues(errorType, source).Inc()
}

// RecordOperationDuration records operation duration.
func RecordOperationDuration(operation string, duration time.Duration) {
	OperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordActiveConnections records active connection count.
func RecordActiveConnections(connectionType string, count int) {
	ActiveConnections.WithLabelValues(connectionType).Set(float64(count))
}

