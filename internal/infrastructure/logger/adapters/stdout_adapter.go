package adapters

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
)

// StdoutLogger - реализация логгера с выводом в stdout
type StdoutLogger struct {
	logger *log.Logger
	level  string
}

// NewStdoutLogger - создает новый логгер для вывода в stdout
func NewStdoutLogger(level string) interfaces.Logger {
	return &StdoutLogger{
		logger: log.New(os.Stdout, "", 0),
		level:  strings.ToUpper(level),
	}
}

// Debug - реализует метод интерфейса Logger
func (s *StdoutLogger) Debug(msg string, fields ...interfaces.Field) {
	if s.level != "DEBUG" && s.level != "TRACE" {
		return
	}
	s.log("DEBUG", msg, fields...)
}

// Info - реализует метод интерфейса Logger
func (s *StdoutLogger) Info(msg string, fields ...interfaces.Field) {
	if s.level == "WARN" || s.level == "ERROR" || s.level == "FATAL" {
		return
	}
	s.log("INFO", msg, fields...)
}

// Warn - реализует метод интерфейса Logger
func (s *StdoutLogger) Warn(msg string, fields ...interfaces.Field) {
	if s.level == "ERROR" || s.level == "FATAL" {
		return
	}
	s.log("WARN", msg, fields...)
}

// Error - реализует метод интерфейса Logger
func (s *StdoutLogger) Error(msg string, fields ...interfaces.Field) {
	if s.level == "FATAL" {
		return
	}
	s.log("ERROR", msg, fields...)
}

// Fatal - реализует метод интерфейса Logger
func (s *StdoutLogger) Fatal(msg string, fields ...interfaces.Field) {
	s.log("FATAL", msg, fields...)
	os.Exit(1)
}

// log - форматирует и выводит лог
func (s *StdoutLogger) log(level string, msg string, fields ...interfaces.Field) {
	// Форматируем текущее время
	now := time.Now().Format("2006-01-02 15:04:05.000")

	// Создаем поля для вывода
	var sb strings.Builder

	for _, field := range fields {
		sb.WriteString(fmt.Sprintf(" %s=%v", field.Key, field.Value))
	}

	// Выводим сообщение с цветами
	if level == "ERROR" || level == "FATAL" {
		s.logger.Printf("\033[1;31m%s [%s] %s%s\033[0m", now, level, msg, sb.String())
	} else if level == "WARN" {
		s.logger.Printf("\033[1;33m%s [%s] %s%s\033[0m", now, level, msg, sb.String())
	} else if level == "DEBUG" {
		s.logger.Printf("\033[1;36m%s [%s] %s%s\033[0m", now, level, msg, sb.String())
	} else {
		s.logger.Printf("\033[1;32m%s [%s] %s%s\033[0m", now, level, msg, sb.String())
	}
}
