package repositories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"tripo-backend/internal/apperror"
	"tripo-backend/internal/models"
)

// UserRepository defines data access operations for users.
type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindAll() ([]*models.User, error)
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

var errUserNotFound = fmt.Errorf("user: %w", apperror.ErrNotFound)
var errEmailTaken = fmt.Errorf("email: %w", apperror.ErrConflict)

// mapErr translates GORM errors into layer-independent apperror values.
func mapErr(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errUserNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return errEmailTaken
	default:
		return err
	}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, mapErr(err)
	}
	return user, nil
}

func (r *userRepository) FindAll() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, mapErr(err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, mapErr(err)
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) (*models.User, error) {
	existing, err := r.FindByID(user.ID)
	if err != nil {
		return nil, err
	}

	existing.DisplayName = user.DisplayName
	existing.Email = user.Email
	if err := r.db.Save(existing).Error; err != nil {
		return nil, mapErr(err)
	}
	return existing, nil
}

func (r *userRepository) Delete(id string) error {
	res := r.db.Delete(&models.User{}, "id = ?", id)
	if res.Error != nil {
		return mapErr(res.Error)
	}
	if res.RowsAffected == 0 {
		return errUserNotFound
	}
	return nil
}
