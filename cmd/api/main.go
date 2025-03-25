package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository/postgres"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/usecases"
	"github.com/MayhemI7I/URL-Shortener-Service/pkg/utils/urlutil"
)

func main() {
	// Инициализация конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация репозиториев
	urlRepo, err := postgres.NewURLRepository(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to create URL repository: %v", err)
	}

	userRepo, err := postgres.NewUserRepository(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	// Инициализация утилит
	urlGenerator := urlutil.NewURLGenerator()

	// Создание ядра приложения
	core := usecases.NewCore(
		urlRepo,
		userRepo,
		urlGenerator,
		cfg.Auth,
	)

	// Пример использования ядра
	ctx := context.Background()

	// Использование URL use case
	urlService := core.URLService()
	shortURL, err := urlService.CreateShortURL(ctx, "https://example.com")
	if err != nil {
		log.Printf("Failed to create short URL: %v", err)
	}
	log.Printf("Created short URL: %s", shortURL)

	// Использование Auth use case
	authService := core.AuthService()
	user, err := authService.Register(ctx)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
	}
	log.Printf("Registered user: %s", user.ID)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
} 