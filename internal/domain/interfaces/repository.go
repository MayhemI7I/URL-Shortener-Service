package interfaces

import (
	"context"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
)

// URLRepository определяет методы для работы с URL в хранилище
type URLRepository interface {
	// Create создает новую запись URL
	Create(ctx context.Context, url *models.URLData) error

	// Get получает URL по короткой ссылке
	Get(ctx context.Context, shortURL string) (*models.URLPair, error)

	// List получает список URL для пользователя
	List(ctx context.Context, userID string) ([]*models.URLPair, error)

	// Delete удаляет URL
	Delete(ctx context.Context, shortURL string) error
}

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	// Create создает нового пользователя
	Create(ctx context.Context, user *models.User) error

	// Get получает пользователя по ID
	Get(ctx context.Context, userID string) (*models.User, error)

	// SaveRefreshToken сохраняет токен обновления для пользователя
	SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error

	// GetUserByRefreshToken получает пользователя по токену обновления
	GetUserByRefreshToken(ctx context.Context, token string) (*models.User, error)
}
