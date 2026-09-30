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
	ID                  string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID              string          `gorm:"type:uuid;not null;index:idx_itineraries_status,priority:1"`
	User                User            `gorm:"constraint:OnDelete:CASCADE"`
	Title               string          `gorm:"size:150;not null"`
	StartDate           time.Time       `gorm:"type:date;not null"`
	EndDate             time.Time       `gorm:"type:date;not null"`
	TotalBudgetEstimate float64         `gorm:"type:numeric(12,2);default:0"`
	Status              ItineraryStatus `gorm:"size:20;not null;default:draft;index:idx_itineraries_status,priority:2"`
	CoverImageURL       *string         `gorm:"type:text"`
	CreatedAt           time.Time       `gorm:"not null"`
	UpdatedAt           time.Time       `gorm:"not null"`
}

type ItineraryDay struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItineraryID string    `gorm:"type:uuid;not null;uniqueIndex:idx_itinerary_day_number"`
	Itinerary   Itinerary `gorm:"constraint:OnDelete:CASCADE"`
	DayNumber   int       `gorm:"not null;uniqueIndex:idx_itinerary_day_number"`
	DayDate     time.Time `gorm:"type:date;not null"`
}

type StopStatus string

const (
	StopPending StopStatus = "pending"
	StopEnRoute StopStatus = "en_route"
	StopArrived StopStatus = "arrived"
	StopSkipped StopStatus = "skipped"
)

type ItineraryStop struct {
	ID             string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItineraryDayID string       `gorm:"type:uuid;not null;uniqueIndex:idx_stop_day_order"`
	ItineraryDay   ItineraryDay `gorm:"constraint:OnDelete:CASCADE"`
	PlaceID        string       `gorm:"type:uuid;not null"`
	Place          Place
	StopOrder      int        `gorm:"not null;uniqueIndex:idx_stop_day_order"`
	PlannedTime    *string    `gorm:"type:time"` // "15:04:05"
	BudgetEstimate float64    `gorm:"type:numeric(10,2);default:0"`
	Status         StopStatus `gorm:"size:20;not null;default:pending"`
	ArrivedAt      *time.Time
	DepartedAt     *time.Time
	Rating         *int16    `gorm:"check:rating BETWEEN 1 AND 5"`
	ActualPrice    *float64  `gorm:"type:numeric(10,2)"`
	CreatedAt      time.Time `gorm:"not null"`
}

type RouteLeg struct {
	ID              string        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItineraryDayID  string        `gorm:"type:uuid;not null;index"`
	ItineraryDay    ItineraryDay  `gorm:"constraint:OnDelete:CASCADE"`
	FromStopID      string        `gorm:"type:uuid;not null"`
	FromStop        ItineraryStop `gorm:"foreignKey:FromStopID"`
	ToStopID        string        `gorm:"type:uuid;not null"`
	ToStop          ItineraryStop `gorm:"foreignKey:ToStopID"`
	DistanceMeters  *int
	DurationSeconds *int
	Polyline        *string   `gorm:"type:text"`
	ComputedAt      time.Time `gorm:"not null;default:now()"`
}

type ItineraryShare struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItineraryID   string    `gorm:"type:uuid;not null;index"`
	Itinerary     Itinerary `gorm:"constraint:OnDelete:CASCADE"`
	ShareImageURL *string   `gorm:"type:text"`
	ShareSlug     *string   `gorm:"size:50;uniqueIndex"`
	Platform      *string   `gorm:"size:50"`
	CreatedAt     time.Time `gorm:"not null"`
}
