package urlhandler

import (
   "context"
   "encoding/json"
   "errors"
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
   OrigURL  string `json:"orig_url"`
}

func NewURLRequest(origURL string) *URLRequest {

   return &URLRequest{
   	ShortURL: "",
   	OrigURL:  origURL,
   }

}


// URLStorage — интерфейс для хранения URL.
type URLStorage interface {
   Get(ctx context.Context, shortURL string) (string, error)
   Save(ctx context.Context, shortURL, origURL string) error
   Close() error
   FindByLongURL(ctx context.Context, shortURL string) (string, error)
}

// URLGenerator — интерфейс для генерации коротких URL.
type URLGenerator interface {
   GenerateShortURL(origURL string) (string, error)
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
   origUrl, err := h.storage.Get(ctx, shortURL)
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
   requestURLs := make([]URLRequest, 0)
   responseURLs := make([]URLRequest, 0)

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
   	requestURLs = append(requestURLs, URLRequest{OrigURL: origUrl})

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
   	shortURL, err := h.storage.FindByLongURL(ctx, url.OrigURL)
   	if err != nil && !errors.Is(err, postgres.ErrURLNotFound) {
   		http.Error(w, "Error checking for existing short URL", http.StatusInternalServerError)
   		return
   	}

   	// Если короткий URL уже существует, добавляем его в ответ
   	if shortURL != "" {
   		responseURLs = append(responseURLs, URLRequest{ShortURL: shortURL, OrigURL: url.OrigURL})
   	} else {
   		shortURL, err = h.urlGenerator.GenerateShortURL(url.OrigURL)
   		if err != nil && !errors.Is(err, postgres.ErrURLNotFound) {
   			logger.Log.Error("Ошибка после функции", zap.Error(err), zap.String(url.OrigURL, shortURL))
   		}

   		responseURLs = append(responseURLs, URLRequest{ShortURL: shortURL, OrigURL: url.OrigURL})

   		// Сохранение нового URL в базу данных
   		err = h.storage.Save(ctx, shortURL, url.OrigURL)
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

