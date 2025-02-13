package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"local/logger"
	"local/utils"
)

// Логирование запросов
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

type URLRequest struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

// Интерфейс для хранилища
type URLStorage interface {
	GetURL(shortURL string) (string, error)
	SaveURL(shortURL, longURL string) error
}

// Обработчик для работы с URL
type URLHandler struct {
	storage URLStorage
}

func NewURLHandler(storage URLStorage) *URLHandler {
	return &URLHandler{storage: storage}
}

// Получение длинного URL
func (h *URLHandler) getLongURL(shortURL string) (string, error) {
	return h.storage.GetURL(shortURL)
}

// Обработка GET-запроса
func (h *URLHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	longURL, err := h.getLongURL(shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	logger.Log.Debug("redirection", "to", longURL)
}

// Обработка POST-запроса
func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Log.Debug("Invalid request body", err)
		return
	}

	longURL := string(body)
	if longURL == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	shortURL, err := utils.GenerateShortURL(longURL)
	if err != nil {
		http.Error(w, "Error generating short URL", http.StatusInternalServerError)
		return
	}

	h.storage.SaveURL(shortURL, longURL)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *URLHandler) HandJsonPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Error("Method not allowed", http.StatusMethodNotAllowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	logger.Log.Debugw("decoding request")
	var req URLRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", err)
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}
	shortURL, err := utils.GenerateShortURL(req.LongURL)
	if err != nil {
		logger.Log.Debug("cannot generate short URL", err)
		http.Error(w, "Error generating short URL", http.StatusInternalServerError)
		return
	}
	h.storage.SaveURL(shortURL, req.LongURL)
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
