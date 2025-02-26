package file

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"local/internal/storage/postgres"
	"local/logger"
	"log"
	"os"
	"sync"
)

type Storage struct {
	urls     map[string]string
	longURLs map[string]string
	mu       sync.Mutex
	file     *os.File
}

func (us *Storage) Load() error {
	us.mu.Lock()
	defer us.mu.Unlock()

	decoder := json.NewDecoder(us.file)
	if err := decoder.Decode(&us.urls); err != nil && err != io.EOF {
		return err
	}

	// Восстанавливаем `longURLs` для быстрого поиска по длинному URL
	us.longURLs = make(map[string]string)
	for short, long := range us.urls {
		us.longURLs[long] = short
	}

	return nil
}

func NewFileStorage(filename string) (*Storage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	storage := &Storage{
		urls:     make(map[string]string),
		longURLs: make(map[string]string),
		mu:       sync.Mutex{},
		file:     file,
	}
	storage.Load()
	return storage, nil
}

func (us *Storage) Close() error {
	return us.file.Close()
}

func (us *Storage) Save(ctx context.Context, shortURL, longURL string) error {
	select {
	case <-ctx.Done():
		return ctx.Err() // Возвращаем ошибку, если контекст отменён
	default:
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	if shortURL == "" || longURL == "" {
		logger.Log.Errorf("Invalid argument: %s, %s", shortURL, longURL)
		return errors.New("invalid argument")
	}
	if _, exists := us.urls[shortURL]; exists {
		logger.Log.Infof("URL already exists: %s", shortURL)
		return errors.New("URL already exists")
	}
	us.urls[shortURL] = longURL

	// Вторичная проверка, чтобы не писать в файл, если контекст отменён
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	encoder := json.NewEncoder(us.file)
	if err := encoder.Encode(us.urls); err != nil {
		return err
	}

	logger.Log.Info("Saved: %s -> %s", shortURL, longURL)
	return nil
}

func (us *Storage) Get(ctx context.Context, shortUrl string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	if shortUrl == "" {
		log.Printf("Invalid argument: %s", shortUrl)
		return "", errors.New("invalid short URL argument")
	}
	value, ok := us.urls[shortUrl]
	if !ok {
		return "", errors.New("URL not found in storage")
	}

	logger.Log.Info("Retrieved: %s -> %s", shortUrl, value)
	return value, nil
}

func (us *Storage) FindByLongURL(ctx context.Context, longURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	us.mu.Lock()
	defer us.mu.Unlock()
	shortURL, ok := us.longURLs[longURL]
	if !ok {
		return "", postgres.ErrURLNotFound
	}
	return shortURL, nil

}
