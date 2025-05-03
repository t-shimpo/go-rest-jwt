package service

import (
	"github.com/t-shimpo/go-rest-standard-library-layered/models"
	"github.com/t-shimpo/go-rest-standard-library-layered/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	return s.repo.GetUserByID(id)
}
