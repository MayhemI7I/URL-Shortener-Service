package utils

import (
   "crypto/sha256"
   "encoding/base64"
   "errors"
   "local/logger"
)

type GeneratorShortURL struct {
   lenght uint16
}

func NewGeneratorShortURL(lenght uint16) *GeneratorShortURL {
   return &GeneratorShortURL{lenght: lenght}
}

func (gen *GeneratorShortURL) GenerateShortURL(longurl string) (string, error) {
   if longurl == "" || longurl == " " {
   	return "", errors.New("invalid URL for generate")
   }
   hash := sha256.Sum256([]byte(longurl))
   shortURL := base64.URLEncoding.EncodeToString(hash[:])
   if len(shortURL) < int(gen.lenght) {
   	return "", errors.New("generated short URL is too short")
   }
   logger.Log.Debug("Generated short URL: ", shortURL[:gen.lenght])
   return shortURL[:gen.lenght], nil
}
