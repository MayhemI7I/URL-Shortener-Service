package urlhandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"

	"local/domain"
	"local/internal/storage"
	"local/logger"

	"go.uber.org/zap"
)

// URLGenerator — интерфейс для генерации коротких URL.
type URLGenerator interface {
	GenerateShortURL(origURL string) (string, error)
}

// URLHandler — обработчик для работы с URL.
type URLHandler struct {
	storage      storage.Storage
	urlGenerator URLGenerator
}

// NewURLHandler создает новый URLHandler.
func NewURLHandler(storage storage.Storage, urlGenerator URLGenerator) *URLHandler {
	return &URLHandler{storage: storage, urlGenerator: urlGenerator}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
		token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("your-secret-key"), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// HandleGet обрабатывает GET-запрос.
func (h *URLHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	userID := strings.TrimPrefix(r.URL.Path, "/")
	logger.Log.Info("shortURL", zap.String("shortURL", shortURL))
	origUrl, err := h.storage.Get(ctx, shortURL, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		} else {
			logger.Log.Error("URL not found", zap.Error(err))
			http.Error(w, "URL not found", http.StatusNotFound)
		}
		return
	}

	w.Header().Set("Location", origUrl)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusTemporaryRedirect)
	logger.Log.Info("redirection", zap.String("to", origUrl))

}
func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	defer func() {
		if r.Body != nil {
			r.Body.Close()
		}
	}()

	var origUrl string
	requestURLs := make([]domain.URLData, 0)
	responseURLs := make([]domain.URLData, 0)

	contentType := r.Header.Get("Content-Type")

	// Обработка FormData
	if contentType == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		origUrl = r.FormValue("url")
		if origUrl == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}
		requestURLs = append(requestURLs, domain.URLData{URLPair: domain.URLPair{OrigURL: origUrl}})

		// Обработка JSON
	} else if contentType == "application/json" {
		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&requestURLs); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	// Создание сокращенных URL для каждого из запросов
	for _, url := range requestURLs {
		shortURL, err := h.storage.FindByLongURL(ctx, url.OrigURL, url.UserID)
		if err != nil && !errors.Is(err, domain.ErrURLNotFound) {
			http.Error(w, "Error checking for existing short URL", http.StatusInternalServerError)
			return
		}

		// Если короткий URL уже существует, добавляем его в ответ
		if shortURL != "" {
			responseURLs = append(responseURLs, domain.URLData{URLPair: domain.URLPair{ShortURL: shortURL, OrigURL: url.OrigURL}})
		} else {
			shortURL, err = h.urlGenerator.GenerateShortURL(url.OrigURL)
			if err != nil && !errors.Is(err, domain.ErrURLNotFound) {
				logger.Log.Error("Ошибка после функции", zap.Error(err), zap.String(url.OrigURL, shortURL))
			}

			responseURLs = append(responseURLs, domain.URLData{URLPair: domain.URLPair{ShortURL: shortURL, OrigURL: url.OrigURL}})

			// Сохранение нового URL в базу данных
			err = h.storage.Save(ctx, shortURL, url.OrigURL, url.UserID)
			if err != nil {
				http.Error(w, "Error saving URL", http.StatusInternalServerError)
				return
			}
		}
	}

	// Ответ клиенту
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(responseURLs)
	if err != nil {
		logger.Log.Error("Error encoding JSON", zap.Error(err))
	}
}

func (h *URLHandler) HandURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.HandleGet(w, r)
	case http.MethodPost:
		h.HandlePost(w, r)
	default:
		logger.Log.Error("Method not allowed", zap.Int("status", http.StatusMethodNotAllowed))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}
