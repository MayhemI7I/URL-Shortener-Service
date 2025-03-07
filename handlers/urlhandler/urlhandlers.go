package urlhandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"strings"
	"time"

	"go.uber.org/zap"

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
	return &URLHandler{
		storage:      storage,
		urlGenerator: urlGenerator,
	}
}


// HandleGet processes GET requests to redirect from a short URL to the original URL.
// It retrieves the user ID from the context and uses it to fetch the original URL.
func (h *URLHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userID, err := extractUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	logger.Log.Info("handling GET request", zap.String("shortURL", shortURL), zap.String("userID", userID))

	origURL, err := h.storage.Get(ctx, shortURL, userID)
	if err != nil {
		respondWithError(w, err)
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



}

// HandlePost processes POST requests to create short URLs from original URLs.
// It supports both form data and JSON input, using the user ID from the context.
func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	defer closeBody(r)

	userID, err := extractUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	urls, contentType, err := parseRequestBody(r, userID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	responseURLs, err := h.processURLs(ctx, urls, userID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(responseURLs); err != nil {
		logger.Log.Error("failed to encode response", zap.Error(err))
	}
}

// HandURL routes incoming requests to the appropriate handler based on the HTTP method
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
func extractUserID(r *http.Request) (string, error) {
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
		urls = []domain.URLData{{URLPair: domain.URLPair{OrigURL: origURL}, UserID: userID}}

	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			return nil, "", fmt.Errorf("error decoding JSON: %v", err)
		}
		for i := range urls {
			urls[i].UserID = userID
		}

	default:
		return nil, "", errors.New("unsupported Content-Type")
	}
	return urls, contentType, nil
}

// processURLs generates or retrieves short URLs for the given list
func (h *URLHandler) processURLs(ctx context.Context, urls []domain.URLData, userID string) ([]domain.URLData, error) {
	var responseURLs []domain.URLData
	for _, url := range urls {
		shortURL, err := h.storage.FindByLongURL(ctx, url.OrigURL, userID)
		if err != nil && !errors.Is(err, domain.ErrURLNotFound) {
			logger.Log.Error("failed to check existing URL", zap.String("origURL", url.OrigURL), zap.Error(err))
			return nil, err
		}

		if shortURL != "" {
			responseURLs = append(responseURLs, domain.URLData{URLPair: domain.URLPair{ShortURL: shortURL, OrigURL: url.OrigURL}, UserID: userID})
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
		responseURLs = append(responseURLs, domain.URLData{URLPair: domain.URLPair{ShortURL: shortURL, OrigURL: url.OrigURL}, UserID: userID})
	}
	return responseURLs, nil
}

// closeBody ensures the request body is closed
func closeBody(r *http.Request) {
	if r.Body != nil {
		r.Body.Close()
	}
}

// respondWithError sends an appropriate HTTP error response based on the error type
func respondWithError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrURLNotFound):
		http.Error(w, "URL not found", http.StatusNotFound)
	case errors.Is(err, domain.ErrURLExists):
		http.Error(w, "URL already exists", http.StatusConflict)
	case errors.Is(err, domain.ErrTokenNotFound):
		http.Error(w, "Unauthorized: refresh token not found", http.StatusUnauthorized)
	case errors.Is(err, domain.ErrTokenExpired):
		http.Error(w, "Unauthorized: refresh token expired", http.StatusUnauthorized)
	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
	default:
		logger.Log.Error("internal error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
