package interfaces

import (
	"context"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
)

// URLRepository определяет методы для работы с URL в хранилище
type URLRepository interface {
	// Create создает новую запись URL
	AddUserURL(ctx context.Context, url *models.URLData) error

	// FindByOriginalURL получает сокращенный URL по оригинальному URL
	FindByOriginalURL(ctx context.Context, originalURL string) (*models.URLData, error)

	// FindByShortURL получает оригинальный URL по короткой ссылке
	FindByShortURL(ctx context.Context, shortURL string) (*models.URLData, error)

	// List получает список всех URL пользователя
	ListByUser(ctx context.Context, userID string) ([]*models.URLData, error)

	// Помечает URL как удаленный
	MarkAsDeleted(ctx context.Context, shortURL string) error

}

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	// GetUserById получает пользователя по ID
	GetUserById(ctx context.Context, userID string) (*models.User, error)


	// SaveRefreshToken сохраняет токен обновления для пользователя
	// Если пользователь не существует, он будет создан
	SaveRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error

	// GetUserByRefreshToken получает пользователя по токену обновления
	GetUserByRefreshToken(ctx context.Context, token string) (*models.User, error)
}
