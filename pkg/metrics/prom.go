package prom_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequests *prometheus.CounterVec
	HTTPDuration *prometheus.HistogramVec
)

func Init() {
	//prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	//prometheus.MustRegister(collectors.NewGoCollector())

	// Создаём метрики:
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gamehub",
			Subsystem: "http_server",
			Name:      "requests_total",
			Help:      "Количество HTTP-запросов.",
		},
		[]string{"method", "path", "status"},
	)
	HTTPDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "gamehub",
			Subsystem: "http_server",
			Name:      "request_duration_seconds",
			Help:      "Длительность HTTP-запроса.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	prometheus.MustRegister(HTTPRequests, HTTPDuration)
}
