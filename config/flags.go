package config

import (
	"local/logger"
	"os"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Config struct {
	ServerAdress string
	ServerPort   string
	BaseURL      string
	LogLevel     string
	FileStorage  string

}

func InitConfig() *Config {
	cfg := &Config{}
	pflag.StringVarP(&cfg.ServerAdress, "SERVER_ADDRESS", "s", "localhost", "server address")
	pflag.StringVarP(&cfg.ServerPort, "SERVER_PORT", "p", "8080", "server port")
	pflag.StringVarP(&cfg.BaseURL, "BASE_URL", "b", "http://localhost:8080", "Base URL for return server")
	pflag.StringVar(&cfg.LogLevel, "LOG_LEVEL", "1", "log level")
	pflag.StringVar(&cfg.FileStorage, "FILE_STORAGE", "f", "/tmp/short-url-db.json")

	if envServerAdress := os.Getenv("SERVER_ADDRESS"); envServerAdress != "" {
		cfg.ServerAdress = envServerAdress
		logger.Log.Infof("Server address set to ", zap.String("asress", envServerAdress))
	}
	if envServerPort := os.Getenv("SERVER_PORT"); envServerPort != "" {
		cfg.ServerPort = envServerPort
		logger.Log.Infof("Server port set to ", zap.String("port", envServerPort))
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseURL = envBaseURL
		logger.Log.Infof("Base URL set to ", zap.String("url", envBaseURL))
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
		logger.Log.Infof("Log level set to ", zap.String("level", envLogLevel))
	}
	if envFileStorage := os.Getenv("FILE_STORAGE"); envFileStorage != "" {
		cfg.FileStorage = envFileStorage
		logger.Log.Infof("File storage set to ", zap.String("file", envFileStorage))
  
	}
	pflag.Parse()

	return cfg
}
