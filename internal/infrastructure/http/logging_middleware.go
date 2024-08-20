package http

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests.
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom ResponseWriter to capture the status code
			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call the next handler
			next.ServeHTTP(ww, r)

			// Log the request details
			logger.Info("Handled request",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.Int("status", ww.statusCode),
				slog.String("duration", time.Since(start).String()),
				slog.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

// responseWriter is a wrapper around http.ResponseWriter that allows us to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
