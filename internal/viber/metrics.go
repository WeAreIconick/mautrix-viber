package viber

import "github.com/prometheus/client_golang/prometheus"

var (
    metricWebhookRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{Name: "viber_webhook_requests_total", Help: "Viber webhook requests"},
        []string{"event"},
    )
    metricForwardedMessages = prometheus.NewCounterVec(
        prometheus.CounterOpts{Name: "viber_messages_forwarded_total", Help: "Messages forwarded to Matrix"},
        []string{"type"},
    )
    metricSignatureFailures = prometheus.NewCounter(
        prometheus.CounterOpts{Name: "viber_signature_failures_total", Help: "Signature verification failures"},
    )
)

func init() {
    prometheus.MustRegister(metricWebhookRequests)
    prometheus.MustRegister(metricForwardedMessages)
    prometheus.MustRegister(metricSignatureFailures)
}


