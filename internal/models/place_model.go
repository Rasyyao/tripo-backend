package models

import "time"

type Place struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string    `gorm:"size:200;not null"`
	Address       *string   `gorm:"type:text"`
	Latitude      float64   `gorm:"not null;index:idx_places_coords,priority:1"`
	Longitude     float64   `gorm:"not null;index:idx_places_coords,priority:2"`
	Category      *string   `gorm:"size:80"`
	GooglePlaceID *string   `gorm:"size:150;uniqueIndex"`
	ImageURL      *string   `gorm:"type:text"`
	Description   *string   `gorm:"type:text"`
	CreatedAt     time.Time `gorm:"not null"`
}

type Memory struct {
	ID              string        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItineraryStopID string        `gorm:"type:uuid;not null;index"`
	ItineraryStop   ItineraryStop `gorm:"constraint:OnDelete:CASCADE"`
	PhotoURL        string        `gorm:"type:text;not null"`
	Caption         *string       `gorm:"size:300"`
	TakenAt         *time.Time
	CreatedAt       time.Time `gorm:"not null"`
}

type Activity struct {
	ID              string        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItineraryStopID string        `gorm:"type:uuid;not null;index"`
	ItineraryStop   ItineraryStop `gorm:"constraint:OnDelete:CASCADE"`
	Description     string        `gorm:"size:300;not null"`
	OccurredAt      time.Time     `gorm:"not null;default:now()"`
}
