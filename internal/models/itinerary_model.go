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
	ID                  string
	UserID              string
	Title               string
	StartDate           time.Time
	EndDate             time.Time
	TotalBudgetEstimate float64
	Status              ItineraryStatus
	CoverImageURL       *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type ItineraryDay struct {
	ID          string
	ItineraryID string
	DayNumber   int
	DayDate     time.Time
}

type StopStatus string

const (
	StopPending StopStatus = "pending"
	StopEnRoute StopStatus = "en_route"
	StopArrived StopStatus = "arrived"
	StopSkipped StopStatus = "skipped"
)

type ItineraryStop struct {
	ID             string
	ItineraryDayID string
	PlaceID        string
	StopOrder      int
	PlannedTime    *string // "15:04:05"
	BudgetEstimate float64
	Status         StopStatus
	ArrivedAt      *time.Time
	DepartedAt     *time.Time
	Rating         *int16
	ActualPrice    *float64
	CreatedAt      time.Time
}

type RouteLeg struct {
	ID              string
	ItineraryDayID  string
	FromStopID      string
	ToStopID        string
	DistanceMeters  *int
	DurationSeconds *int
	Polyline        *string
	ComputedAt      time.Time
}

type ItineraryShare struct {
	ID            string
	ItineraryID   string
	ShareImageURL *string
	ShareSlug     *string
	Platform      *string
	CreatedAt     time.Time
}
