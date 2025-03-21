package models

import (
	"time"
)

// User represents a user in the system.
// It contains a unique identifier (UUID) for the user.
type User struct {
	ID string `json:"id" db:"user_id"` // UUID of the user
	CreatedAt time.Time `json:"created_at" db:"created_at_user"` // Timestamp when the user was created
	LastLogin time.Time `json:"last_login" db:"last_login_user"` // Timestamp when the user was last logged in
}

// URLPair represents a pair of a short URL and its original (long) URL, including the creation timestamp.
type URLInfo struct {
	ShortURL  string    `json:"short_url" db:"short_url"`    // The shortened URL
	OrigURL   string    `json:"original_url" db:"original_url"` // The original long URL
	DeletedFlag bool `json:"deleted_flag" db:"deleted_flag"` // Flag to check if the URL is deleted
	CreatedAt time.Time `json:"created_at" db:"created_at_url_pair"`   // Timestamp when the URL pair was created

}

// URLData contains data for working with URLs and their associated user.
type URLData struct {
	User
	URLInfo
}


