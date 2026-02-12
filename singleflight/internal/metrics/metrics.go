package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// HTTP Requests Metrics
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	// HTTP Latency Metrics
	HTTPLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	// Cache Hits Metrics
	CacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total cache hits",
		},
	)

	// Cache Misses Metrics
	CacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total cache misses",
		},
	)

	// DB Hits Metrics
	DBHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_hits_total",
			Help: "Total database hits",
		},
	)

	// Singleflight Shared Metrics
	SingleflightShared = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "singleflight_shared_total",
			Help: "Total shared singleflight responses",
		},
	)

	// In-Flight Requests Metrics
	InFlightRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_in_flight_requests",
			Help: "Number of in-flight HTTP requests",
		},
	)
)

func InitMetrics() {
	prometheus.MustRegister(
		HTTPRequests,
		HTTPLatency,
		CacheHits,
		CacheMisses,
		DBHits,
		SingleflightShared,
		InFlightRequests,
	)
}
