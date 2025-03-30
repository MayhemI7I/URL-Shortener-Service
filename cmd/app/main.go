package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/auth"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/http"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository/postgres"
)

func main() {
	// Инициализация всех конфигураций компонентов
	httpConfig := http.NewHTTPConfig()
	loggerConfig := logger.NewLoggerConfig()
	dbConfig := postgres.NewDBConfig()
	jwtConfig := auth.NewJWTConfig()
	
	// Собираем их в общий контейнер конфигурации
	appConfig := config.NewAppConfig(
		httpConfig,
		loggerConfig, 
		dbConfig,
		jwtConfig,
	)
	
	// Добавляем флаги и парсим их
	appConfig.AddFlags()
	if err := appConfig.ParseFlags(); err != nil {
		log.Fatalf("Ошибка при парсинге флагов: %v", err)
	}
	
	// Загружаем значения из переменных окружения
	appConfig.LoadFromEnv()
	
	// Валидируем конфигурацию
	if err := appConfig.Validate(); err != nil {
		log.Fatalf("Ошибка валидации конфигурации: %v", err)
	}
	
	// Здесь можно инициализировать компоненты приложения
	log.Printf("Сервер запускается на %s", appConfig.GetServerAddr())
	
	// Создаем канал для сигналов завершения
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	
	<-done
	fmt.Println("Завершение работы сервера...")
} 