package usecases

import (
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
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

	// Утилиты
	urlGenerator interfaces.URLGenerator
}

// NewCore создает новый экземпляр Core со всеми необходимыми зависимостями
func NewCore(
	urlRepo interfaces.URLRepository,
	userRepo interfaces.UserRepository,
	urlGenerator interfaces.URLGenerator,
) *Core {
	core := &Core{
		urlRepo:      urlRepo,
		userRepo:     userRepo,
		urlGenerator: urlGenerator,
	}

	// Инициализируем сервисы с репозиториями
	urlUseCase := NewURLUseCase(urlRepo)
	authUseCase := NewAuthUseCase(userRepo)

	core.urlService = urlUseCase
	core.authService = authUseCase

	return core
}

// NewURLService создает новый сервис для работы с URL
func NewURLService(repo interfaces.URLRepository, generator interfaces.URLGenerator) interfaces.URLService {
	return NewURLUseCase(repo)
}

// NewAuthService создает новый сервис для аутентификации
func NewAuthService(repo interfaces.UserRepository) interfaces.AuthService {
	return NewAuthUseCase(repo)
}

// URLService возвращает сервис для работы с URL
func (c *Core) URLService() interfaces.URLService {
	return c.urlService
}

// AuthService возвращает сервис для работы с аутентификацией
func (c *Core) AuthService() interfaces.AuthService {
	return c.authService
}

// NewAuthUseCase создает новый экземпляр сервиса аутентификации
func NewAuthUseCase(repo interfaces.UserRepository) interfaces.AuthService {
	return &AuthUseCase{repo: repo}
}

// AuthUseCase реализует интерфейс AuthService
type AuthUseCase struct {
	repo interfaces.UserRepository
}

// Методы AuthUseCase должны реализовать интерфейс AuthService
