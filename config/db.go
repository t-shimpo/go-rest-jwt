package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("データベースURLが設定されていません")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("データベース接続に失敗しました: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("データベースの接続確認に失敗しました: %w", err)
	}

	DB = db

	if err := applyMigrations(); err != nil {
		return fmt.Errorf("マイグレーション適用に失敗しました: %w", err)
	}

	return nil
}

func applyMigrations() error {
	createUsersTables := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(createUsersTables)
	if err != nil {
		return err
	}

	fmt.Println("マイグレーションの適用に成功しました")

	return nil
}
