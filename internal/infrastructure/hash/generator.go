package hash

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
)

// HashGenerator реализует генератор URL на основе хеширования
type HashGenerator struct{}

// NewHashGenerator создает новый генератор URL
func NewHashGenerator() infrastructure.URLGenerator {
	return &HashGenerator{}
}

// Generate генерирует короткий URL на основе оригинального URL
func (g *HashGenerator) Generate(ctx context.Context, origURL string) (string, error) {
	// Создаем хеш от оригинального URL
	hash := sha256.Sum256([]byte(origURL))

	// Кодируем хеш в base64 и берем первые 8 символов
	shortURL := base64.URLEncoding.EncodeToString(hash[:])[:8]

	// Удаляем небезопасные символы
	shortURL = strings.ReplaceAll(shortURL, "+", "-")
	shortURL = strings.ReplaceAll(shortURL, "/", "_")

	return shortURL, nil
}
