package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/t-shimpo/go-rest-standard-library/config"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateUser(name, email string) (*User, error) {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at`
	user := &User{}

	err := config.DB.QueryRow(query, name, email).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("ユーザーの作成に失敗しました: %w", err)
	}

	return user, nil
}

func GetUserByID(id int) (*User, error) {
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	user := &User{}

	err := config.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sql: no rows in result set") // ユーザーが見つからない場合
		}
		return nil, fmt.Errorf("ユーザーの取得に失敗しました: %w", err) // その他エラー
	}

	return user, nil
}

func GetUsers() ([]User, error) {
	query := `SELECT id, name, email, created_at FROM users`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ユーザー一覧の取得に失敗しました: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("ユーザーのデータ取得中にエラーが発生しました: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ユーザー一覧の取得中にエラーが発生しました: %w", err)
	}

	return users, nil
}
