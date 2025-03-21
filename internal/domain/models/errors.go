package models

import(
	"errors"

)

var (
    ErrInvalidURL      = errors.New("invalid URL format")//ошибка неверного формата URL
	ErrUserNotFound    = errors.New("user not found")//ошибка неверного пользователя
	ErrTokenNotFound   = errors.New("refresh token not found")//ошибка неверного токена
	ErrTokenExpired    = errors.New("refresh token expired")//ошибка просроченного токена
	ErrURLNotFound = errors.New("URL not found")//ошибка неверного URL
	ErrURLExists = errors.New("URL already exists")//ошибка существующего URL
	ErrURLDeleted = errors.New("URL deleted")//ошибка URL удлен
	
)
