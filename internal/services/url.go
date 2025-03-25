package services

import (
	"context"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
)

// URLService реализует интерфейс interfaces.URLService
type URLService struct {
	repo      interfaces.URLRepository
	generator interfaces.URLGenerator
}

func NewURLService(repo interfaces.URLRepository, generator interfaces.URLGenerator) interfaces.URLService {
	return &URLService{repo: repo, generator: generator}
}

// CreateShortURL реализует метод создания короткой ссылки из interfaces.URLService
func (s *URLService) CreateShortURL(ctx context.Context, originalURL, userID string) (*models.URLData, error) {
	// Проверяем существование URL
	if existing, err := s.repo.FindByOriginalURL(ctx, originalURL); err == nil {
		return existing, nil
	}

	// Генерируем короткий URL
	shortURL, err := s.generator.Generate(ctx, originalURL)
	if err != nil {
		return nil, err
	}

	// Создаем новый URL
	url := &models.URLData{
		User: models.User{ID: userID},
		URLInfo: models.URLInfo{
			ShortURL: shortURL,
			OrigURL:  originalURL,
		},
	}

	// Сохраняем URL
	if err := s.repo.AddUserURL(ctx, url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLService) GetShortURL(ctx context.Context, originalURL string) (*models.URLInfo, error) {
	url, err := s.repo.FindByOriginalURL(ctx, originalURL)
	if err != nil {
		return nil, err
	}
	return &url.URLInfo, nil
}

// GetOriginalURL реализует метод получения оригинального URL по короткой ссылке из interfaces.URLService
func (s *URLService) GetOriginalURL(ctx context.Context, shortURL string) (*models.URLInfo, error) {
	url, err := s.repo.FindByShortURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}
	return &url.URLInfo, nil
}

// GetUserURLs реализует метод получения списка URL пользователя из interfaces.URLService
func (s *URLService) GetUserURLs(ctx context.Context, userID string) ([]*models.URLInfo, error) {
	urls, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*models.URLInfo, len(urls))
	for i, url := range urls {
		result[i] = &url.URLInfo
	}
	return result, nil
}

// DeleteURL реализует метод удаления URL из interfaces.URLService
func (s *URLService) DeleteURL(ctx context.Context, shortURL string) error {
	return s.repo.MarkAsDeleted(ctx, shortURL)
}
