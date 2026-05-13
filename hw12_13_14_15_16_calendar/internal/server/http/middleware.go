package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	}, []string{"method", "endpoint", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "endpoint"})
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		latency := time.Since(start)
		endpoint := normalizePath(r.URL.Path)

		httpRequestsTotal.WithLabelValues(r.Method, endpoint, fmt.Sprintf("%d", rw.status)).Inc()
		httpRequestDuration.WithLabelValues(r.Method, endpoint).Observe(latency.Seconds())
	})
}

func normalizePath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return path
	}

	if parts[1] == "events" {
		if len(parts) == 3 {
			return "/events/{id}"
		}
		if len(parts) == 4 && parts[2] == "interval" {
			return "/events/interval/{interval}"
		}
	}

	return path
}

func loggingMiddleware(logger Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		latency := time.Since(start)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		logMsg := fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\" %s",
			ip,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			rw.status,
			0,
			r.UserAgent(),
			latency,
		)

		logger.Info(logMsg)
	})
}
