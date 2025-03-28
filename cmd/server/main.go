package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/http/middleware"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/usecases"
	"github.com/MayhemI7I/URL-Shortener-Service/pkg/utils/urlutil"
	"go.uber.org/zap"
)

// initApp выполняет все необходимые иниты и возвращает готовые зависимости.
func initApp() (*config.Config, *usecases.Core, error) {
	// Загружаем конфиг
	cfg := config.InitConfig()

	// Инициализируем логгер
	logger := logger.NewLogger(cfg.LogLevel)

	// Инициализация репозиториев через селектор
	storageSelector := repository.NewStorageSelector(cfg)
	urlRepo, err := storageSelector.SelectURLStorage()
	if err != nil {
		return nil, nil, err
	}

	// Инициализация утилит
	urlGenerator := urlutil.NewURLGenerator()

	// Создание ядра приложения
	core := usecases.NewCore(
		urlRepo,
		urlGenerator,
		cfg.JWT,
	)

	return cfg, core, nil
}

func main() {
	cfg, core, err := initApp()
	if err != nil {
		logger.Log.Fatalf("failed to initialize application: %v", err)
	}
	defer logger.CloseLogger()

	// Создаем HTTP multiplexer
	mux := http.NewServeMux()

	// Получаем use cases из ядра
	urlService := core.URLService()
	authService := core.AuthService()

	// Унифицированный конвейер для всех маршрутов
	mux.Handle("/", middleware.Conveyor(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlService.HandURL(w, r)
		}),
		middleware.WithLog,
		zstd.Decompression,
		zstd.Compression,
	))

	mux.Handle("/api/login", middleware.Conveyor(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authService.Login(w, r)
		}),
		middleware.WithLog,
	))

	mux.Handle("/api/shorten", middleware.Conveyor(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlService.HandlePost(w, r)
		}),
		middleware.WithLog,
		zstd.Decompression,
		zstd.Compression,
		middleware.Auth(core),
	))

	mux.Handle("/api/user/urls", middleware.Conveyor(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlService.HandleGetUserAllURLs(w, r)
		}),
		middleware.WithLog,
		zstd.Decompression,
		zstd.Compression,
		middleware.Auth(core),
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
