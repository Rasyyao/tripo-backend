package models

import "time"

type Place struct {
	ID            string
	Name          string
	Address       *string
	Latitude      float64
	Longitude     float64
	Category      *string
	GooglePlaceID *string
	ImageURL      *string
	Description   *string
	CreatedAt     time.Time
}

type Memory struct {
	ID              string
	ItineraryStopID string
	PhotoURL        string
	Caption         *string
	TakenAt         *time.Time
	CreatedAt       time.Time
}

type Activity struct {
	ID              string
	ItineraryStopID string
	Description     string
	OccurredAt      time.Time
}
