package middleware

import (
	"time"
	"net/http"
	"github.com/MayhemI7I/URL-Shortener-Service/logger"
)

func WithLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		logger.Log.Debugw("request",
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
		)
	})
}
