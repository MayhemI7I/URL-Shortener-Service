package main

import (
	"local/compression/zstd"
	"local/config"
	"local/handlers/loghandler"
	"local/handlers/urlhandler"
	"local/internal/storage"
	"local/logger"
	"local/utils"
	"net/http"
	"time"
)

// initApp выполняет все необходимые иниты и возвращает готовые зависимости.
func initApp() (*config.Config, *urlhandler.URLHandler, error) {
	// Загружаем конфиг
	cfg := config.InitConfig()

	// Инициализируем логгер
	logger.InitLogger(cfg.LogLevel)

	// Инициализируем хранилище
	store, err := storage.NewStorage(*cfg)
	if err != nil {
		return nil, nil, err
	}

	// Создаем генератор коротких URL
	genUrl := utils.NewGeneratorShortURL(cfg.URLLength)

	// Создаем обработчик URL
	urlHandler := urlhandler.NewURLHandler(store, genUrl)

	return cfg, urlHandler, nil
}

func main() {
	cfg, urlHandler, err := initApp()
	if err != nil {
		logger.Log.Fatalf("failed to initialize application: %v", err)
	}
	defer logger.CloseLogger()

	// Создаем HTTP multiplexer и регистрируем хендлеры
	mux := http.NewServeMux()
	compressedHandler := loghandler.WithLog(
		zstd.Decompression(
			zstd.Compression(
				http.HandlerFunc(urlHandler.HandURL),
			),
		),
	)
	mux.Handle("/", compressedHandler)

	mux.Handle("/api/shorten", loghandler.WithLog(
		zstd.Decompression(
			zstd.Compression(
				http.HandlerFunc(urlHandler.HandlePost),
			),
		),
	))

	// Запускаем сервер
	if err := runServer(cfg, mux); err != nil {
		logger.Log.Fatalf("failed to start server: %v", err)
	}
}

// runServer запускает HTTP-сервер
func runServer(cfg *config.Config, mux *http.ServeMux) error {
	addr := cfg.ServerAdress + ":" + cfg.ServerPort
	logger.Log.Infof(time.Now().Format("2006-01-02 15:04:05")+"Server started on %s", addr)
	return http.ListenAndServe(addr, mux)
}
