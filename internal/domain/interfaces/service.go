package interfaces

import (
	"context"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
)

// URLService определяет бизнес-логику для работы с URL
type URLService interface {
	// Create создает новую короткую ссылку
	Create(ctx context.Context, longURL string, userID string) (*models.URLData, error)

	// Get получает оригинальный URL по короткой ссылке
	Get(ctx context.Context, shortURL string) (*models.URLInfo, error)

	// List получает список URL пользователя
	List(ctx context.Context, userID string) ([]*models.URLInfo, error)

	// Delete удаляет URL
	Delete(ctx context.Context, shortURL string) error
}

// AuthService определяет бизнес-логику для аутентификации
type AuthService interface {
	// Register регистрирует нового пользователя
	Register(ctx context.Context, user *models.RefreshToken) (*models.User, error)

	// RefreshTokens обновляет токены пользователя
	RefreshTokens(ctx context.Context, user *models.User) (*models.User, error)

	// ValidateToken проверяет токен
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}


