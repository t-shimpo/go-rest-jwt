package models

import (
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
