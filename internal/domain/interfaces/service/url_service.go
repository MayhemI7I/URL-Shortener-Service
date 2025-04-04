package service

import (
	"context"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
)

// URLService определяет бизнес-логику для работы с URL
type URLService interface {
	// CreateShortURL создает новую короткую ссылку
	CreateShortURL(ctx context.Context, longURL, userID string) (*models.URLData, error)

	// GetOriginalURL получает оригинальный URL по короткой ссылке
	GetOriginalURL(ctx context.Context, shortURL string) (*models.URLInfo, error)

	// GetShortURL получает короткую ссылку по оригинальному URL
	GetShortURL(ctx context.Context, originalURL string) (*models.URLInfo, error)

	// GetUserURLs получает список URL пользователя
	GetUserURLs(ctx context.Context, userID string) ([]*models.URLInfo, error)

	// DeleteURL удаляет URL
	DeleteURL(ctx context.Context, shortURL string) error
}


