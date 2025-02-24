package file

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"local/logger"
	"log"
	"os"
	"sync"
)

type FileStorage struct {
	urls map[string]string
	mu   sync.Mutex
	file *os.File
}

func (us *FileStorage) Load() error {
	us.mu.Lock()
	defer us.mu.Unlock()
	decoder := json.NewDecoder(us.file)
	for {
		var entry map[string]string
		if err := decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for key, value := range entry {
			us.urls[key] = value
		}
	}
	return nil
}

func NewFileStorage(filename string) (*FileStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	storage := &FileStorage{
		urls: make(map[string]string),
		mu:   sync.Mutex{},
		file: file,
	}
	storage.Load()
	return storage, nil
}

func (us *FileStorage) Close() error {
	return us.file.Close()
}

func (us *FileStorage) Save(ctx context.Context, shortURL, longURL string) error {
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

func (us *FileStorage) Get(ctx context.Context, shortUrl string) (string, error) {
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
