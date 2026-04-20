package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// contextKey is a custom type for context keys
type contextKey string

const requestIDKey contextKey = "requestID"

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// GetRequestID gets request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// RequestID adds a unique request ID to each request
func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateID()
		}
		ctx := WithRequestID(r.Context(), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func generateID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

// Logger is a middleware that logs requests (placeholder)
func Logger(next http.Handler) http.Handler {
	return next
}

// Recoverer recovers from panics
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				panicErr, ok := err.(error)
				if !ok {
					panicErr = fmt.Errorf("%v", err)
				}
				log.Printf("panic: %v", panicErr)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Internal Server Error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Compress is a placeholder for GZIP compression
func Compress(next http.Handler, level string) http.Handler {
	return next
}

// Heartbeat returns a simple heartbeat response
func Heartbeat(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// CORS adds CORS headers
func CORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// NoCache adds no-cache headers
func NoCache(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// RequireJSON ensures the request has JSON content type
func RequireJSON(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if contentType != "" && contentType != "application/json" {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Content-Type must be application/json",
				})
				return
			}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// ResponseRecorder wraps http.ResponseWriter to capture status code
type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Written    bool
}

// NewResponseRecorder creates a new ResponseRecorder
func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

// WriteHeader captures the status code
func (rec *ResponseRecorder) WriteHeader(code int) {
	rec.StatusCode = code
	rec.Written = true
	rec.ResponseWriter.WriteHeader(code)
}

// Middleware chaining helper
type MiddlewareFunc func(http.Handler) http.Handler

// Chain combines multiple middleware into one
func Chain(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
