package postgres

import (
   "context"
   "database/sql"
   "fmt"
   "local/logger"
   "local/domain"
   "time"

   "errors"
   "github.com/jmoiron/sqlx"
   //"github.com/jackc/pgerrcode"
   "go.uber.org/zap"

   _ "github.com/jackc/pgx/v5/stdlib"
)


type PostgresStorage struct {
   db *sqlx.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
   db, err := sqlx.Open("pgx", dsn) // Используем sqlx.Open вместо sql.Open
   if err != nil {
   	logger.Log.Error(err)
   	return nil, err
   }
   if err := db.Ping(); err != nil {
   	logger.Log.Error(err)
   	return nil, err
   }

   queryInitShortURLs := `
		CREATE TABLE IF NOT EXISTS short_urls (
			id SERIAL PRIMARY KEY,
			user_id UUID NOT NULL,
			short_url VARCHAR(255) UNIQUE NOT NULL,
			long_url VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = db.Exec(queryInitShortURLs)
	if err != nil {
		logger.Log.Error("error creating short_urls table", zap.Error(err))
		return nil, err
	}

	queryInitRefreshTokens := `
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			refresh_token VARCHAR(255) UNIQUE NOT NULL,
			user_id UUID NOT NULL,
			expires_at TIMESTAMP NOT NULL
		);
	`
	_, err = db.Exec(queryInitRefreshTokens)
	if err != nil {
		logger.Log.Error("error creating refresh_tokens table", zap.Error(err))
		return nil, err
	}

   queryCreateIndexes := `
		CREATE INDEX idx_short_urls_user_id ON short_urls(user_id);
		CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
		CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
	`
	_, err = db.Exec(queryCreateIndexes)
	if err != nil {
		logger.Log.Error("error creating indexes", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("DB, tables, and indexes are ready")

	logger.Log.Info("DB and tables are ready")

	return &PostgresStorage{
		db: db,
	}, nil
}

func (pg *PostgresStorage) Get(ctx context.Context, origURL string, userID string) (string, error) {
   queryGet := `SELECT long_url FROM short_urls WHERE short_url = $1`
   var longURL string

   // Используем sqlx.QueryRowx, который поддерживает более удобную работу с результатами
   err := pg.db.GetContext(ctx, &longURL, queryGet, origURL)
   if err != nil {
   	logger.Log.Error("short URL not found", zap.String("short_url", origURL))
   	return "", fmt.Errorf("short URL not found")
   }
   return longURL, nil
}
func(pg *PostgresStorage) GetUserURLs(ctx context.Context, userID string) ([]domain.URLData, error) {
   queryGetURLs := `SELECT short_url, long_url FROM short_urls WHERE user_id = $1`
   var urls []domain.URLData
   err := pg.db.GetContext(ctx, &urls, queryGetURLs, userID)
   if err != nil {
      return nil, err
   }
   logger.Log.Debug("taked all urls for user:", userID)
   return urls, nil

}

func (pg *PostgresStorage) Save(ctx context.Context, shortURL string, origURL string, userID string) error {
   querySave := `INSERT INTO short_urls (short_url, long_url, user_id) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING WHERE user_id IN == 3$`
   // Используем ExecContext для выполнения запроса
   _, err := pg.db.ExecContext(ctx, querySave, shortURL, origURL, userID)
   if err != nil {
   	logger.Log.Debug("error saving short url", zap.Error(err))
   	return err
   }
   logger.Log.Debug("short URL saved", zap.String("shortURL", shortURL), zap.String("userID", userID))
   return nil
}

func (pg *PostgresStorage) FindByLongURL(ctx context.Context, shortURL string, userID string) (string, error) {
   select {
   case <-ctx.Done():
   	return "", ctx.Err()
   default:
   }
   queryGet := `SELECT short_url FROM short_urls WHERE user_id = $2 AND short_url = $1`
   var longURL string
   err := pg.db.GetContext(ctx, &longURL, queryGet, shortURL, userID)
   if err != nil {
   	if errors.Is(err,sql.ErrNoRows) {
   		return "", domain.ErrURLNotFound // Если записи нет, возвращаем ошибку "не найдено"
   	}
   	logger.Log.Debugf("error getting short URL for user ID %s", userID, zap.Error(err))
   	return "", err
   }
   return shortURL, nil
}

func(pg *PostgresStorage )GetUserIDByRefreshToken(ctx context.Context, refreshToken string)(string,error){
   select{
   case <-ctx.Done():return "", ctx.Err()
   default:
   }
   query := `SELECT user_id FROM refresh_tokens WHERE refresh_token = $1`
   var userID string
   err := pg.db.GetContext(ctx, &userID, query, refreshToken)
   if err != nil {
      return "", err
   }
   return userID, nil

}
func (pg *PostgresStorage) SaveRefreshToken(ctx context.Context,refreshToken, userID string, expiresAt time.Time)error{
   select{
      case <-ctx.Done():
         return ctx.Err()
      default:   
   }
   querySaveFT := `INSERT INTO refresh_tocens (refresh_token, user_id, expires_at) VALUES ($1, $2, $3)`
   _,err := pg.db.ExecContext(ctx, querySaveFT, refreshToken, userID, expiresAt)
   if err != nil {
      logger.Log.Debug("errors insert refresh token", zap.Error(err))
      return err
   }
   return nil
}

func (pg *PostgresStorage) DeleteRefreshToken(ctx context.Context,refreshToken string)error{
   select{
   case <-ctx.Done():
      return ctx.Err()
   default: 
   }
     


}
func(pg *PostgresStorage)Close() error{}


