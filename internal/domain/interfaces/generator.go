package interfaces

import "context"

// URLGenerator определяет интерфейс для генерации коротких URL
type URLGenerator interface {
    // Generate генерирует короткую ссылку
    Generate(ctx context.Context, origURL string) (string, error)
}
