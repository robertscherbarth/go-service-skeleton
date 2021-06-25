package middleware

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RequestLogger(logger *zap.Logger, serviceName string) func(next http.Handler) http.Handler {
	fields := []zap.Field{
		zap.String("application_type", "service"),
		zap.String("log_type", "access"),
	}
	childLogger := logger.With(fields...)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			duration := time.Since(start).Seconds()
			fields := []zap.Field{
				zap.String("service", serviceName),
				zap.String("remote_address", r.RemoteAddr),
				zap.String("protocol", r.Proto),
				zap.String("request_method", r.Method),
				zap.String("query_string", r.URL.RawQuery),
				zap.Int64("bytes_received", r.ContentLength),
				zap.Int("status", ww.Status()),
				zap.String("uri", r.URL.Path),
				zap.String("response_time", fmt.Sprintf("%f", duration)),
				zap.Int("bytes_sent", ww.BytesWritten()),
				zap.String("remote_client_id", r.Header.Get("remoteClientId")),
				zap.String("user_agent", r.Header.Get("User-Agent")),
			}
			childLogger.Info("-", fields...)
		}
		return http.HandlerFunc(fn)
	}
}
