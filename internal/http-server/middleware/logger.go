package middleware

import (
	"log/slog"
	"net/http"
	"time"
	"url-shortener/internal/lib/consts"
)

func NewLoggingMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String(consts.RequestIdKey, r.Context().Value(consts.RequestIdKey).(string)),
			)
			wrw := &wrappedResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", wrw.statusCode),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(wrw, r)
		})
	}
}
