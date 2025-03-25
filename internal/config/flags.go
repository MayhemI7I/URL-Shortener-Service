package config

import (
	"os"
	"strconv"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"

	"time"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type JWT struct {
	Secret                 string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

// GetJWTSecret возвращает секретный ключ JWT
func (j JWT) GetJWTSecret() string {
	return j.Secret
}

// GetAccessTokenExpiration возвращает время жизни access-токена
func (j JWT) GetAccessTokenExpiration() time.Duration {
	return j.AccessTokenExpiration
}

// GetRefreshTokenExpiration возвращает время жизни refresh-токена
func (j JWT) GetRefreshTokenExpiration() time.Duration {
	return j.RefreshTokenExpiration
}

type DB struct {
	DSN             string
	PoolSize        int
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Driver          string
}

type FileStorage struct {
	Path string
}

// Config represents the configuration for the application.
type Config struct {
	ServerAdress    string
	ServerPort      string
	BaseURL         string
	LogLevel        string
	FileStorage     FileStorage
	DB              DB
	URLLength       uint16
	JWT             JWT
	ShutdownTimeout string
}

// InitConfig initializes the configuration for the application.
func InitConfig() *Config {
	cfg := &Config{
		JWT: JWT{
			Secret:                 "secret",
			AccessTokenExpiration:  24 * time.Hour,
			RefreshTokenExpiration: 720 * time.Hour, // 30 days
		},
	}

	// Define command-line flags
	pflag.StringVarP(&cfg.ServerAdress, "server-address", "s", "localhost", "Server address")
	pflag.StringVarP(&cfg.ServerPort, "server-port", "p", "8080", "Server port")
	pflag.StringVarP(&cfg.BaseURL, "base-url", "b", "http://localhost:8080", "Base URL for return server")
	pflag.StringVar(&cfg.LogLevel, "log-level", "debug", "Log level")
	pflag.StringVarP(&cfg.FileStorage.Path, "file-storage", "f", "short-url-db.json", "Path to file storage")
	pflag.StringVarP(&cfg.DB.DSN, "database-dsn", "d", "postgres://postgres:1@localhost:5432/usvideos", "PostgreSQL DSN")
	pflag.IntVarP(&cfg.DB.PoolSize, "db-pool-size", "", 10, "Database connection pool size")
	pflag.IntVarP(&cfg.DB.MaxIdleConns, "db-max-idle-conns", "", 5, "Maximum number of idle database connections")
	pflag.IntVarP(&cfg.DB.MaxOpenConns, "db-max-open-conns", "", 25, "Maximum number of open database connections")
	pflag.DurationVarP(&cfg.DB.ConnMaxLifetime, "db-conn-max-lifetime", "", time.Hour, "Maximum lifetime of database connections")
	pflag.DurationVarP(&cfg.DB.ConnMaxIdleTime, "db-conn-max-idle-time", "", 30*time.Minute, "Maximum idle time of database connections")
	pflag.StringVarP(&cfg.DB.Driver, "db-driver", "", "postgres", "Database driver")
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
		cfg.FileStorage.Path = envFileStorage
		logger.Log.Infof("File storage set to ", zap.String("file", envFileStorage))
	}
	if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
		cfg.DB.DSN = envDataBaseDSN
		logger.Log.Infof("DATABASE_DSN set to ", zap.String("database", envDataBaseDSN))
	}
	if envDBPoolSize := os.Getenv("DB_POOL_SIZE"); envDBPoolSize != "" {
		if poolSize, err := strconv.Atoi(envDBPoolSize); err == nil {
			cfg.DB.PoolSize = poolSize
			logger.Log.Infof("Database pool size set to ", zap.Int("pool_size", poolSize))
		}
	}
	if envDBMaxIdleConns := os.Getenv("DB_MAX_IDLE_CONNS"); envDBMaxIdleConns != "" {
		if maxIdleConns, err := strconv.Atoi(envDBMaxIdleConns); err == nil {
			cfg.DB.MaxIdleConns = maxIdleConns
			logger.Log.Infof("Database max idle connections set to ", zap.Int("max_idle_conns", maxIdleConns))
		}
	}
	if envDBMaxOpenConns := os.Getenv("DB_MAX_OPEN_CONNS"); envDBMaxOpenConns != "" {
		if maxOpenConns, err := strconv.Atoi(envDBMaxOpenConns); err == nil {
			cfg.DB.MaxOpenConns = maxOpenConns
			logger.Log.Infof("Database max open connections set to ", zap.Int("max_open_conns", maxOpenConns))
		}
	}
	if envDBConnMaxLifetime := os.Getenv("DB_CONN_MAX_LIFETIME"); envDBConnMaxLifetime != "" {
		if duration, err := time.ParseDuration(envDBConnMaxLifetime); err == nil {
			cfg.DB.ConnMaxLifetime = duration
			logger.Log.Infof("Database connection max lifetime set to ", zap.Duration("conn_max_lifetime", duration))
		}
	}
	if envDBConnMaxIdleTime := os.Getenv("DB_CONN_MAX_IDLE_TIME"); envDBConnMaxIdleTime != "" {
		if duration, err := time.ParseDuration(envDBConnMaxIdleTime); err == nil {
			cfg.DB.ConnMaxIdleTime = duration
			logger.Log.Infof("Database connection max idle time set to ", zap.Duration("conn_max_idle_time", duration))
		}
	}
	if envDBDriver := os.Getenv("DB_DRIVER"); envDBDriver != "" {
		cfg.DB.Driver = envDBDriver
		logger.Log.Infof("Database driver set to ", zap.String("driver", envDBDriver))
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
