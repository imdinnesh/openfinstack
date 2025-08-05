package observability

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Time taken to process HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)

	HTTPRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(HTTPRequestDuration, HTTPRequestCount)
}
