package services

import (
	"strings"

	"golang.org/x/crypto/bcrypt"

	"tripo-backend/internal/dto"
	"tripo-backend/internal/models"
	"tripo-backend/internal/repositories"
)

type UserService interface {
	CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
	ListUsers() ([]*dto.UserResponse, error)
	GetUser(id string) (*dto.UserResponse, error)
	UpdateUser(id string, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := normalizeEmail(req.Email)

	if err := collect(validateName(name), validateEmail(email), validatePassword(req.Password)); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.Create(&models.User{
		DisplayName:  &name,
		Email:        email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponse(user), nil
}

func (s *userService) ListUsers() ([]*dto.UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponses(users), nil
}

func (s *userService) GetUser(id string) (*dto.UserResponse, error) {
	if err := collect(validateID(id)); err != nil {
		return nil, err
	}
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponse(user), nil
}

func (s *userService) UpdateUser(id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := normalizeEmail(req.Email)

	if err := collect(validateID(id), validateName(name), validateEmail(email)); err != nil {
		return nil, err
	}

	user, err := s.repo.Update(&models.User{
		ID:          id,
		DisplayName: &name,
		Email:       email,
	})
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponse(user), nil
}

func (s *userService) DeleteUser(id string) error {
	if err := collect(validateID(id)); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
