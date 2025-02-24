package pinghandler

import (
	"local/internal/db"
	"local/internal/storage"
	"local/logger"
	"net/http"

	"go.uber.org/zap"
)

type PingHandler struct {
   db storage.Storage
}

func NewPingHandler(db storage.Storage) *PingHandler {
   return &PingHandler{
   	db: storage.Storage
   }
}

func (p *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   if r.Method != http.MethodGet {
   	logger.Log.Warn("Method not allowed", zap.String("method", r.Method))
   	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
   	return
   }
   if err := p.dbConnector.PingDataBase(); err != nil {
   	logger.Log.Error("Failed to ping database", zap.Error(err))
   	http.Error(w, err.Error(), http.StatusInternalServerError)
   	return
   }
   logger.Log.Info("Ping successful")
   w.WriteHeader(http.StatusOK)
   w.Write([]byte("OK"))
}
