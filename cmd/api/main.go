package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/server"
)

func main() {
	// Инициализация конфигурации
	appConfig := config.NewAppConfig()
	appConfig.AddFlags()

	// Парсинг флагов
	if err := appConfig.ParseFlags(); err != nil {
		log.Fatalf("Ошибка парсинга флагов: %v", err)
	}

	// Загрузка конфигурации из переменных окружения
	appConfig.LoadFromEnv()

	// Валидация конфигурации
	if err := appConfig.Validate(); err != nil {
		log.Fatalf("Ошибка валидации конфигурации: %v", err)
	}

	// Инициализация логгера
	logger, err := logger.NewZapLogger(appConfig.Logger)
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer logger.Sync()

	// Создание контекста с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация сервера
	srv := server.NewServer(appConfig, logger)

	// Запуск сервера в горутине
	go func() {
		logger.Info("Запуск сервера", "addr", appConfig.GetServerAddr())
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Ошибка запуска сервера", "error", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Получен сигнал завершения работы")

	// Создание контекста с таймаутом для graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), appConfig.ShutdownTimeout)
	defer shutdownCancel()

	// Остановка сервера
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Ошибка при остановке сервера", "error", err)
	}

	logger.Info("Сервер успешно остановлен")
}
