package models

import "time"

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	DisplayName  *string   `gorm:"size:100"`
	AvatarURL    *string   `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

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
