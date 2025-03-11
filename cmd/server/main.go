package main

import (
   "context"
   "net/http"
   "os"
   "os/signal"
   "strconv"
   "syscall"
   "time"

   "github.com/MayhemI7I/URL-Shortener-Service/compression/zstd"
   "github.com/MayhemI7I/URL-Shortener-Service/config"
   "github.com/MayhemI7I/URL-Shortener-Service/handlers/authhandler"
   "github.com/MayhemI7I/URL-Shortener-Service/handlers/middleware"
   "github.com/MayhemI7I/URL-Shortener-Service/handlers/urlhandler"
   "github.com/MayhemI7I/URL-Shortener-Service/internal/storage"
   "github.com/MayhemI7I/URL-Shortener-Service/logger"
   "github.com/MayhemI7I/URL-Shortener-Service/utils"
   "go.uber.org/zap"
)

// initApp выполняет все необходимые иниты и возвращает готовые зависимости.
func initApp() (*config.Config, *urlhandler.URLHandler, *authhandler.AuthHandler, storage.Storage, error) {
   // Загружаем конфиг
   cfg := config.InitConfig()

   // Инициализируем логгер
   logger.InitLogger(cfg.LogLevel)

   // Инициализируем хранилище
   store, err := storage.NewStorage(*cfg)
   if err != nil {
   	return nil, nil, nil, nil, err
   }

   // Создаем генератор коротких URL
   genUrl := utils.NewGeneratorShortURL(cfg.URLLength)

   // Создаем обработчик URL
   urlHandler := urlhandler.NewURLHandler(store, genUrl)
   authHandler := authhandler.NewAuthHandler(store)

   return cfg, urlHandler, authHandler, store, nil
}

func main() {
   cfg, urlHandler, authHandler, store, err := initApp()
   if err != nil {
   	logger.Log.Fatalf("failed to initialize application: %v", err)
   }
   defer logger.CloseLogger()
   defer func() {
   	if err := store.Close(); err != nil {
   		logger.Log.Error("failed to close storage :", zap.Error(err))
   	}
   }()

   // Создаем HTTP multiplexer
   mux := http.NewServeMux()

   // Унифицированный конвейер для всех маршрутов
   mux.Handle("/", middleware.Conveyor(
   	http.HandlerFunc(urlHandler.HandURL),
   	middleware.WithLog,
   	zstd.Decompression,
   	zstd.Compression,
   ))

   mux.Handle("/api/login", middleware.Conveyor(
	http.HandlerFunc(authHandler.Login),
	middleware.WithLog,
	))

   mux.Handle("/api/shorten", middleware.Conveyor(
   	http.HandlerFunc(urlHandler.HandlePost),
   	middleware.WithLog,
   	zstd.Decompression,
   	zstd.Compression,
      middleware.Auth(store),
   ))

   mux.Handle("/api/user/urls", middleware.Conveyor(
   	http.HandlerFunc(urlHandler.HandleGetUserAllURLs),
   	middleware.WithLog,
      zstd.Decompression,
   	zstd.Compression,
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
   server := &http.Server{
   	Addr:         addr,
   	Handler:      mux,
   	ReadTimeout:  10 * time.Second,
   	WriteTimeout: 10 * time.Second,
   }
   logger.Log.Infof(time.Now().Format("2006-01-02 15:04:05")+"Server started on %s", addr)

   errChan := make(chan error, 1)
   go func() {
   	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
   		errChan <- err
   	}
   }()

   signChan := make(chan os.Signal, 1)
   signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

   select {
   case err := <-errChan:
   	return err
   case <-signChan:
   	logger.Log.Info("Shutting down server")
   	timeout, err := strconv.Atoi(cfg.ShutdownTimeout)
   	if err != nil {
   		logger.Log.Error("failed to parse shutdown timeout", zap.Error(err))
   		return err
   	}
   	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
   	defer cancel()
   	if err := server.Shutdown(ctx); err != nil {
   		logger.Log.Error("shutdown error", zap.Error(err))
   		return err
   	}
   	logger.Log.Info("Server stopped")
   	return nil
   }
}
