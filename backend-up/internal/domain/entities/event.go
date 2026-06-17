package entities

import "time"

// Event e um evento da agenda da escola (reuniao, festa, passeio etc)
type Event struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	RSVPCount   int64     `json:"rsvp_count"`
	CreatedAt   time.Time `json:"created_at"`
}
