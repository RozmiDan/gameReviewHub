package middleware_metrics

import (
	"net/http"
	"strconv"

	prom_metrics "github.com/RozmiDan/gameReviewHub/pkg/metrics"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		prom_metrics.HTTPInFlight.WithLabelValues(method, path).Inc()
        defer prom_metrics.HTTPInFlight.WithLabelValues(method, path).Dec()

		histTimer := prometheus.NewTimer(prom_metrics.HTTPDuration.WithLabelValues(method, path))
        defer func() {
            histTimer.ObserveDuration()
        }()

		// timer := prometheus.NewTimer(prom_metrics.HTTPDuration.WithLabelValues(method, path))
		// defer timer.ObserveDuration()

		next.ServeHTTP(rw, r)

		status := strconv.Itoa(rw.Status())
		prom_metrics.HTTPRequests.WithLabelValues(method, path, status).Inc()
	})
}
