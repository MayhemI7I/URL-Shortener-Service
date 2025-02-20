package urlhandler

import (
   "encoding/json"
   "io"
   "net/http"

   "local/utils"
   "local/logger"

   "go.uber.org/zap"
)

// URLRequest represents a request for a URL.
type URLRequest struct {
   ShortURL string `json:"short_url"`
   LongURL  string `json:"long_url"`
}

// URLStorage is an interface for a URL storage.
type URLStorage interface {
   GetURL(shortURL string) (string, error)
   SaveURL(shortURL, longURL string) error
}

// URLHandler is a handler for working with URLs.
type URLHandler struct {
   storage URLStorage
}

// NewURLHandler creates a new URLHandler.
func NewURLHandler(storage URLStorage) *URLHandler {
   return &URLHandler{storage: storage}
}

// getLongURL gets a long URL from the storage.
func (h *URLHandler) getLongURL(shortURL string) (string, error) {
   return h.storage.GetURL(shortURL)
}

// HandleGet handles a GET request.
func (h *URLHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
   shortURL := r.URL.Path[1:]
   longURL, err := h.getLongURL(shortURL)
   if err != nil {
   	logger.Log.Error("URL not found", zap.Error(err))
   	http.Error(w, "URL not found", http.StatusNotFound)
   	return
   }

   w.Header().Set("Location", longURL)
   w.WriteHeader(http.StatusTemporaryRedirect)
   logger.Log.Debug("redirection", zap.String("to", longURL))
}

// HandlePost handles a POST request.
func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
   body, err := io.ReadAll(r.Body)
   if err != nil || len(body) == 0 {
   	logger.Log.Error("Invalid request body", zap.Error(err))
   	http.Error(w, "Invalid request body", http.StatusBadRequest)
   	return
   }

   longURL := string(body)
   if longURL == "" {
   	http.Error(w, "Invalid URL", http.StatusBadRequest)
   	return
   }

   shortURL, err := utils.GenerateShortURL(longURL)
   if err != nil {
   	logger.Log.Error("Error generating short URL", zap.Error(err))
   	http.Error(w, "Error generating short URL", http.StatusInternalServerError)
   	return
   }

   err = h.storage.SaveURL(shortURL, longURL)
   if err != nil {
   	logger.Log.Error("Error saving URL", zap.Error(err))
   	http.Error(w, "Error saving URL", http.StatusInternalServerError)
   	return
   }

   w.Header().Set("Content-Type", "text/plain; charset=utf-8")
   w.WriteHeader(http.StatusCreated)
   w.Write([]byte(shortURL))
}

// HandJsonPost handles a JSON POST request.
func (h *URLHandler) HandJsonPost(w http.ResponseWriter, r *http.Request) {
   if r.Method != http.MethodPost {
   	logger.Log.Error("Method not allowed", zap.Int("status", http.StatusMethodNotAllowed))
   	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
   	return
   }
   logger.Log.Debug("decoding request")
   var req URLRequest
   dec := json.NewDecoder(r.Body)
   if err := dec.Decode(&req); err != nil {
   	logger.Log.Error("cannot decode request JSON body", zap.Error(err))
   	http.Error(w, "Error decoding request", http.StatusBadRequest)
   	return
   }
   shortURL, err := utils.GenerateShortURL(req.LongURL)
   if err != nil {
   	logger.Log.Error("cannot generate short URL", zap.Error(err))
   	http.Error(w, "Error generating short URL", http.StatusInternalServerError)
   	return
   }
   err = h.storage.SaveURL(shortURL, req.LongURL)
   if err != nil {
   	logger.Log.Error("Error saving URL", zap.Error(err))
   	http.Error(w, "Error saving URL", http.StatusInternalServerError)
   	return
   }
   response := map[string]string{"result": shortURL}
   w.Header().Set("Content-Type", "application/json")
   w.WriteHeader(http.StatusCreated)
   json.NewEncoder(w).Encode(response)
}

// HandURL handles a URL request.
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
