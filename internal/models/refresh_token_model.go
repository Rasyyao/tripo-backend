package models

import "time"

// RefreshToken stores only a hash of the token, never the raw value.
type RefreshToken struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"type:uuid;not null;index"`
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
	TokenHash string    `gorm:"type:text;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time `gorm:"not null"`
}
