package usecases

import (
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/repository"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/services"
)

// Core представляет основную бизнес-логику приложения,
// объединяя в себе все необходимые зависимости
type Core struct {
	// Репозитории
	urlRepo  repository.URLRepository
	userRepo repository.UserRepository

	// Use Cases
	urlUseCase  *URLUseCase
	authUseCase *AuthUseCase
	authConfig  infrastructure.JWTConfigProvider

	// Утилиты
	urlGenerator infrastructure.URLGenerator
}

// NewCore создает новый экземпляр Core со всеми необходимыми зависимостями
func NewCore(
	urlRepo repository.URLRepository,
	userRepo repository.UserRepository,
	urlGenerator infrastructure.URLGenerator,
	authConfig infrastructure.JWTConfigProvider,
) *Core {
	core := &Core{
		urlRepo:      urlRepo,
		userRepo:     userRepo,
		urlGenerator: urlGenerator,
		authConfig:   authConfig,
	}

	// Создаем сервисы
	urlService := services.NewURLService(urlRepo, urlGenerator)
	authService := services.NewAuthService(userRepo, authConfig)

	// Создаем use cases
	core.urlUseCase = NewURLUseCase(urlService)
	core.authUseCase = NewAuthUseCase(authService)

	return core
}

// URLService возвращает use case для работы с URL
func (c *Core) URLService() *URLUseCase {
	return c.urlUseCase
}

// AuthService возвращает use case для работы с аутентификацией
func (c *Core) AuthService() *AuthUseCase {
	return c.authUseCase
}
