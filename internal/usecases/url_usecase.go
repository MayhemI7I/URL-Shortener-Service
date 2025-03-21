// internal/usecase/url_usecase.go
package usecases

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
)

// URLUseCase содержит бизнес-логику для работы с URL
type URLUseCase struct {
	repo interfaces.URLRepository
}

// NewURLUseCase создаёт новый экземпляр URLUseCase
func NewURLUseCase(repo interfaces.URLRepository) *URLUseCase {
	return &URLUseCase{repo: repo}
}

// CreateShortURL создаёт короткую ссылку для оригинального URL
func (uc *URLUseCase) CreateShortURL(ctx context.Context, originalURL, userID string) (*models.URLData, error) {
	// Проверка наличия URL в базе
	existingURL, err := uc.repo.FindByOriginalURL(ctx, originalURL)
	if err == nil && existingURL != nil {
		return existingURL, nil
	}

	// Генерация короткой ссылки
	shortURL := uc.generateShortURL(originalURL)

	// Создание новой записи
	urlData := &models.URLData{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
	}

	// Сохранение в репозитории
	if err := uc.repo.AddUserURL(ctx, urlData); err != nil {
		return nil, err
	}

	return urlData, nil
}

// GetOriginalURL получает оригинальный URL по короткой ссылке
func (uc *URLUseCase) GetOriginalURL(ctx context.Context, shortURL string) (*models.URLData, error) {
	return uc.repo.FindByShortURL(ctx, shortURL)
}

// ListUserURLs получает все URL пользователя
func (uc *URLUseCase) ListUserURLs(ctx context.Context, userID string) ([]*models.URLData, error) {
	return uc.repo.FindByUserID(ctx, userID)
}

// DeleteURL удаляет URL по короткой ссылке
func (uc *URLUseCase) DeleteURL(ctx context.Context, shortURL, userID string) error {
	// Проверка владельца
	urlData, err := uc.repo.FindByShortURL(ctx, shortURL)
	if err != nil {
		return err
	}

	if urlData.UserID != userID {
		return errors.New("доступ запрещён")
	}

	return uc.repo.Delete(ctx, shortURL)
}

// Приватный метод для генерации короткого URL
func (uc *URLUseCase) generateShortURL(longURL string) string {
	hash := sha256.Sum256([]byte(longURL))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8]
}
