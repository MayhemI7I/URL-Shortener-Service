package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/usecases"
)

// Conveyor создает цепочку middleware
func Conveyor(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// WithLog добавляет логирование запросов
func WithLog(logger infrastructure.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Создаем response writer для отслеживания статуса
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Вызываем следующий обработчик
			next.ServeHTTP(rw, r)

			// Логируем информацию о запросе
			logger.Info("HTTP запрос",
				infrastructure.NewField("method", r.Method),
				infrastructure.NewField("path", r.URL.Path),
				infrastructure.NewField("status", rw.status),
				infrastructure.NewField("duration", time.Since(start).String()),
				infrastructure.NewField("remote_addr", r.RemoteAddr),
			)
		})
	}
}

// Auth добавляет проверку аутентификации
func Auth(core *usecases.Core) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из заголовка
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Проверяем токен
			user, err := core.AuthService().ValidateToken(token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Добавляем пользователя в контекст
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// responseWriter обертка для http.ResponseWriter для отслеживания статуса
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
