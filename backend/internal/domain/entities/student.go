package entities

import "time"

type Student struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	PresenceStatus string     `json:"presence_status"`
	CheckInAt      *time.Time `json:"check_in_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
