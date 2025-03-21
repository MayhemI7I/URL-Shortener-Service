package models

import (
    "github.com/golang-jwt/jwt/v4"
	"time"
)

// Claims представляет собой данные в JWT-токене
type Claims struct {
    UserID string `json:"user_id"`//UUID of the user
    jwt.RegisteredClaims
}

type RefreshToken struct {
	User
	Token     string    `json:"token" db:"refresh_token"`      // The refresh token value
	ExpiresAt time.Time `json:"expires_at" db:"refresh_token_expires_at"` // The expiration time of the refresh token
}
