package jwtutil

import (
	"os"
	"time"
	"fmt"
	"crypto/rand"
	"encoding/hex"
	"go.uber.org/zap"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"

	"github.com/golang-jwt/jwt/v4"

)
var (
    AccessTokenExpiration  = 15 * time.Minute // Время жизни access-токена: 15 минут
    RefreshTokenExpiration = 7 * 24 * time.Hour // Время жизни refresh-токена: 7 дней
)

// GenerateAccessToken creates a new JWT access token
func GenerateAccessToken(userID, secretKey string) (string, error) {
	claims := &models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		logger.Log.Error("Error generating refresh token", zap.Error(err))
		return "", err
	}
	return hex.EncodeToString(b), nil
}


// ParseJWT parses and validates a JWT access token, returning the claims
func ParseJWT(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// IsTokenExpired checks if the token is expired
func IsTokenExpired(err error) bool {
	if ve, ok := err.(*jwt.ValidationError); ok {
		return ve.Errors&jwt.ValidationErrorExpired != 0
	}
	return false
}

