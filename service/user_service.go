package service

import (
	"errors"

	"github.com/t-shimpo/go-rest-jwt/models"
	"github.com/t-shimpo/go-rest-jwt/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrValidation        = errors.New("validation error")
	ErrorNotFound        = errors.New("not found")
	ErrorInvalidPassword = errors.New("invalid password")
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *models.User, password string) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, ErrValidation
	}
	if password == "" {
		return nil, ErrValidation
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = hashedPassword

	return s.repo.CreateUser(user)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) GetUsers(limit, offset int) ([]*models.User, error) {
	return s.repo.GetUsers(limit, offset)
}

func (s *UserService) PatchUser(id int64, name, email *string) (*models.User, error) {
	return s.repo.PatchUser(id, name, email)
}

func (s *UserService) DeleteUser(id int64) error {
	return s.repo.DeleteUser(id)
}

func (s *UserService) Authenticate(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, ErrorNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrorInvalidPassword
	}
	return user, nil
}
