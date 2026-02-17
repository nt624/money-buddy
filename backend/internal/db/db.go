package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=appuser password=password dbname=expense_db sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// コネクションプール設定（Neon対応）
	db.SetMaxOpenConns(10)                 // 最大同時接続数（Neonの制限内）
	db.SetMaxIdleConns(5)                  // アイドル接続数
	db.SetConnMaxLifetime(5 * time.Minute) // 接続の最大寿命（Neonのタイムアウト対策）
	db.SetConnMaxIdleTime(2 * time.Minute) // アイドル接続の最大時間

	// 接続テスト
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Println("Database connection established and tested successfully")
	return db, nil
}
