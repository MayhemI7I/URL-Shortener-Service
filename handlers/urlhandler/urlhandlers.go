package urlhandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
<<<<<<< HEAD
	"github.com/golang-jwt/jwt/v4"
=======
>>>>>>> 8392c39f29fbebc6ac979bb47028df03598833e3
	"net/http"

	"strings"
	"time"

<<<<<<< HEAD
	"local/domain"
	"local/internal/storage"
	"local/logger"

=======
>>>>>>> 8392c39f29fbebc6ac979bb47028df03598833e3
	"go.uber.org/zap"
	"github.com/MayhemI7I/URL-Shortener-Service/utils/httputil"

	"github.com/MayhemI7I/URL-Shortener-Service/domain"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/storage"
	"github.com/MayhemI7I/URL-Shortener-Service/logger"
)


// contextKey is a custom type to avoid key collisions in context
type contextKey string

// UserIDKey is the key used to store user ID in the request context
const UserIDKey contextKey = "userID"

// URLGenerator defines an interface for generating short URLs
type URLGenerator interface {
	GenerateShortURL(origURL string) (string, error)
}

// URLHandler manages URL-related HTTP requests
type URLHandler struct {
	storage      storage.Storage
	urlGenerator URLGenerator
}

// NewURLHandler initializes a new URLHandler with the provided dependencies
func NewURLHandler(storage storage.Storage, urlGenerator URLGenerator) *URLHandler {
<<<<<<< HEAD
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
=======
	return &URLHandler{
		storage:      storage,
		urlGenerator: urlGenerator,
	}
>>>>>>> 8392c39f29fbebc6ac979bb47028df03598833e3
}

func(h *URLHandler)DeleteURL(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userID, err := extractUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	


// HandleGet processes GET requests to redirect from a short URL to the original URL.
// It retrieves the user ID from the context and uses it to fetch the original URL.
func (h *URLHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
<<<<<<< HEAD

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
=======

	userID, err := extractUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	logger.Log.Info("handling GET request", zap.String("shortURL", shortURL), zap.String("userID", userID))

	origURL, err := h.storage.Get(ctx, shortURL, userID)
	if err != nil {
		httputil.RespondWithError(w, err)
		return
	}

	w.Header().Set("Location", origURL)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusTemporaryRedirect)
	logger.Log.Info("redirecting to original URL", zap.String("to", origURL))
}

func(h *URLHandler)HandleGetUserAllURLs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
defer cancel()

userID, err := extractUserID(r)
if err != nil {
	http.Error(w, err.Error(), http.StatusUnauthorized)
	return
}

responseURLs, err := h.storage.GetUserAllURLs(ctx, userID)
if err != nil {
	logger.Log.Error("failed to get user URLs", zap.Error(err))
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
if len(responseURLs) == 0 {
	w.WriteHeader(http.StatusNoContent)
		return 
	
}

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
if err = json.NewEncoder(w).Encode(responseURLs); err != nil {
	logger.Log.Error("failed to encode response", zap.Error(err))
	return
}
logger.Log.Info("successfully retrieved user URLs", zap.String("userID", userID), zap.Int("count", len(responseURLs)))


>>>>>>> 8392c39f29fbebc6ac979bb47028df03598833e3

}

// HandlePost processes POST requests to create short URLs from original URLs.
// It supports both form data and JSON input, using the user ID from the context.
func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
<<<<<<< HEAD

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

=======
	defer closeBody(r)

	userID, err := extractUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	urls, contentType, err := parseRequestBody(r, userID)
	if err != nil {
		httputil.RespondWithError(w, err)
		return
	}

	responseURLs, err := h.ProcessURLs(ctx, urls, userID)
	if err != nil {
		httputil.RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(responseURLs); err != nil {
		logger.Log.Error("failed to encode response", zap.Error(err))
	}
}

// HandURL routes incoming requests to the appropriate handler based on the HTTP method
>>>>>>> 8392c39f29fbebc6ac979bb47028df03598833e3
func (h *URLHandler) HandURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.HandleGet(w, r)
	case http.MethodPost:
		h.HandlePost(w, r)
	default:
		logger.Log.Error("method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}





// Helper functions



// extractUserID retrieves the user ID from the request context
func ExtractUserID(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", errors.New("unauthorized: missing or invalid user ID")
	}
	return userID, nil
}

// parseRequestBody parses the request body based on Content-Type
func parseRequestBody(r *http.Request, userID string) ([]domain.URLData, string, error) {
	contentType := r.Header.Get("Content-Type")
	var urls []domain.URLData

	switch contentType {
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			return nil, "", fmt.Errorf("invalid form data: %v", err)
		}
		origURL := r.FormValue("url")
		if origURL == "" {
			return nil, "", errors.New("URL is required")
		}
		urls = []domain.URLData{{User: domain.User{ID: userID}, URLPair: domain.URLPair{OrigURL: origURL}},

}

	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			return nil, "", fmt.Errorf("error decoding JSON: %v", err)
		}
		for i := range urls {
			urls[i].User.ID = userID
		}

	default:
		return nil, "", errors.New("unsupported Content-Type")
	}
	return urls, contentType, nil
}

// ProcessURLs generates or retrieves short URLs for the given list
func (h *URLHandler) ProcessURLs(ctx context.Context, urls []domain.URLData, userID string) ([]domain.URLData, error) {
	var responseURLs []domain.URLData
	for _, url := range urls {
		shortURL, err := h.storage.FindByOriginalURL(ctx, url.OrigURL, userID)
		if err != nil && !errors.Is(err, domain.ErrURLNotFound) {
			logger.Log.Error("failed to check existing URL", zap.String("origURL", url.OrigURL), zap.Error(err))
			return nil, err
		}

		if shortURL != "" {
			responseURLs = append(responseURLs, domain.URLData{domain.User{ID: userID},domain.URLPair{ShortURL: shortURL, OrigURL: url.OrigURL}})
			continue
		}

		shortURL, err = h.urlGenerator.GenerateShortURL(url.OrigURL)
		if err != nil {
			logger.Log.Error("failed to generate short URL", zap.String("origURL", url.OrigURL), zap.Error(err))
			return nil, err
		}

		if err := h.storage.Save(ctx, shortURL, url.OrigURL, userID); err != nil {
			if errors.Is(err, domain.ErrURLExists) {
				return nil, fmt.Errorf("URL already exists: %v", err)
			}
			logger.Log.Error("failed to save URL", zap.String("shortURL", shortURL), zap.Error(err))
			return nil, err
		}
		responseURLs = append(responseURLs, domain.URLData{domain.User{ID: userID},domain.URLPair{ShortURL: shortURL, OrigURL: url.OrigURL}})
	}
	return responseURLs, nil
}

// closeBody ensures the request body is closed
func closeBody(r *http.Request) {
	if r.Body != nil {
		r.Body.Close()
	}
}
