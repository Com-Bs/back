package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// RequestLog entry database
type RequestLog struct {
	ID        int64
	Method    string
	Path      string
	Status    int
	Duration  time.Duration
	IP        string
	UserAgent string
	Headers   string
	Body      string
	CreatedAt time.Time
}

// Custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// DBLoggingMiddleware returns a middleware that logs HTTP requests to a database
func DBLoggingMiddleware(db string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Pre-request operations
			start := time.Now()
			var bodyBytes []byte
			if r.Body != nil {
				bodyBytes, _ = json.Marshal(r.Body)
			}

			// Create a custom response writer to capture the status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Execute the next handler (the actual request)
			next.ServeHTTP(rw, r)

			// Post-request operations
			logEntry := RequestLog{
				Method:    r.Method,
				Path:      r.URL.Path,
				Status:    rw.statusCode,
				Duration:  time.Since(start),
				IP:        r.RemoteAddr,
				UserAgent: r.UserAgent(),
				Body:      string(bodyBytes),
				CreatedAt: time.Now(),
			}

			log.Printf("Request completed: %v", logEntry)
		})
	}
}

// AuthenticateMiddleware returns a middleware that handles authentication
func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-request: Check authentication
		// TODO: Implement actual authentication logic

		// Execute the next handler
		next.ServeHTTP(w, r)

		// Post-request: Clean up or additional operations
		// TODO: Add any post-authentication operations
	})
}
