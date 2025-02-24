package zstd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZstdCompress(t *testing.T) {
	// Создаем тестовый обработчик
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	zstdCompress := Compression(nextHandler)

	// Таблица тестов
	tests := []struct {
		name           string
		acceptEncoding string
		expectedHeader string
		expectedBody   string
	}{
		{
			name:           "Accept-Encoding: zstd",
			acceptEncoding: "zstd",
			expectedHeader: "zstd",
			expectedBody:   "Hello, world!",
		},
		{
			name:           "Accept-Encoding: gzip",
			acceptEncoding: "gzip",
			expectedHeader: "",
			expectedBody:   "Hello, world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый запрос
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Accept-Encoding", tt.acceptEncoding)

			// Создаем тестовый ответ
			w := httptest.NewRecorder()

			// Вызываем ZstdCompress
			zstdCompress.ServeHTTP(w, req)

			// Проверяем, что ответ был сжат с использованием zstd
			assert.Equal(t, tt.expectedHeader, w.Header().Get("Content-Encoding"))
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestZstdDecompress(t *testing.T) {
	// Создаем тестовый обработчик
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	
	zstdDecompress := Decompression(nextHandler)

	// Таблица тестов
	tests := []struct {
		name           string
		contentEncoding string
		expectedHeader string
		expectedBody   string
	}{
		{
			name:           "Content-Encoding: zstd",
			contentEncoding: "zstd",
			expectedHeader: "",
			expectedBody:   "Hello, world!",
		},
		{
			name:           "Content-Encoding: gzip",
			contentEncoding: "gzip",
			expectedHeader: "",
			expectedBody:   "Hello, world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый запрос
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Content-Encoding", tt.contentEncoding)

			// Создаем тестовый ответ
			w := httptest.NewRecorder()

			// Вызываем ZstdDecompress
			zstdDecompress.ServeHTTP(w, req)

			// Проверяем, что ответ был распакован с использованием zstd
			assert.Equal(t, tt.expectedHeader, w.Header().Get("Content-Encoding"))
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
