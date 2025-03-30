package models

import "errors"

var (
	// ErrURLNotFound возвращается, когда URL не найден
	ErrURLNotFound = errors.New("URL not found")

	// ErrURLExists возвращается, когда URL уже существует
	ErrURLExists = errors.New("URL already exists")

	// ErrTokenNotFound возвращается, когда токен не найден
	ErrTokenNotFound = errors.New("token not found")

	// ErrInvalidToken возвращается, когда токен недействителен
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired возвращается, когда токен истек
	ErrTokenExpired = errors.New("token expired")

	// ErrUserNotFound возвращается, когда пользователь не найден
	ErrUserNotFound = errors.New("user not found")

	// ErrURLDeleted возвращается, когда URL удален
	ErrURLDeleted = errors.New("URL deleted")


	// ErrInvalidURLFormat возвращается, когда URL недействителен
	ErrInvalidURLFormat = errors.New("invalid URL format")
	
	
)
