package models

import "time"

type ItineraryStatus string

const (
	ItineraryDraft     ItineraryStatus = "draft"
	ItineraryPlanned   ItineraryStatus = "planned"
	ItineraryOngoing   ItineraryStatus = "ongoing"
	ItineraryCompleted ItineraryStatus = "completed"
	ItineraryCancelled ItineraryStatus = "cancelled"
)

type Itinerary struct {
	ID                  string          `json:"id"`
	UserID              string          `json:"user_id"`
	Title               string          `json:"title"`
	StartDate           time.Time       `json:"start_date"`
	EndDate             time.Time       `json:"end_date"`
	TotalBudgetEstimate float64         `json:"total_budget_estimate"`
	Status              ItineraryStatus `json:"status"`
	CoverImageURL       *string         `json:"cover_image_url"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type ItineraryDay struct {
	ID          string    `json:"id"`
	ItineraryID string    `json:"itinerary_id"`
	DayNumber   int       `json:"day_number"`
	DayDate     time.Time `json:"day_date"`
}

type StopStatus string

const (
	StopPending StopStatus = "pending"
	StopEnRoute StopStatus = "en_route"
	StopArrived StopStatus = "arrived"
	StopSkipped StopStatus = "skipped"
)

type ItineraryStop struct {
	ID             string     `json:"id"`
	ItineraryDayID string     `json:"itinerary_day_id"`
	PlaceID        string     `json:"place_id"`
	StopOrder      int        `json:"stop_order"`
	PlannedTime    *string    `json:"planned_time"` // "15:04:05"
	BudgetEstimate float64    `json:"budget_estimate"`
	Status         StopStatus `json:"status"`
	ArrivedAt      *time.Time `json:"arrived_at"`
	DepartedAt     *time.Time `json:"departed_at"`
	Rating         *int16     `json:"rating"`
	ActualPrice    *float64   `json:"actual_price"`
	CreatedAt      time.Time  `json:"created_at"`
}

type RouteLeg struct {
	ID              string    `json:"id"`
	ItineraryDayID  string    `json:"itinerary_day_id"`
	FromStopID      string    `json:"from_stop_id"`
	ToStopID        string    `json:"to_stop_id"`
	DistanceMeters  *int      `json:"distance_meters"`
	DurationSeconds *int      `json:"duration_seconds"`
	Polyline        *string   `json:"polyline"`
	ComputedAt      time.Time `json:"computed_at"`
}

type ItineraryShare struct {
	ID            string    `json:"id"`
	ItineraryID   string    `json:"itinerary_id"`
	ShareImageURL *string   `json:"share_image_url"`
	ShareSlug     *string   `json:"share_slug"`
	Platform      *string   `json:"platform"`
	CreatedAt     time.Time `json:"created_at"`
}
