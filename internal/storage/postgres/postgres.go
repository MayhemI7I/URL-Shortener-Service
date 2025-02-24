package postgres

import (
	"database/sql"
	"local/logger"
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		logger.Log.Error(err)
		return nil, err
	}

	queryInitTable := `
	CREATE TABLE IF NOT EXISTS short_urls (
	id SERIAL PRIMARY KEY,
	short_url VARCHAR(255) UNIQUE NOT NULL,
	long_url VARCHAR(255) NOT NULL);
	`
	_ ,err = db.Exec(queryInitTable)
	if err != nil {
		logger.Log.Error("error creating table", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("DB and tables is ready")

	return &PostgresStorage{
		db: db,
	}, nil
}

func (pg *PostgresStorage) PingDataBase() error {
	return pg.db.Ping()
}

func (pg *PostgresStorage) Close() error {
	 if err := pg.db.Close(); err != nil {
		logger.Log.Error("error closing database connection:", zap.Error(err))
		return err
	 }
	 return nil
}

func(pg *PostgresStorage) Get(ctx context.Context, shortURL string) (string, error) {
	queryGet := `SELECT long_url FROM short_urls WHERE short_url = $1`
	row := pg.db.QueryRowContext(ctx, queryGet, shortURL)
	var longURL string
    err := row.Scan(longURL)
	if err != nil {
		logger.Log.Debug("error getting long url", zap.Error(err))
		return "", err
	}
	return longURL, nil

}

func (pg *PostgresStorage) Save(ctx context.Context, shortURL string, longURL string) error {
	querySave := `INSERT INTO short_urls (short_url, long_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING`
	_, err := pg.db.ExecContext(ctx, querySave, shortURL, longURL)
	if err != nil {
		logger.Log.Debug("error saving short url", zap.Error(err))
		return err
	}
	logger.Log.Debug("short url saved", zap.String("shortURL", shortURL))
	return nil
}



