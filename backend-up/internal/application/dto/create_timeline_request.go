package dto

import "time"

// CreateTimelineRequest e o corpo esperado no POST /api/timeline
type CreateTimelineRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OccurredAt  time.Time `json:"occurred_at"`
}
