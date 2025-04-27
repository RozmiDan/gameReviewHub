package middleware_main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func MyLogger(log *zap.Logger) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		log = log.With(zap.String("component", "middleware/logger"))
		log.Info("logger middleware enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			curLog := log.With(zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				curLog.Info("request completed",
					zap.Int("status", ww.Status()),
					zap.Int("bytes", ww.BytesWritten()),
					zap.Duration("request time", time.Since(t1)),
				)
			}()
			next.ServeHTTP(w, r)
		})
	}
}
