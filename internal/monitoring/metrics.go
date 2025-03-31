package monitoring

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "geo_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"endpoint"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "geo_request_duration_seconds",
			Help:    "Duration of HTTP requests in second",
			Buckets: prometheus.LinearBuckets(0.01, 0.05, 20),
		},
		[]string{"endpoint"},
	)

	CacheRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "geo_cache_duration_seconds",
			Help:    "Duration of cache operations in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
		[]string{"method"},
	)

	ExternalAPIRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "geo_external_api_duration_seconds",
			Help:    "Duration of external API operations in seconds",
			Buckets: prometheus.ExponentialBuckets(0.01, 2, 15),
		},
		[]string{"method"},
	)
)

func RegisterMetricsRoute(r chi.Router) {
	r.Handle("/metrics", promhttp.Handler())
}

func RegisterMetrics() {
	prometheus.MustRegister(
		HTTPRequestTotal,
		HTTPRequestDuration,
		CacheRequestDuration,
		ExternalAPIRequestDuration,
	)
}
