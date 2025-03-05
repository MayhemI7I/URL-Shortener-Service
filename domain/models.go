package domain

import "time"

// User represents a user in the system.
// It contains a unique identifier (UUID) for the user.
type User struct {
	ID string `json:"id"` // UUID of the user
}

// URLPair represents a pair of a short URL and its original (long) URL, including the creation timestamp.
type URLPair struct {
	ShortURL  string    `json:"short_url"`  // The shortened URL
	OrigURL   string    `json:"original_url"` // The original long URL
	CreatedAt time.Time `json:"created_at"`  // Timestamp when the URL pair was created
}

// URLData contains data for working with URLs and their associated user.
// It embeds URLPair to include short URL, original URL, and creation time, along with the user ID.
type URLData struct {
	UserID string `json:"user_id"` // The ID (UUID) of the user who owns the URL
	URLPair
}


// RefreshToken represents a refresh token and its metadata.
// It is used for refreshing access tokens and includes the token value, user ID, and expiration time.
type RefreshToken struct {
	Token     string    `json:"token"`     // The refresh token value
	UserID    string    `json:"user_id"`   // The ID (UUID) of the user associated with the token
	ExpiresAt time.Time `json:"expires_at"` // The expiration time of the refresh token
}

