package infrastructure

// Logger определяет интерфейс для логирования в приложении
type Logger interface {
	// Debug логирует сообщение на уровне отладки
	Debug(msg string, fields ...Field)

	// Info логирует информационное сообщение
	Info(msg string, fields ...Field)

	// Warn логирует предупреждающее сообщение
	Warn(msg string, fields ...Field)

	// Error логирует сообщение об ошибке
	Error(msg string, fields ...Field)

	// Fatal логирует критическое сообщение и завершает работу приложения
	Fatal(msg string, fields ...Field)
}

// Field представляет собой структуру для дополнительных полей в логе
type Field struct {
	Key   string
	Value interface{}
}

// NewField создает новое поле для лога
func NewField(key string, value interface{}) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

