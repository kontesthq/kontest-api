package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		slog.Info("Request", slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.String("duration", time.Since(start).String()))
	})
}
