package config

import (
	"os"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"

	"time"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type JWT struct {
	Secret                 string
	AccessTokenExpiration        time.Duration
	RefreshTokenExpiration time.Duration
}

// Config represents the configuration for the application.
type Config struct {
	ServerAdress    string
	ServerPort      string
	BaseURL         string
	LogLevel        string
	FileStorage     string
	DataBaseDSN     string
	URLLength       uint16
	JWT             JWT
	ShutdownTimeout string
}

// InitConfig initializes the configuration for the application.
func InitConfig() *Config {
	cfg := &Config{
		JWT: JWT{
			Secret:                 "secret",
			AccessTokenExpiration:        24 * time.Hour,
			RefreshTokenExpiration: 720 * time.Hour, // 30 days
		},
	}

	// Define command-line flags
	pflag.StringVarP(&cfg.ServerAdress, "server-address", "s", "localhost", "Server address")
	pflag.StringVarP(&cfg.ServerPort, "server-port", "p", "8080", "Server port")
	pflag.StringVarP(&cfg.BaseURL, "base-url", "b", "http://localhost:8080", "Base URL for return server")
	pflag.StringVar(&cfg.LogLevel, "log-level", "debug", "Log level")
	pflag.StringVarP(&cfg.FileStorage, "file-storage", "f", "short-url-db.json", "Path to file storage")
	pflag.StringVarP(&cfg.DataBaseDSN, "database-dsn", "d", "postgres://postgres:1@localhost:5432/usvideos", "PostgreSQL DSN")
	pflag.Uint16VarP(&cfg.URLLength, "url-length", "l", 8, "URL length")
	pflag.StringVarP(&cfg.JWT.Secret, "jwt-secret", "j", "secret", "JWT secret")
	pflag.DurationVarP(&cfg.JWT.AccessTokenExpiration, "jwt-token-expiration", "e", 24*time.Hour, "JWT token expiration")
	pflag.DurationVarP(&cfg.JWT.RefreshTokenExpiration, "jwt-refresh-expiration", "r", 720*time.Hour, "JWT refresh token expiration")
	pflag.StringVarP(&cfg.ShutdownTimeout, "shutdown-timeout", "t", "10", "Shutdown timeout")
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
	if envJWTSecret := os.Getenv("JWT_SECRET"); envJWTSecret != "" {
		cfg.JWT.Secret = envJWTSecret
		logger.Log.Infof("JWT secret set from environment")
	}
	if envTokenExpiration := os.Getenv("JWT_TOKEN_EXPIRATION"); envTokenExpiration != "" {
		if duration, err := time.ParseDuration(envTokenExpiration); err == nil {
			cfg.JWT.AccessTokenExpiration = duration
			logger.Log.Infof("JWT token expiration set to ", zap.Duration("expiration", duration))
		}
	}
	if envRefreshExpiration := os.Getenv("JWT_REFRESH_EXPIRATION"); envRefreshExpiration != "" {
		if duration, err := time.ParseDuration(envRefreshExpiration); err == nil {
			cfg.JWT.RefreshTokenExpiration = duration
			logger.Log.Infof("JWT refresh token expiration set to ", zap.Duration("expiration", duration))
		}
	}
	if envShutdownTimeout := os.Getenv("SHUTDOWN_TIMEOUT"); envShutdownTimeout != "" {
		cfg.ShutdownTimeout = envShutdownTimeout
	}

	// Parse command-line flags
	pflag.Parse()

	return cfg
}
