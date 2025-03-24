package usecases

import (
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/services"
)

// Core представляет основную бизнес-логику приложения,
// объединяя в себе все необходимые зависимости
type Core struct {
	// Репозитории
	urlRepo  interfaces.URLRepository
	userRepo interfaces.UserRepository

	// Сервисы
	urlService  interfaces.URLService
	authService interfaces.AuthService
	authConfig  interfaces.AuthConfig

	// Утилиты
	urlGenerator interfaces.URLGenerator
}

// NewCore создает новый экземпляр Core со всеми необходимыми зависимостями
func NewCore(
	urlRepo interfaces.URLRepository,
	userRepo interfaces.UserRepository,
	urlGenerator interfaces.URLGenerator,
	authConfig interfaces.AuthConfig,
) *Core {
	core := &Core{
		urlRepo:      urlRepo,
		userRepo:     userRepo,
		urlGenerator: urlGenerator,
		authConfig:   authConfig,
	}

	// Создаем сервисы через фабрики или DI-контейнер
	urlService := services.NewURLService(urlRepo, urlGenerator)
	authService := services.NewAuthService(userRepo, authConfig)

	// Передаем сервисы в use cases
	urlUseCase := NewURLUseCase(urlService)
	authUseCase := NewAuthUseCase(authService)

	core.urlService = urlUseCase
	core.authService = authUseCase

	return core
}

// URLService возвращает сервис для работы с URL
func (c *Core) URLService() interfaces.URLService {
	return c.urlService
}

// AuthService возвращает сервис для работы с аутентификацией
func (c *Core) AuthService() interfaces.AuthService {
	return c.authService
}
