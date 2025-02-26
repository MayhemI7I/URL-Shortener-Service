package urlhandler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"local/internal/storage/postgres"
	"local/logger"

	"go.uber.org/zap"
)

// URLRequest представляет запрос на URL.
type URLRequest struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func NewURLRequest(longURL string) *URLRequest {

	return &URLRequest{
		ShortURL: "",
		LongURL: longURL,
	}
 
}


func(r *URLRequest)Encode()([]byte, error){
	return json.Marshal(r)
}


// URLStorage — интерфейс для хранения URL.
type URLStorage interface {
	Get(ctx context.Context, shortURL string) (string, error)
	Save(ctx context.Context, shortURL, longURL string) error
	Close() error
	IfExistUrl(ctx context.Context, shortURL string) (string, error)
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

	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	logger.Log.Info("shortURL", zap.String("shortURL", shortURL))
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
	logger.Log.Info("redirection", zap.String("to", longURL))
}


// HandJsonPost обрабатывает JSON POST-запрос.
func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var longURL string
	// Обрабатываем FormData
	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		longURL = r.FormValue("url")
	} else {
		// Обрабатываем JSON
		var req URLRequest
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		longURL = req.LongURL
	}

	// Проверяем, существует ли уже короткий URL для этого длинного
	shortURL, err := h.storage.IfExistUrl(ctx, longURL)
	if err != nil && err != postgres.ErrURLNotFound {
		http.Error(w, "Error checking for existing short URL", http.StatusInternalServerError)
		return
	}

	// Если короткий URL уже существует, возвращаем его
	if shortURL != "" {
		if r.Header.Get("Content-Type") == "application/json" {
			// Для JSON возвращаем ответ в формате JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(URLRequest{ShortURL: shortURL, LongURL: longURL})
		} else {
			// Для FormData возвращаем просто текст
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(shortURL))
		}
		return
	}

	// Генерируем новый короткий URL
	shortURL, err = h.urlGenerator.GenerateShortURL(longURL)
	if err != nil {
		http.Error(w, "Error generating short URL", http.StatusInternalServerError)
		return
	}

	// Сохраняем новый короткий URL в базе данных
	err = h.storage.Save(ctx, shortURL, longURL)
	if err != nil {
		http.Error(w, "Error saving URL", http.StatusInternalServerError)
		return
	}

	// Возвращаем короткий URL
	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(URLRequest{ShortURL: shortURL, LongURL: longURL})
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
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
