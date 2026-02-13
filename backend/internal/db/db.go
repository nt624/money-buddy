package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=appuser password=password dbname=expense_db sslmode=disable"
	}
	return sql.Open("postgres", dsn)
}
