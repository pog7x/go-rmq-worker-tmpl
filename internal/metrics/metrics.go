package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const NameSpace = "ns"

var (
	buckets = []float64{1, 3, 5, 10, 15, 20, 30, 50, 70, 100, 150, 200, 300, 400, 500, 750, 1000, 1400, 2000, 3000, 5000, 7000, 9000, 11000, 13000, 15000, 20000}

	ResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: NameSpace,
			Name:      "http_response_time_milliseconds",
			Help:      "Histogram of RT for HTTP requests (ms).",
			Buckets:   buckets,
		},
		[]string{"code", "handler", "method"},
	)

	ResponseTimeExternal = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: NameSpace,
			Name:      "external_http_response_time_milliseconds",
			Help:      "Histogram of RT for outgoing HTTP requests (ms).",
			Buckets:   buckets,
		},
		[]string{"service", "code", "handler", "method"},
	)

	HandlerDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: NameSpace,
			Name:      "handler_duration_ms",
			Help:      "Histogram of duration (ms) of handlers.",
			Buckets:   buckets,
		},
		[]string{"handler", "err", "retries"},
	)

	NetworkErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: NameSpace,
			Name:      "network_error_counter",
			Help:      "Counter of network errors.",
		},
		[]string{"service", "handler", "type"},
	)
)

func init() {
	// Metrics have to be registered to be exposed.
	prometheus.MustRegister(ResponseTimeExternal, ResponseTime, HandlerDuration, NetworkErrorCounter)
}
