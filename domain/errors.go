package domain

import(
	"errors"

)

var (
    ErrInvalidURL      = errors.New("invalid URL format")
	ErrUserNotFound    = errors.New("user not found")
	ErrTokenNotFound   = errors.New("refresh token not found")
	ErrTokenExpired    = errors.New("refresh token expired")
	ErrURLNotFound = errors.New("URL not found")
	ErrURLExists = errors.New("URL already exists")
	
)