package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/t-shimpo/go-rest-standard-library-layered/config"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func GetUsers(limit, offset int) ([]User, error) {
	query := `SELECT id, name, email, created_at FROM users ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := config.DB.Query(query, limit, offset)
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

func UpdateUser(id int, name, email *string) (*User, error) {
	updates := map[string]interface{}{}

	if name != nil {
		updates["name"] = *name
	}

	if email != nil {
		updates["email"] = *email
	}

	// SET句を動的に構築
	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d RETURNING id, name, email, created_at",
		strings.Join(setClauses, ", "), argIndex,
	)
	args = append(args, id)

	user := &User{}
	err := config.DB.QueryRow(query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sql: no rows in result set")
		}
		return nil, fmt.Errorf("ユーザーの更新に失敗しました: %w", err)
	}

	return user, nil
}

func DeleteUser(id int) error {
	result, err := config.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("削除に失敗しました: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("削除結果の確認に失敗しました: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("sql: no rows in result set")
	}

	return nil
}
