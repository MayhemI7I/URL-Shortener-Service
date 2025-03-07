package main

import (
	"github.com/MayhemI7I/URL-Shortener-Service/compression/zstd"
	"github.com/MayhemI7I/URL-Shortener-Service/config"
	"github.com/MayhemI7I/URL-Shortener-Service/handlers/middleware"
	"github.com/MayhemI7I/URL-Shortener-Service/handlers/urlhandler"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/storage"
	"github.com/MayhemI7I/URL-Shortener-Service/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/utils"
	"net/http"
	"time"
)

// initApp выполняет все необходимые иниты и возвращает готовые зависимости.
func initApp() (*config.Config, *urlhandler.URLHandler, storage.Storage, error) {
	// Загружаем конфиг
	cfg := config.InitConfig()

	// Инициализируем логгер
	logger.InitLogger(cfg.LogLevel)

	// Инициализируем хранилище
	store, err := storage.NewStorage(*cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	// Создаем генератор коротких URL
	genUrl := utils.NewGeneratorShortURL(cfg.URLLength)

	// Создаем обработчик URL
	urlHandler := urlhandler.NewURLHandler(store, genUrl)

	return cfg, urlHandler,store, nil
}

func main() {
	cfg, urlHandler,store, err := initApp()
    if err != nil {
        logger.Log.Fatalf("failed to initialize application: %v", err)
    }
    defer logger.CloseLogger()

    // Создаем HTTP multiplexer
    mux := http.NewServeMux()

    // Унифицированный конвейер для всех маршрутов
    mux.Handle("/", middleware.Conveyor(
        http.HandlerFunc(urlHandler.HandURL),
        middleware.WithLog,
        zstd.Decompression,
        zstd.Compression,
    ))

    mux.Handle("/api/shorten", middleware.Conveyor(
        http.HandlerFunc(urlHandler.HandlePost),
        middleware.WithLog,
        zstd.Decompression,
        zstd.Compression,
    ))

    mux.Handle("/api/users/urls", middleware.Conveyor(
        http.HandlerFunc(urlHandler.HandleGetUserAllURLs),
        middleware.WithLog,
        middleware.Auth(store),
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
