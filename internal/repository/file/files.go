package file

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/domain"
	"github.com/MayhemI7I/URL-Shortener-Service/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/utils/jwtutil"
	"go.uber.org/zap"
)

type Storage struct {
	urls     map[string]domain.URLData // Хранит данные URL с учетом userID
	mu       sync.Mutex
	file     *os.File
	filename string // Для переоткрытия файла при записи
}

func NewFileStorage(filename string) (*Storage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("failed to open file storage", zap.Error(err))
		return nil, err
	}
	storage := &Storage{
		urls:     make(map[string]domain.URLData),
		mu:       sync.Mutex{},
		file:     file,
		filename: filename,
	}
	if err := storage.Load(); err != nil {
		logger.Log.Error("failed to load initial data from file", zap.Error(err))
		return nil, err
	}
	return storage, nil
}

func (s *Storage) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	decoder := json.NewDecoder(s.file)
	if err := decoder.Decode(&s.urls); err != nil && err != io.EOF {
		logger.Log.Error("failed to decode file data", zap.Error(err))
		return err
	}

	// Сбрасываем указатель файла в начало после чтения
	if _, err := s.file.Seek(0, 0); err != nil {
		logger.Log.Error("failed to reset file pointer", zap.Error(err))
		return err
	}

	return nil
}

func (s *Storage) saveToFile() error {
	// Переоткрываем файл для записи, чтобы переписать его полностью
	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("failed to open file for writing", zap.Error(err))
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(s.urls); err != nil {
		logger.Log.Error("failed to encode data to file", zap.Error(err))
		return err
	}
	return nil
}

func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.file.Close()
}

// Save сохраняет короткий URL для пользователя
func (s *Storage) Save(ctx context.Context, shortURL, origURL, userID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if shortURL == "" || origURL == "" || userID == "" {
		logger.Log.Error("invalid arguments",
			zap.String("short_url", shortURL),
			zap.String("orig_url", origURL),
			zap.String("user_id", userID))
		return errors.New("invalid argument")
	}

	// Проверяем, существует ли shortURL для этого пользователя
	for _, data := range s.urls {
		if data.ShortURL == shortURL && data.User.ID == userID {
			logger.Log.Info("URL already exists for user",
				zap.String("short_url", shortURL),
				zap.String("user_id", userID))
			return errors.New("URL already exists")
		}
	}

	urlData := domain.URLData{
		User: domain.User{ID: userID},
		URLPair: domain.URLPair{
			ShortURL:  shortURL,
			OrigURL:   origURL,
			CreatedAt: time.Now(),
		},
	}
	s.urls[shortURL] = urlData

	if err := s.saveToFile(); err != nil {
		return err
	}

	logger.Log.Info("short URL saved",
		zap.String("short_url", shortURL),
		zap.String("orig_url", origURL),
		zap.String("user_id", userID))
	return nil
}

// Get возвращает длинный URL по короткому для пользователя
func (s *Storage) Get(ctx context.Context, shortURL, userID string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if shortURL == "" || userID == "" {
		logger.Log.Error("invalid arguments",
			zap.String("short_url", shortURL),
			zap.String("user_id", userID))
		return "", errors.New("invalid argument")
	}

	data, ok := s.urls[shortURL]
	if !ok || data.User.ID != userID {
		logger.Log.Debug("short URL not found",
			zap.String("short_url", shortURL),
			zap.String("user_id", userID))
		return "", domain.ErrURLNotFound
	}

	logger.Log.Info("retrieved URL",
		zap.String("short_url", shortURL),
		zap.String("long_url", data.OrigURL),
		zap.String("user_id", userID))
	return data.OrigURL, nil
}

// FindByLongURL ищет короткий URL по длинному для пользователя
func (s *Storage) FindByLongURL(ctx context.Context, longURL, userID string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for shortURL, data := range s.urls {
		if data.OrigURL == longURL && data.User.ID == userID {
			return shortURL, nil
		}
	}

	logger.Log.Debug("short URL not found for long URL",
		zap.String("long_url", longURL),
		zap.String("user_id", userID))
	return "", domain.ErrURLNotFound
}

// GetUserAllURLs возвращает все URL пользователя
func (s *Storage) GetUserAllURLs(ctx context.Context, userID string) ([]domain.URLData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var userURLs []domain.URLData
	for _, data := range s.urls {
		if data.User.ID == userID {
			userURLs = append(userURLs, data)
		}
	}

	if len(userURLs) == 0 {
		logger.Log.Debug("no URLs found for user", zap.String("user_id", userID))
		return nil, nil
	}

	logger.Log.Debug("retrieved URLs for user",
		zap.String("user_id", userID),
		zap.Int("count", len(userURLs)))
	return userURLs, nil
}

// SaveRefreshToken сохраняет refresh-токен
func (s *Storage) SaveRefreshToken(ctx context.Context, refreshToken *domain.RefreshToken) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Обновляем данные пользователя или создаем новые
	data, exists := s.urls[refreshToken.User.ID] // Используем userID как ключ для простоты
	if !exists {
		data = domain.URLData{
			User: domain.User{ID: refreshToken.User.ID},
		}
	}
	data.User.RefreshToken = refreshToken.Token
	data.User.RefreshExpiresAt = refreshToken.ExpiresAt
	s.urls[refreshToken.User.ID] = data

	if err := s.saveToFile(); err != nil {
		return err
	}

	logger.Log.Debug("refresh token saved",
		zap.String("refresh_token", refreshToken.Token),
		zap.String("user_id", refreshToken.User.ID))
	return nil
}

// GetUserIDByRefreshToken возвращает user_id по refresh-токену
func (s *Storage) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, data := range s.urls {
		if data.User.RefreshToken == refreshToken {
			return data.User.ID, nil
		}
	}

	logger.Log.Debug("refresh token not found",
		zap.String("refresh_token", refreshToken))
	return "", domain.ErrTokenNotFound
}

// GetNewAccessToken генерирует новый access-токен
func (s *Storage) GetNewAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
	select {
	case <-ctx.Done():
		return "", "", ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var userData domain.URLData
	var found bool
	for _, data := range s.urls {
		if data.User.RefreshToken == refreshToken {
			userData = data
			found = true
			break
		}
	}

	if !found {
		logger.Log.Debug("refresh token not found",
			zap.String("refresh_token", refreshToken))
		return "", "", domain.ErrTokenNotFound
	}

	newRefreshToken := refreshToken
	if time.Now().After(userData.User.RefreshExpiresAt) {
		logger.Log.Debug("refresh token expired",
			zap.String("refresh_token", refreshToken))
		var err error
		newRefreshToken, err = jwtutil.GenerateRefreshToken()
		if err != nil {
			logger.Log.Error("failed to generate new refresh token", zap.Error(err))
			return "", "", err
		}
		newExpiresAt := time.Now().Add(jwtutil.RefreshTokenExpiration)
		
		// Создаем объект RefreshToken для сохранения
		tokenObj := &domain.RefreshToken{
			User: domain.User{ID: userData.User.ID},
			Token: newRefreshToken,
			ExpiresAt: newExpiresAt,
		}
		
		// Сохраняем через обновленный метод
		if err := s.SaveRefreshToken(ctx, tokenObj); err != nil {
			return "", "", err
		}
	}

	secretKey := os.Getenv("JWT_SECRET")
	accessToken, err := jwtutil.GenerateAccessToken(userData.User.ID, secretKey)
	if err != nil {
		logger.Log.Error("failed to generate access token", zap.Error(err))
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// DeleteRefreshToken удаляет refresh-токен
func (s *Storage) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for shortURL, data := range s.urls {
		if data.User.RefreshToken == refreshToken {
			data.User.RefreshToken = ""
			data.User.RefreshExpiresAt = time.Time{}
			s.urls[shortURL] = data
			if err := s.saveToFile(); err != nil {
				return err
			}
			logger.Log.Debug("refresh token deleted",
				zap.String("refresh_token", refreshToken))
			return nil
		}
	}

	logger.Log.Debug("refresh token not found",
		zap.String("refresh_token", refreshToken))
	return domain.ErrTokenNotFound
}