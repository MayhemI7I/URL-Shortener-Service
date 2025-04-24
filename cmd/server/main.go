package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/app"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/middleware"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/zstd"
)

func main() {
	// Получаем экземпляр фасада приложения
	appFacade := app.GetInstance()

	// Инициализируем приложение
	if err := appFacade.Init(); err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}
	defer appFacade.GetDB().Close()

	// Получаем логгер
	logger := appFacade.GetLogger()

	// Создаем HTTP multiplexer
	mux := http.NewServeMux()

	// Получаем use cases из ядра
	core := appFacade.GetCore()
	urlService := core.URLService()
	authService := core.AuthService()

	// Унифицированный конвейер для всех маршрутов
	mux.Handle("/", middleware.Conveyor(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlService.HandURL(w, r)
		}),
		middleware.WithLog(logger),
		zstd.Decompression,
		zstd.Compression,
	))

	mux.Handle("/api/login", middleware.Conveyor(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authService.Login(w, r)
		}),
		middleware.WithLog(logger),
	))

	mux.Handle("/api/shorten", middleware.Conveyor(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlService.HandlePost(w, r)
		}),
		middleware.WithLog(logger),
		zstd.Decompression,
		zstd.Compression,
		middleware.Auth(core),
	))

	mux.Handle("/api/user/urls", middleware.Conveyor(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlService.HandleGetUserAllURLs(w, r)
		}),
		middleware.WithLog(logger),
		zstd.Decompression,
		zstd.Compression,
		middleware.Auth(core),
	))

	// Запускаем сервер
	cfg := appFacade.GetConfig()
	server := &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.HTTPConfig.GetReadTimeout()) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTPConfig.GetWriteTimeout()) * time.Second,
	}

	// Канал для сигналов завершения
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		logger.Info("Сервер запущен", "address", cfg.GetServerAddr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Ошибка запуска сервера", "error", err)
		}
	}()

	// Ожидание сигнала завершения
	<-done
	logger.Info("Получен сигнал завершения")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.HTTPConfig.GetShutdownTimeout())*time.Second)
	defer cancel()

	// Graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Ошибка при завершении работы сервера", "error", err)
	}

	logger.Info("Сервер остановлен")
}
