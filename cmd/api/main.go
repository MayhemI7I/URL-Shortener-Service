package main

import (
	"log"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/auth"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/http"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository/file"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository/postgres"
)

func main() {
	// Создаем конфигурации компонентов
	httpConfig := http.NewHTTPConfig()
	loggerConfig := logger.NewLoggerConfig()
	dbConfig := postgres.NewDBConfig()
	fileStorageConfig := file.NewFileStorageConfig()
	jwtConfig := auth.NewJWTConfig()

	// Инициализируем конфигурацию приложения
	cfg, err := config.Init(
		httpConfig,
		loggerConfig,
		dbConfig,
		fileStorageConfig,
		jwtConfig,
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации конфигурации: %v", err)
	}

	// Создаем логгер
	loggerFactory := logger.NewLoggerFactory(cfg.LoggerConfig)
	appLogger, err := loggerFactory.CreateDefaultLogger()
	if err != nil {
		log.Fatalf("Ошибка создания логгера: %v", err)
	}

	// Создаем и запускаем приложение
	app := NewApp(cfg, appLogger)
	if err := app.Run(); err != nil {
		log.Fatalf("Ошибка запуска приложения: %v", err)
	}
}
