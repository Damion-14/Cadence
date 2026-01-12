package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type LogEntry struct {
	Timestamp string  `json:"timestamp"`
	RequestID string  `json:"request_id"`
	Method    string  `json:"method"`
	Path      string  `json:"path"`
	Status    int     `json:"status"`
	Duration  float64 `json:"duration_ms"`
	UserID    int     `json:"user_id,omitempty"`
	Error     string  `json:"error,omitempty"`
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			status:         200,
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Milliseconds()

		entry := LogEntry{
			Timestamp: start.Format(time.RFC3339),
			RequestID: GetRequestID(r.Context()),
			Method:    r.Method,
			Path:      r.URL.Path,
			Status:    wrapped.status,
			Duration:  float64(duration),
		}

		if userID := GetUserID(r.Context()); userID != 0 {
			entry.UserID = userID
		}

		json.NewEncoder(os.Stdout).Encode(entry)
	})
}
