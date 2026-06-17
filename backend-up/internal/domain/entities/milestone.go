package entities

import "time"

// Milestone e uma conquista/marco do desenvolvimento do aluno (ex: "comeu sozinho", "andou")
type Milestone struct {
	ID          int64      `json:"id"`
	StudentID   int64      `json:"student_id"`
	Title       string     `json:"title"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	AchievedAt  *time.Time `json:"achieved_at,omitempty"`
	Done        bool       `json:"done"`
	CreatedAt   time.Time  `json:"created_at"`
}
