package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	dsn := "host=localhost port=5432 user=appuser password=password dbname=expense_db sslmode=disable"
	return sql.Open("postgres", dsn)
}
