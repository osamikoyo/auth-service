package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	once            sync.Once
)

func InitMetrics() {
	once.Do(func() {
		RequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "request_count",
			Help: "Request count in auth service",
		}, []string{"method", "endpoint", "status"})

		RequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Time to route request in seconds",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "endpoint", "status"})

		prometheus.MustRegister(RequestCount, RequestDuration)
	})
}