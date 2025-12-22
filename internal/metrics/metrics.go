package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestsTotal counts total requests by service and method
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pulumicost_requests_total",
			Help: "Total number of requests by service and method",
		},
		[]string{"service", "method"},
	)

	// RequestDuration tracks request latencies
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pulumicost_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method"},
	)

	// ErrorsTotal counts errors by service and method
	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pulumicost_errors_total",
			Help: "Total number of errors by service and method",
		},
		[]string{"service", "method", "error_type"},
	)

	// CostQueriesTotal tracks cost analysis queries
	CostQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pulumicost_cost_queries_total",
			Help: "Total cost queries by type",
		},
		[]string{"query_type"},
	)

	// ResourcesAnalyzed tracks number of resources analyzed
	ResourcesAnalyzed = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "pulumicost_resources_analyzed",
			Help:    "Number of resources analyzed per query",
			Buckets: []float64{1, 10, 50, 100, 500, 1000, 5000},
		},
	)

	// PluginCallsTotal tracks plugin invocations
	PluginCallsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pulumicost_plugin_calls_total",
			Help: "Total plugin calls by plugin name and status",
		},
		[]string{"plugin", "status"},
	)

	// PluginLatency tracks plugin call latencies
	PluginLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pulumicost_plugin_latency_seconds",
			Help:    "Plugin call latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"plugin"},
	)
)

// RecordRequest records a request with duration
func RecordRequest(service, method string, duration time.Duration) {
	RequestsTotal.WithLabelValues(service, method).Inc()
	RequestDuration.WithLabelValues(service, method).Observe(duration.Seconds())
}

// RecordError records an error occurrence
func RecordError(service, method, errorType string) {
	ErrorsTotal.WithLabelValues(service, method, errorType).Inc()
}

// RecordCostQuery records a cost query
func RecordCostQuery(queryType string) {
	CostQueriesTotal.WithLabelValues(queryType).Inc()
}

// RecordResourceCount records the number of resources analyzed
func RecordResourceCount(count int) {
	ResourcesAnalyzed.Observe(float64(count))
}

// RecordPluginCall records a plugin call
func RecordPluginCall(plugin, status string, duration time.Duration) {
	PluginCallsTotal.WithLabelValues(plugin, status).Inc()
	PluginLatency.WithLabelValues(plugin).Observe(duration.Seconds())
}

// Handler returns the Prometheus metrics HTTP handler
func Handler() http.Handler {
	return promhttp.Handler()
}
