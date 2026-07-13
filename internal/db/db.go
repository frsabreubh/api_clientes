package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// Connect opens and validates a PostgreSQL connection pool.
func Connect(databaseURL string) (*sql.DB, error) {
	database, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(25)
	database.SetMaxIdleConns(10)
	database.SetConnMaxLifetime(5 * time.Minute)

	if err := database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}
