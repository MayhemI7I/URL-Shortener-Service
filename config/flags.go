package config

import (
   "local/logger"
   "os"

   "github.com/spf13/pflag"
   "go.uber.org/zap"
)

// Config represents the configuration for the application.
type Config struct {
   ServerAdress string
   ServerPort   string
   BaseURL      string
   LogLevel     string
   FileStorage  string
   DataBaseDSN string
}

// InitConfig initializes the configuration for the application.
func InitConfig() *Config {
   cfg := &Config{}

   // Define command-line flags
   pflag.StringVarP(&cfg.ServerAdress, "server-address", "s", "localhost", "Server address")
   pflag.StringVarP(&cfg.ServerPort, "server-port", "p", "8080", "Server port")
   pflag.StringVarP(&cfg.BaseURL, "base-url", "b", "http://localhost:8080", "Base URL for return server")
   pflag.StringVar(&cfg.LogLevel, "log-level", "1", "Log level")
   pflag.StringVarP(&cfg.FileStorage, "file-storage", "f", "/short-url-db.json", "Path to file storage")
   pflag.StringVarP(&cfg.DataBaseDSN, "database-dsn", "d", "host=localhost user=postgres password=1 dbname=usvideos sslmode=disable", "PostgreSQL DSN")

   // Override configuration with environment variables if they are set
   if envServerAdress := os.Getenv("SERVER_ADDRESS"); envServerAdress != "" {
   	cfg.ServerAdress = envServerAdress
   	logger.Log.Infof("Server address set to ", zap.String("address", envServerAdress))
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
   if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
   	cfg.DataBaseDSN = envDataBaseDSN
   	logger.Log.Infof("DATABASE_DSN set to ", zap.String("database", envDataBaseDSN))
   }

   // Parse command-line flags
   pflag.Parse()

   return cfg
}
