package repositories

import (
	"fmt"
	"sync"
	"time"

	"tripo-backend/internal/apperror"
	"tripo-backend/internal/models"

	"github.com/google/uuid"
)

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

var errUserNotFound = fmt.Errorf("user: %w", apperror.ErrNotFound)
var errEmailTaken = fmt.Errorf("email: %w", apperror.ErrConflict)

// emailTaken must be called with the lock held.
func (r *userRepository) emailTaken(email, exceptID string) bool {
	for id, u := range r.store {
		if id != exceptID && u.Email == email {
			return true
		}
	}
	return false
}

// Values are copied on the way in and out so callers can never mutate stored state.
func clone(u *models.User) *models.User {
	c := *u
	return &c
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.emailTaken(user.Email, "") {
		return nil, errEmailTaken
	}

	now := time.Now()
	stored := clone(user)
	stored.ID = uuid.NewString()
	stored.CreatedAt = now
	stored.UpdatedAt = now

	r.store[stored.ID] = stored
	return clone(stored), nil
}

func (r *userRepository) FindAll() ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.store))
	for _, u := range r.store {
		users = append(users, clone(u))
	}
	return users, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.store[id]
	if !ok {
		return nil, errUserNotFound
	}
	return clone(user), nil
}

func (r *userRepository) Update(user *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.store[user.ID]
	if !ok {
		return nil, errUserNotFound
	}
	if r.emailTaken(user.Email, user.ID) {
		return nil, errEmailTaken
	}

	existing.DisplayName = user.DisplayName
	existing.Email = user.Email
	existing.UpdatedAt = time.Now()

	return clone(existing), nil
}

func (r *userRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return errUserNotFound
	}
	delete(r.store, id)
	return nil
}
