package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/christophercampbell/dseq/app/logging"
	"github.com/christophercampbell/dseq/app/metrics"
	"github.com/christophercampbell/dseq/app/tracing"
	"go.opentelemetry.io/otel/attribute"
)

// Middleware represents an HTTP middleware function
type Middleware func(http.Handler) http.Handler

// Chain chains multiple middleware functions
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Logging creates a logging middleware
func Logging(logger *logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create response writer wrapper
			rw := newResponseWriter(w)

			// Process request
			next.ServeHTTP(rw, r)

			// Log request
			duration := time.Since(start)
			logger.Info("HTTP request",
				map[string]interface{}{
					"method":     r.Method,
					"path":       r.URL.Path,
					"status":     rw.status,
					"duration":   duration,
					"user_agent": r.UserAgent(),
					"remote_ip":  r.RemoteAddr,
				})
		})
	}
}

// Metrics creates a metrics middleware
func Metrics() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create response writer wrapper
			rw := newResponseWriter(w)

			// Process request
			next.ServeHTTP(rw, r)

			// Record metrics
			duration := time.Since(start).Seconds()
			metrics.TransactionDuration.WithLabelValues(r.Method).Observe(duration)
			metrics.TransactionTotal.WithLabelValues(fmt.Sprintf("%d", rw.status)).Inc()
		})
	}
}

// Tracing creates a tracing middleware
func Tracing() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracing.StartSpan(r.Context(), "http.request")
			defer span.End()

			// Add request attributes
			tracing.SetAttributes(span,
				attribute.String("http.method", r.Method),
				attribute.String("http.path", r.URL.Path),
				attribute.String("http.user_agent", r.UserAgent()),
				attribute.String("http.remote_ip", r.RemoteAddr),
			)

			// Create response writer wrapper
			rw := newResponseWriter(w)

			// Process request
			next.ServeHTTP(rw, r.WithContext(ctx))

			// Add response attributes
			tracing.SetAttributes(span,
				attribute.Int("http.status", rw.status),
			)
		})
	}
}

// Recovery creates a recovery middleware
func Recovery(logger *logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(fmt.Errorf("panic recovered: %v", err), "HTTP request panic",
						map[string]interface{}{
							"method":    r.Method,
							"path":      r.URL.Path,
							"remote_ip": r.RemoteAddr,
						})
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter is a wrapper around http.ResponseWriter
type responseWriter struct {
	http.ResponseWriter
	status int
}

// newResponseWriter creates a new responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// WriteHeader implements http.ResponseWriter
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
