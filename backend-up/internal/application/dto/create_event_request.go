package dto

import "time"

// CreateEventRequest e o corpo esperado no POST /api/events
type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
}
