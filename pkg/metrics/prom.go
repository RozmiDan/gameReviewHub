package prom_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequests        *prometheus.CounterVec
	HTTPDuration        *prometheus.HistogramVec
	HTTPInFlight        *prometheus.GaugeVec
	DBErrors            *prometheus.CounterVec
	KafkaPublishErrors  *prometheus.CounterVec
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
			Help:      "Распределение длительности HTTP запросов",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	HTTPInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "gamehub",
			Subsystem: "http_server",
			Name:      "inflight_requests",
			Help:      "Текущее число обрабатываемых HTTP запросов",
		},
		[]string{"method", "path"},
	)
	DBErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gamehub",
			Subsystem: "repository",
			Name:      "db_errors_total",
			Help:      "Число ошибок при запросах в БД",
		},
		[]string{"operation"},
	)
	KafkaPublishErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gamehub",
			Subsystem: "kafka_producer",
			Name:      "publish_errors_total",
			Help:      "Число ошибок при публикации в Kafka",
		},
		[]string{"topic"},
	)

	prometheus.MustRegister(
		HTTPRequests, HTTPDuration, HTTPInFlight, DBErrors, KafkaPublishErrors,
	)
}
