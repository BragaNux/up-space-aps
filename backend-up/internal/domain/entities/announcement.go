package entities

import "time"

// Announcement e um comunicado/aviso enviado pra comunidade da escola (recados, avisos importantes etc)
type Announcement struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	Sender         string    `json:"sender"`
	Priority       string    `json:"priority"`
	Preview        string    `json:"preview"`
	Body           string    `json:"body"`
	AttachmentName *string   `json:"attachment_name,omitempty"`
	Read           bool      `json:"read"`
	CreatedAt      time.Time `json:"created_at"`
}
