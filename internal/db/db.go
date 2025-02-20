package db

import (
   "database/sql"
   "fmt"
   "local/logger"
   _ "github.com/jackc/pgx/v5/stdlib"
)

type DbConnector struct {
   db *sql.DB
}

func NewDBConnector(dsn string) (*DbConnector, error) {  
   db, err := sql.Open("pgx", dsn)
   if err != nil {
   	return nil, err
   }
   if err := db.Ping(); err != nil {
   	return nil, fmt.Errorf("unable to connect to database: %v", err)
   }

   return &DbConnector{
   	db: db,
   }, nil
}

func (d *DbConnector) PingDataBase() error {
   return d.db.Ping()
}

func (d *DbConnector) CloseDataBase() error {
   if err := d.db.Close(); err != nil {
   	logger.Log.Errorf("error closing database connection: %v", err)
   	return err
   }
   return nil
}
