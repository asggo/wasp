package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

// Logger creates a logging middleware using slog.
func Logger() func(next http.Handler) http.Handler {
	l := httplog.NewLogger(
		"webapp-log",
		httplog.Options{
			LogLevel: slog.LevelDebug,
		},
	)

	return httplog.RequestLogger(l)
}

// Timeout adds a timeout to the request context on each request.
func Timeout(secs int) func(next http.Handler) http.Handler {
	dur := time.Duration(secs) * time.Second

	return middleware.Timeout(dur)
}

// SecurityHeaders sets our security headers on the response.
func SecurityHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
