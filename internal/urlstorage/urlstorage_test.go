package urlstorage

import (
	"errors"
	"testing"
	"sync"

	"github.com/stretchr/testify/assert"
)

func TestSaveURL(t *testing.T) {
	storage := NewURLStorage()
	tests := []struct {
		name       string
		shortURL   string
		longURL    string
		wantExists bool
		wantURL    string
	}{
		{
			name:       "empty long URL",
			shortURL:   "abc123",
			longURL:    "",
			wantExists: false,
			wantURL:    "",
		},
		{
			name:       "successful save",
			shortURL:   "abc123",
			longURL:    "/long-url.com",
			wantExists: true,
			wantURL:    "/long-url.com",
		},
		{
			name:       "empty short URL",
			shortURL:   "",
			longURL:    "/long-url.com",
			wantExists: false,
			wantURL:    "",
		},
		{
			name:       "URL already exists",
			shortURL:   "abc123",
			longURL:    "/long-url.com",
			wantExists: true,
			wantURL:    "/long-url.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.SaveURL(tt.shortURL, tt.longURL)

			gotURL, exists := storage.urls[tt.shortURL]
			assert.Equal(t, tt.wantExists, exists, "exists mismatch")
			assert.Equal(t, tt.wantURL, gotURL, "URL mismatch in storage")
		})
	}
}

func TestGetURL(t *testing.T) {
	storage := NewURLStorage()
	storage.SaveURL("abc123", "/long-url.com")

	type want struct {
		longURL string
		err     error
	}
	tests := []struct {
		name     string
		shortURL string
		want     want
	}{
		{
			name:     "successful get",
			shortURL: "abc123",
			want: want{
				longURL: "/long-url.com",
				err:     nil,
			},
		},
		{
			name:     "empty short URL",
			shortURL: "",
			want: want{
				longURL: "",
				err:     errors.New("invalid short URL argument"),
			},
		},
		{
			name:     "URL not exists",
			shortURL: "not-exist",
			want: want{
				longURL: "",
				err:     errors.New("URL not found in storage"),
			},
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			longURL, err := storage.GetURL(tt.shortURL)
			assert.Equal(t, tt.want.longURL, longURL, "mismatched URL in storage")

			if tt.want.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConcurrency(t *testing.T) {
	storage := NewURLStorage()
	var wg sync.WaitGroup

	shortURL := "abc123"
	longURL := "/long-url.com"

	wg.Add(2)

	go func() {
		defer wg.Done()
		storage.SaveURL(shortURL, longURL)
	}()

	go func() {
		defer wg.Done()
		_, _ = storage.GetURL(shortURL)
	}()

	wg.Wait()

	// Проверяем, что URL сохранился
	gotURL, err := storage.GetURL(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, longURL, gotURL)
}
