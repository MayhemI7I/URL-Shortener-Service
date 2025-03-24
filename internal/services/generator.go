package services

import (
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"
)

type GeneratorService struct {
	repo interfaces.URLRepository
	lenght uint16 // длина короткой ссылки
}

func NewGeneratorService(repo interfaces.URLRepository, lenght uint16) *GeneratorService {
	return &GeneratorService{repo: repo, lenght: lenght}
}

func (s *GeneratorService) Generate(ctx context.Context, longURL string) (string, error) {
	hash := sha256.Sum256([]byte(longURL))
		shortURL := base64.URLEncoding.EncodeToString(hash[:])
		if len(shortURL) < int(s.lenght) {
			return "", errors.New("generated short URL is too short")
		}
		logger.Log.Debug("Generated short URL: ", shortURL[:s.lenght])
		return shortURL[:s.lenght], nil
	 }
