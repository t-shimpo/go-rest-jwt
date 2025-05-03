package repository

import (
	"database/sql"

	"github.com/t-shimpo/go-rest-standard-library-layered/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRow(query, user.Name, user.Email).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
