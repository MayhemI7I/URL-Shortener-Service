package zstd

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/klauspost/compress/zstd"
)

// Compression добавляет сжатие ответов
func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, поддерживает ли клиент сжатие zstd
		if !supportsZstd(r) {
			next.ServeHTTP(w, r)
			return
		}

		// Создаем буфер для сжатого ответа
		buf := &bytes.Buffer{}
		writer, err := zstd.NewWriter(buf)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer writer.Close()

		// Создаем response writer для записи сжатого ответа
		rw := &responseWriter{
			ResponseWriter: w,
			writer:         writer,
			buf:            buf,
		}

		// Устанавливаем заголовки сжатия
		w.Header().Set("Content-Encoding", "zstd")
		w.Header().Add("Vary", "Accept-Encoding")

		// Вызываем следующий обработчик
		next.ServeHTTP(rw, r)

		// Завершаем сжатие и отправляем ответ
		writer.Close()
		w.Write(buf.Bytes())
	})
}

// Decompression добавляет распаковку запросов
func Decompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, сжат ли запрос
		if r.Header.Get("Content-Encoding") != "zstd" {
			next.ServeHTTP(w, r)
			return
		}

		// Создаем reader для распаковки
		reader, err := zstd.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer reader.Close()

		// Заменяем тело запроса на распакованное
		r.Body = io.NopCloser(reader)
		next.ServeHTTP(w, r)
	})
}

// responseWriter обертка для http.ResponseWriter для сжатия ответа
type responseWriter struct {
	http.ResponseWriter
	writer *zstd.Encoder
	buf    *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.writer.Write(b)
}

// supportsZstd проверяет, поддерживает ли клиент сжатие zstd
func supportsZstd(r *http.Request) bool {
	acceptEncoding := r.Header.Get("Accept-Encoding")
	return strings.Contains(acceptEncoding, "zstd")
}
