package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTimMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)
		duration := time.Since(start)
		// Log request details and response time
		fmt.Printf("Request: %s %s, Status: %d, Response Time: %s, \n", r.Method, r.URL.Path, wrappedWriter.status, duration)
	})
}

// responseWriter
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
