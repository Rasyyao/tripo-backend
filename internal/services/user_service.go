package services

import (
	"errors"

	"tripo-backend/internal/models"
	"tripo-backend/internal/repositories"
)

var ErrInvalidInput = errors.New("invalid input")

// UserService holds business logic and orchestrates repository calls.
type UserService interface {
	CreateUser(name, email string) (*models.User, error)
	ListUsers() ([]*models.User, error)
	GetUser(id string) (*models.User, error)
	UpdateUser(id, name, email string) (*models.User, error)
	DeleteUser(id string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(name, email string) (*models.User, error) {
	if name == "" || email == "" {
		return nil, ErrInvalidInput
	}

	user := &models.User{
		DisplayName: &name,
		Email:       email,
	}

	return s.repo.Create(user)
}

func (s *userService) ListUsers() ([]*models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUser(id string) (*models.User, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.FindByID(id)
}

func (s *userService) UpdateUser(id, name, email string) (*models.User, error) {
	if id == "" || name == "" || email == "" {
		return nil, ErrInvalidInput
	}

	user := &models.User{
		ID:          id,
		DisplayName: &name,
		Email:       email,
	}

	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id string) error {
	if id == "" {
		return ErrInvalidInput
	}
	return s.repo.Delete(id)
}
