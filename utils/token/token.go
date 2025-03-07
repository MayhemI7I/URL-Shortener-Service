package token

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRefreshToken creates a new random refresh token
func GenerateRefreshToken(userID string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
