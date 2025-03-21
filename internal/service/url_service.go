package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
)

// urlService реализует интерфейс interfaces.URLService.
type urlService struct {
	repo interfaces.URLRepository // Используем интерфейс репозитория
}

// NewURLService создает новый экземпляр URLService.
func NewURLService(repo interfaces.URLRepository) interfaces.URLService {
	return &urlService{repo: repo}
}

// Create создает новую короткую ссылку.
func (s *urlService) AddUserURL(ctx context.Context, shortURL, origURL, userID string) error {
	// Проверяем нет ли уже в базе
	existingURL, err := s.repo.GetByOriginalURL(ctx, origURL)
	if err == nil && existingURL != nil {
		return existingURL, nil // Возвращаем уже созданный
	}

	// Генерация короткой ссылки (простой пример)
	shortURL := generateShortURL(origURL)

	urlData := &models.URLData{
		ShortURL:  shortURL,
		OrigURL:   origURL,
		UserID:    userID,
	}

	err = s.repo.AddUserURL(ctx, shortURL, origURL, userID) // Сохраняем через репозиторий
	if err != nil {
		return nil, err
	}

	return urlPair, nil
}

// Get получает оригинальный URL по короткой ссылке.
func (s *urlService) Get(ctx context.Context, shortURL string) (*models.URLPair, error) {
	return s.repo.Get(ctx, shortURL)
}

// List получает список URL пользователя.
func (s *urlService) ListUserURLs(ctx context.Context, userID string) ([]*models.URLData, error) {
	return s.repo.ListByUser(ctx, userID) // Получаем через репозиторий
}

// Delete удаляет URL.
func (s *urlService) Delete(ctx context.Context, shortURL string) error {
	return s.repo.Delete(ctx, shortURL) // Удаляем через репозиторий
}

func generateShortURL(longURL string) string {
	hash := sha256.Sum256([]byte(longURL))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8] // Берем первые 8 символов
} 