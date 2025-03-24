package dto

import "time"

// URLResponse ответ с данными о сокращенном URL
type URLResponse struct {
    ShortURL     string    `json:"short_url"`
    OriginalURL  string    `json:"original_url"`
    CreatedAt    time.Time `json:"created_at"`
    VisitCount   int64     `json:"visit_count,omitempty"`
}

// URLListResponse список URL для пользователя
type URLListResponse struct {
    URLs []URLResponse `json:"urls"`//список URL
    Total int          `json:"total"`//количество URL
}

// UserResponse данные о пользователе
type UserResponse struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// TokenResponse ответ с токенами доступа
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"` 
    ExpiresIn    time.Duration    `json:"expires_in"` // Время жизни в секундах
}
