package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// エラーを防ぐため一時的に残す
// 全ての処理をレイヤー分離すれば不要
var DB *sql.DB

func InitDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("データベースURLが設定されていません")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("データベース接続に失敗しました: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベースの接続確認に失敗しました: %w", err)
	}

	if err := applyMigrations(db); err != nil {
		return nil, fmt.Errorf("マイグレーション適用に失敗しました: %w", err)
	}

	return db, nil
}

func applyMigrations(db *sql.DB) error {
	createUsersTables := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createUsersTables)
	if err != nil {
		return err
	}

	fmt.Println("マイグレーションの適用に成功しました")

	return nil
}
