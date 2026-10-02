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
