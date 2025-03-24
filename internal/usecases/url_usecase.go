package usecases

import (
	"context"


	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/services"
)

// URLUseCase содержит бизнес-логику для работы с URL
type URLUseCase struct {
	service services.URLService
}

// NewURLUseCase создаёт новый экземпляр URLUseCase
func NewURLUseCase(service services.URLService) *URLUseCase {
	return &URLUseCase{service: service}
}

// CreateShortURL создаёт короткую ссылку для оригинального URL
func (uc *URLUseCase) CreateShortURL(ctx context.Context, originalURL, userID string) (*models.URLData, error) {

	return uc.service.CreateShortURL(ctx, originalURL, userID)

}

// GetOriginalURL получает оригинальный URL по короткой ссылке
func (uc *URLUseCase) GetOriginalURL(ctx context.Context, shortURL string) (*models.URLInfo, error) {
	return uc.service.GetOriginalURL(ctx, shortURL)
}

// List получает список URL пользователя (реализация интерфейса URLService)
func (uc *URLUseCase) List(ctx context.Context, userID string) ([]*models.URLInfo, error) {

	return uc.service.GetUserURLs(ctx, userID)
}

// DeleteURL удаляет URL (реализация интерфейса URLService)
func (uc *URLUseCase) DeleteURL(ctx context.Context, shortURL string) error {
	return uc.service.DeleteURL(ctx, shortURL)
}
