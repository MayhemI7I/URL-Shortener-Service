package models

import (
	"fmt"
	"time"
)

func (u *URLInfo) Validate() error {
	if u.OrigURL == "" {
		return fmt.Errorf("original URL is required")
	}
	if u.ShortURL == "" {
		return fmt.Errorf("short URL is required")
	}
	return nil
}

func (ud *URLData) Validate() error {
	if ud.ID == "" {
		return fmt.Errorf("user ID is required")
	}
	if ud.ShortURL == "" {
		return fmt.Errorf("short URL is required")
	}
	if ud.OrigURL == "" {
		return fmt.Errorf("original URL is required")
	}
	return nil
}

func (rt *RefreshToken) Validate() error {
	if rt.Token == "" {
		return fmt.Errorf("refresh token is required")
	}
	if rt.ExpiresAt.IsZero() || rt.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("invalid expiration for refresh token")
	}
	return nil
}

func (u *User) Validate() error {
	if u.ID == "" {
		return fmt.Errorf("user ID is required")
	}
	return nil
}
