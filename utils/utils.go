package utils

import (
	"errors"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateShortURL(longurl string) (string,error) {
	if longurl == "" || longurl == " " {
		return"", errors.New("invalid URL for generate")
	}
	hash := sha256.Sum256([]byte(longurl))
	shortURL := base64.URLEncoding.EncodeToString(hash[:])
	return shortURL[:8],nil

}
