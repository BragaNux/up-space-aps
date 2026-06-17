package dto

import "time"

// CreateMilestoneRequest e o corpo esperado no POST /api/milestones
type CreateMilestoneRequest struct {
	StudentID   int64      `json:"student_id"`
	Title       string     `json:"title"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	AchievedAt  *time.Time `json:"achieved_at,omitempty"`
	Done        bool       `json:"done"`
}
