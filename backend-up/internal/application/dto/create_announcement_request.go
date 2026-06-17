package dto

// CreateAnnouncementRequest e o corpo esperado no POST /api/announcements
type CreateAnnouncementRequest struct {
	Title          string  `json:"title"`
	Sender         string  `json:"sender"`
	Priority       string  `json:"priority"`
	Preview        string  `json:"preview"`
	Body           string  `json:"body"`
	AttachmentName *string `json:"attachment_name,omitempty"`
}
