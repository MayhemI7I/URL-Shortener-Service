package urlhandler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"local/logger"

	"go.uber.org/zap"
)

// URLRequest представляет запрос на URL.
type URLRequest struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

// URLStorage — интерфейс для хранения URL.
type URLStorage interface {
	Get(ctx context.Context, shortURL string) (string, error)
	Save(ctx context.Context, shortURL, longURL string) error
}

// URLGenerator — интерфейс для генерации коротких URL.
type URLGenerator interface {
	GenerateShortURL(longURL string) (string, error)
}

// URLHandler — обработчик для работы с URL.
type URLHandler struct {
	storage      URLStorage
	urlGenerator URLGenerator
}

// NewURLHandler создает новый URLHandler.
func NewURLHandler(storage URLStorage, urlGenerator URLGenerator) *URLHandler {
	return &URLHandler{storage: storage, urlGenerator: urlGenerator}
}

// HandleGet обрабатывает GET-запрос.
func (h *URLHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	shortURL := r.URL.Path[1:]
	longURL, err := h.storage.Get(ctx, shortURL)
	if err != nil {
		if err == context.DeadlineExceeded {
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		} else {
			logger.Log.Error("URL not found", zap.Error(err))
			http.Error(w, "URL not found", http.StatusNotFound)
		}
		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	logger.Log.Debug("redirection", zap.String("to", longURL))
}

// HandlePost обрабатывает POST-запрос.
func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	longURL := string(body)
	shortURL, err := h.urlGenerator.GenerateShortURL(longURL)
	if err != nil {
		http.Error(w, "Error generating short URL", http.StatusInternalServerError)
		return
	}

	err = h.storage.Save(ctx, shortURL, longURL)
	if err != nil {
		http.Error(w, "Error saving URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

// HandJsonPost обрабатывает JSON POST-запрос.
func (h *URLHandler) HandJsonPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req URLRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	shortURL, err := h.urlGenerator.GenerateShortURL(req.LongURL)
	if err != nil {
		http.Error(w, "Error generating short URL", http.StatusInternalServerError)
		return
	}

	err = h.storage.Save(ctx, shortURL, req.LongURL)
	if err != nil {
		http.Error(w, "Error saving URL", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"result": shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
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
