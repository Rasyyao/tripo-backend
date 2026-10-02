package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"tripo-backend/internal/apperror"
	"tripo-backend/internal/models"
)

// RefreshTokenRepository persists hashed refresh tokens.
type RefreshTokenRepository interface {
	Create(t *models.RefreshToken) error
	FindByHash(hash string) (*models.RefreshToken, error)
	Revoke(id string) (bool, error) // true only if this call revoked it
	RevokeAllForUser(userID string) error
}

type refreshTokenRepository struct{ db *gorm.DB }

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(t *models.RefreshToken) error {
	return r.db.Create(t).Error
}

func (r *refreshTokenRepository) FindByHash(hash string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := r.db.Where("token_hash = ?", hash).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Revoke is atomic: of several concurrent callers only one sees true, which is
// what makes refresh-token rotation safe against replays.
func (r *refreshTokenRepository) Revoke(id string) (bool, error) {
	res := r.db.Model(&models.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", time.Now())
	return res.RowsAffected == 1, res.Error
}

func (r *refreshTokenRepository) RevokeAllForUser(userID string) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", time.Now()).Error
}
