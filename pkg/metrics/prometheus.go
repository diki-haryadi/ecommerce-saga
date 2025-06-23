package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Database metrics
	DbConnectionsOpen = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_open",
			Help: "Number of open database connections",
		},
	)

	DbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query", "type"},
	)

	// Cache metrics
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache"},
	)

	// Message broker metrics
	MessagePublished = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_published_total",
			Help: "Total number of messages published",
		},
		[]string{"broker", "topic"},
	)

	MessageConsumed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_consumed_total",
			Help: "Total number of messages consumed",
		},
		[]string{"broker", "topic"},
	)

	// Business metrics
	OrdersCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_created_total",
			Help: "Total number of orders created",
		},
	)

	PaymentsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payments_processed_total",
			Help: "Total number of payments processed",
		},
		[]string{"status"},
	)

	// System metrics
	GoroutinesCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_count",
			Help: "Number of goroutines",
		},
	)

	MemoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
	)
)
