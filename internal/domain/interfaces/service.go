package interfaces

import (
	"context"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
)

// URLService определяет бизнес-логику для работы с URL
type URLService interface {
	// Create создает новую короткую ссылку
	Create(ctx context.Context, longURL string, userID string) (*models.URLPair, error)

	// Get получает оригинальный URL по короткой ссылке
	Get(ctx context.Context, shortURL string) (*models.URLPair, error)

	// List получает список URL пользователя
	List(ctx context.Context, userID string) ([]*models.URLPair, error)

	// Delete удаляет URL
	Delete(ctx context.Context, shortURL string) error
}

// AuthService определяет бизнес-логику для аутентификации
type AuthService interface {
	// Register регистрирует нового пользователя
	Register(ctx context.Context, email, password string) (*models.User, error)

	// Login аутентифицирует пользователя
	Login(ctx context.Context, email, password string) (*models.User, error)

	// RefreshTokens обновляет токены доступа
	RefreshTokens(ctx context.Context, refreshToken string) (*models.User, error)

	// ValidateToken проверяет токен
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}
