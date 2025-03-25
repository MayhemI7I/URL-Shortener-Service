package memory

import (
	"context"
	"sync"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
)

// Storage представляет in-memory хранилище для URL
type Storage struct {
	urls map[string]*models.URLData // Мапа коротких URL на данные URL
	mu   sync.RWMutex               // Мьютекс для потокобезопасного доступа
}

// NewMemoryStorage создает и инициализирует новое in-memory хранилище
func NewMemoryStorage() (interfaces.URLRepository, error) {
	return &Storage{
		urls: make(map[string]*models.URLData),
	}, nil
}

// AddUserURL сохраняет пару URL (короткий и оригинальный) для пользователя
func (s *Storage) AddUserURL(ctx context.Context, url *models.URLData) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, существует ли уже короткий URL
	if _, exists := s.urls[url.ShortURL]; exists {
		return domain.ErrURLExists
	}

	// Сохраняем URL
	s.urls[url.ShortURL] = url
	return nil
}

// FindByOriginalURL ищет короткий URL по оригинальному URL
func (s *Storage) FindByOriginalURL(ctx context.Context, originalURL string) (*models.URLData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, data := range s.urls {
		if data.OrigURL == originalURL {
			return data, nil
		}
	}

	return nil, domain.ErrURLNotFound
}

// FindByShortURL ищет оригинальный URL по короткой ссылке
func (s *Storage) FindByShortURL(ctx context.Context, shortURL string) (*models.URLData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.urls[shortURL]
	if !ok {
		return nil, domain.ErrURLNotFound
	}

	return data, nil
}

// ListByUser получает список всех URL пользователя
func (s *Storage) ListByUser(ctx context.Context, userID string) ([]*models.URLData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var urls []*models.URLData
	for _, data := range s.urls {
		if data.User.ID == userID {
			urls = append(urls, data)
		}
	}

	return urls, nil
}

// MarkAsDeleted помечает URL как удаленный
func (s *Storage) MarkAsDeleted(ctx context.Context, shortURL string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.urls[shortURL]
	if !ok {
		return domain.ErrURLNotFound
	}

	// Помечаем URL как удаленный (в данной реализации просто удаляем)
	delete(s.urls, shortURL)
	return nil
}
