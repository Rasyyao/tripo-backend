package repositories

import (
	"errors"
	"sync"
	"time"

	"tripo-backend/internal/models"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepository defines data access operations for users.
// Swap the in-memory implementation for a database-backed one without
// touching the service layer.
type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindAll() ([]*models.User, error)
	FindByID(id string) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id string) error
}

type userRepository struct {
	mu    sync.RWMutex
	store map[string]*models.User
}

func NewUserRepository() UserRepository {
	return &userRepository{
		store: make(map[string]*models.User),
	}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = uuid.NewString()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.store[user.ID] = user
	return user, nil
}

func (r *userRepository) FindAll() ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.store))
	for _, u := range r.store {
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.store[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *userRepository) Update(user *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.store[user.ID]
	if !ok {
		return nil, ErrUserNotFound
	}

	existing.DisplayName = user.DisplayName
	existing.Email = user.Email
	existing.UpdatedAt = time.Now()

	r.store[existing.ID] = existing
	return existing, nil
}

func (r *userRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return ErrUserNotFound
	}
	delete(r.store, id)
	return nil
}
