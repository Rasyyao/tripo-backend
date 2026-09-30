package models

import "time"

type Place struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Address       *string   `json:"address"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Category      *string   `json:"category"`
	GooglePlaceID *string   `json:"google_place_id"`
	ImageURL      *string   `json:"image_url"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

type Memory struct {
	ID              string     `json:"id"`
	ItineraryStopID string     `json:"itinerary_stop_id"`
	PhotoURL        string     `json:"photo_url"`
	Caption         *string    `json:"caption"`
	TakenAt         *time.Time `json:"taken_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Activity struct {
	ID              string    `json:"id"`
	ItineraryStopID string    `json:"itinerary_stop_id"`
	Description     string    `json:"description"`
	OccurredAt      time.Time `json:"occurred_at"`
}
