package file

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/repository"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"

	"go.uber.org/zap"
)

// FileStorage представляет файловое хранилище для URL
type FileStorage struct {
	urls     map[string]*models.URLData // Мапа коротких URL на данные URL
	mu       sync.Mutex                 // Мьютекс для потокобезопасного доступа
	file     *os.File                   // Файл для хранения данных
	filename string                     // Имя файла
}

// NewFileStorage создает и инициализирует новое файловое хранилище
func NewFileStorage(filename string) (repository.URLRepository, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Ошибка открытия файлового хранилища", zap.Error(err))
		return nil, err
	}

	storage := &FileStorage{
		urls:     make(map[string]*models.URLData),
		mu:       sync.Mutex{},
		file:     file,
		filename: filename,
	}

	if err := storage.load(); err != nil && err != io.EOF {
		logger.Log.Error("Ошибка загрузки данных из файла", zap.Error(err))
		return nil, err
	}

	return storage, nil
}

// load загружает данные из файла в память
func (s *FileStorage) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Сбрасываем указатель в начало файла
	if _, err := s.file.Seek(0, 0); err != nil {
		logger.Log.Error("Ошибка сброса указателя файла", zap.Error(err))
		return err
	}

	decoder := json.NewDecoder(s.file)
	if err := decoder.Decode(&s.urls); err != nil && err != io.EOF {
		logger.Log.Error("Ошибка декодирования файла", zap.Error(err))
		return err
	}

	return nil
}

// save сохраняет данные из памяти в файл
func (s *FileStorage) save() error {
	// Открываем файл заново для записи, перезаписывая содержимое
	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Ошибка открытия файла для записи", zap.Error(err))
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(s.urls); err != nil {
		logger.Log.Error("Ошибка кодирования данных в файл", zap.Error(err))
		return err
	}

	return nil
}

// AddUserURL сохраняет пару URL (короткий и оригинальный) для пользователя
func (s *FileStorage) AddUserURL(ctx context.Context, url *models.URLData) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, существует ли уже короткий URL
	if _, exists := s.urls[url.ShortURL]; exists {
		return models.ErrURLExists
	}

	// Сохраняем URL
	s.urls[url.ShortURL] = url

	// Записываем в файл
	return s.save()
}

// FindByOriginalURL ищет короткий URL по оригинальному URL
func (s *FileStorage) FindByOriginalURL(ctx context.Context, originalURL string) (*models.URLData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, data := range s.urls {
		if data.OrigURL == originalURL {
			return data, nil
		}
	}

	return nil, models.ErrURLNotFound
}

// FindByShortURL ищет оригинальный URL по короткой ссылке
func (s *FileStorage) FindByShortURL(ctx context.Context, shortURL string) (*models.URLData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok := s.urls[shortURL]
	if !ok {
		return nil, models.ErrURLNotFound
	}

	return data, nil
}

// ListByUser получает список всех URL пользователя
func (s *FileStorage) ListByUser(ctx context.Context, userID string) ([]*models.URLData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var urls []*models.URLData
	for _, data := range s.urls {
		if data.User.ID == userID {
			urls = append(urls, data)
		}
	}

	return urls, nil
}

// MarkAsDeleted помечает URL как удаленный
func (s *FileStorage) MarkAsDeleted(ctx context.Context, shortURL string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.urls[shortURL]
	if !ok {
		return models.ErrURLNotFound
	}

	// Удаляем URL из памяти
	delete(s.urls, shortURL)

	// Сохраняем изменения в файл
	return s.save()
}
