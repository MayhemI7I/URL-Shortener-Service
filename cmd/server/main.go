package main

import (
	"local/compression/zstd"
	"local/config"
	"local/handlers/urlhandler"
   "local/handlers/pinghandler"
   "local/internal/db"
   "local/handlers/loghandler"
	"local/internal/urlstorage"
	"local/logger"
	"net/http"
)

func main() {
	// Initialize configuration
	cfg := config.InitConfig()

	// Initialize logger
	logger.InitLogger(cfg.LogLevel)
	defer logger.CloseLogger()

	// Initialize URL storage
	store, err := urlstorage.NewURLStorage(cfg.FileStorage)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer store.Close()

	// Initialize database connector
	db, err := db.NewDBConnector(cfg.DataBaseDSN)
	if err != nil {
		logger.Log.Fatal(err)
	}
   defer db.CloseDataBase()

	// Create ping and URL handlers
	pingHandler := pinghandler.NewPingHandler(db)
	urlHandler := urlhandler.NewURLHandler(store)

	// Create HTTP multiplexer and register handlers
	mux := http.NewServeMux()
	mux.Handle("/", loghandler.WithLog(zstd.ZstdDecompress(zstd.ZstdCompress(http.HandlerFunc(urlHandler.HandURL)))))
	mux.Handle("/ping", pingHandler)

	// Start the server
	if err := run(cfg, mux); err != nil {
		logger.Log.Fatal(err)
	}
}

func run(cfg *config.Config, mux *http.ServeMux) error {
	logger.Log.Infof("Server started on %s:%s", cfg.ServerAdress, cfg.ServerPort)
	return http.ListenAndServe(cfg.ServerAdress+":"+cfg.ServerPort, mux)
}
