package models


// URLCreateRequest запрос на создание нового URL
type URLCreateRequest struct {
    OriginalURL []string `json:"original_url" validate:"required,url"`
   
}

// URLDeleteRequest запрос на удаление URL
type URLDeleteRequest struct {
    ShortURLs []string `json:"short_urls" validate:"required,min=1"`
}

