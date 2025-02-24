package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

)

var Log *zap.SugaredLogger 

func InitLogger(logLevel string) {
	// Устанавливаем уровень логирования
	var level zapcore.Level
	switch logLevel {
	case "1":
		level = zapcore.InfoLevel
	case "2":
		level = zapcore.DebugLevel
	default:
		level = zapcore.InfoLevel
	}

	// Форматирование логов
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"

	// Для консольного вывода
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	// Для файла
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Создаем Writer для вывода в консоль
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Создаем Writer для вывода в файл
	logFile, err := os.OpenFile("/short-url-db.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Не удалось создать файл для логов: " + err.Error())
	}
	fileWriter := zapcore.AddSync(logFile)

	// Создаем Core с комбинированными обработчиками
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, level),
		zapcore.NewCore(fileEncoder, fileWriter, level),
	)

	// Создаем логгер
	logger := zap.New(core)
	Log = logger.Sugar()
}

func CloseLogger() {
	if Log != nil {
		Log.Sync() // Закрытие логгера, запись всех оставшихся логов в файл
	}
}
