package models

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// User represents a user in the system.
// It contains a unique identifier (UUID) for the user.
type User struct {
	ID string `json:"id" db:"user_id"` // UUID of the user
	RefreshToken
	AccessToken
}

// URLPair represents a pair of a short URL and its original (long) URL, including the creation timestamp.
type URLPair struct {
	ShortURL  string    `json:"short_url" db:"short_url"`    // The shortened URL
	OrigURL   string    `json:"original_url" db:"original_url"` // The original long URL
	DeletedFlag bool `json:"deleted_flag" db:"deleted_flag"` // Flag to check if the URL is deleted
	CreatedAt time.Time `json:"created_at" db:"created_at"`   // Timestamp when the URL pair was created

}

// URLData contains data for working with URLs and their associated user.
// It embeds URLPair to include short URL, original URL, and creation time, along with the user ID.
type URLData struct {
	User
	URLPair
}

// RefreshToken represents a refresh token and its metadata.
// It is used for refreshing access tokens and includes the token value, user ID, and expiration time.
type RefreshToken struct {
	Token     string    `json:"refresh_token" db:"refresh_token"`      // The refresh token value
	ExpiresAt time.Time `json:"refresh_expires_at" db:"expires_at"` // The expiration time of the refresh token
}

type AccessToken struct{
	Token string `json:"acces_token"` // The access token value
	ExpiresAt time.Time `json:"access_expires_at"` // The expiration time of the access token
}

type Claims struct {
	UserID string `json:"user_id"` // The ID (UUID) of the user associated with the token
	jwt.RegisteredClaims
}
