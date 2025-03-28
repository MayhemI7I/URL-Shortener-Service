package interfaces

// LoggerAdapter - интерфейс для адаптеров логгеров
type LoggerAdapter interface {
	// AsLogger преобразует конкретный логгер в интерфейс Logger
	AsLogger() Logger
}
