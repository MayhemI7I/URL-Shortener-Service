package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"local/domain"
)

// Storage represents an in-memory storage for URL pairs and refresh tokens.
// It uses maps to store short URLs mapped to URL data and long URLs mapped to short URLs,
// with thread-safe access via a read-write mutex.
type Storage struct {
	urls     map[string]domain.URLData // Maps short URL to URLData
	longURLs map[string]string         // Maps long URL to short URL for reverse lookup
	mu       sync.RWMutex             // Mutex for thread-safe read/write operations
}

// NewMemoryStorage creates and initializes a new in-memory storage instance.
// It returns a pointer to Storage and an error if initialization fails.
func NewMemoryStorage() (*Storage, error) {
	return &Storage{
		urls:     make(map[string]domain.URLData),
		longURLs: make(map[string]string),
	}, nil
}

// newURLData creates a new URLData instance with the given short URL, original URL, and user ID.
// The CreatedAt field is set to the current time.
func newURLData(shortURL, origURL, userID string) *domain.URLData {
	return &domain.URLData{
		UserID: userID,
		URLPair: domain.URLPair{
			ShortURL:  shortURL,
			OrigURL:   origURL,
			CreatedAt: time.Now(),
		},
	}
}

// Save stores a new URL pair (short URL and original URL) for a given user in memory.
// It ensures the short URL is unique and the long URL doesn't already exist for the user.
// Returns an error if the context is canceled, the short URL exists, or the long URL is duplicated.
func (ms *Storage) Save(ctx context.Context, shortURL, origURL, userID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Check if the short URL already exists
	if _, exists := ms.urls[shortURL]; exists {
		return errors.New("short URL already exists")
	}

	// Check if the long URL already exists for this user
	for _, urlData := range ms.urls {
		if urlData.UserID == userID && urlData.URLPair.OrigURL == origURL {
			return errors.New("long URL already exists for user")
		}
	}

	// Create and store the new URL data
	urlData := newURLData(shortURL, origURL, userID)
	ms.urls[shortURL] = *urlData
	ms.longURLs[origURL] = shortURL

	return nil
}

// Get retrieves the original URL (long URL) for a given short URL from memory.
// Returns the original URL and an error if the short URL is not found or the context is canceled.
func (ms *Storage) Get(ctx context.Context, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	urlData, ok := ms.urls[shortURL]
	if !ok {
		return "", errors.New("short URL not found")
	}
	return urlData.URLPair.OrigURL, nil
}

// FindByLongURL retrieves the short URL for a given long URL from memory.
// Returns the short URL and an error if the long URL is not found or the context is canceled.
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

// Close performs any necessary cleanup for the in-memory storage.
// For this implementation, it does nothing and returns nil.
func (ms *Storage) Close() error {
	return nil
}