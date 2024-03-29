package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"url-shortener/internal/lib/consts"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		newCtx := context.WithValue(r.Context(), consts.RequestIdKey, requestID)
		r = r.WithContext(newCtx)
		next.ServeHTTP(w, r)
	})
}
