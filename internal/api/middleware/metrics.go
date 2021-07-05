package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func MeasureResponseDuration(namespace string) func(next http.Handler) http.Handler {
	responseTimeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "http_server_request_duration_seconds",
		Help:      "Histogram of response time for handler in seconds",
	}, []string{"route", "method", "status_code"})

	if err := prometheus.Register(responseTimeHistogram); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			responseTimeHistogram = are.ExistingCollector.(*prometheus.HistogramVec)
		}
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			statusCode := strconv.Itoa(ww.Status())
			route := getRoutePattern(r)
			responseTimeHistogram.WithLabelValues(route, r.Method, statusCode).Observe(duration.Seconds())
		}

		return http.HandlerFunc(fn)
	}
}

// getRoutePattern returns the route pattern from the chi context there are 3 conditions
// a) static routes "/events" => "/events"
// b) dynamic routes "/events/:id" => "/events/{id}"
// c) if nothing matches the output is undefined
func getRoutePattern(r *http.Request) string {
	reqContext := chi.RouteContext(r.Context())
	if pattern := reqContext.RoutePattern(); pattern != "" {
		return pattern
	}

	return "undefined"
}
