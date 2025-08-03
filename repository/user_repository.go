package repository

import (
	"database/sql"
	"fmt"

	"github.com/t-shimpo/go-rest-jwt/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUsers(limit, offset int) ([]*models.User, error)
	PatchUser(id int, name, email *string) (*models.User, error)
	DeleteUser(id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(query, user.Name, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, name, email, created_at FROM users where id = $1`
	var user models.User

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 呼び出し元で 404 として扱う
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(Email string) (*models.User, error) {
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`
	var user models.User

	err := r.db.QueryRow(query, Email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUsers(limit, offset int) ([]*models.User, error) {
	query := `SELECT id, name, email, created_at FROM users ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) PatchUser(id int, name, email *string) (*models.User, error) {
	// 既存データを取得
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	var user models.User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	// 更新内容を反映
	if name != nil {
		user.Name = *name
	}
	if email != nil {
		user.Email = *email
	}

	// 更新クエリ
	query = `UPDATE users SET name = $1, email = $2 WHERE id = $3 RETURNING id, name, email, created_at`
	err = r.db.QueryRow(query, user.Name, user.Email, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 呼び出し元で 404 として扱う
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) DeleteUser(id int) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
