package memory

import (
	"context"
	"errors"
	"sync"
)

type Storage struct {
	urls     map[string]string
	longURLs map[string]string
	mu       sync.RWMutex
}

func NewMemoryStorage() (*Storage, error) {
	return &Storage{
		urls:     make(map[string]string),
		longURLs: make(map[string]string),
	}, nil
}

func (ms *Storage) Save(ctx context.Context, shortURL, longURL string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.urls[shortURL] = longURL
	ms.longURLs[longURL] = shortURL
	return nil
}

func (ms *Storage) Get(ctx context.Context, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	longURL, ok := ms.urls[shortURL]
	if !ok {
		return "", errors.New("short URL not found")
	}
	return longURL, nil
}

func (ms *Storage) FindByLongURL(ctx context.Context, longURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	shortURL, ok := ms.longURLs[longURL]
	if !ok {
		return "", errors.New("long URL not found")
	}
	return shortURL, nil
}
func (ms *Storage) Close() error {
	return nil
}
