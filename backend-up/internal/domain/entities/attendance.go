package entities

import "time"

// Attendance guarda a presenca/falta de um aluno num dia especifico
type Attendance struct {
	ID             int64      `json:"id"`
	StudentID      int64      `json:"student_id"`
	Date           time.Time  `json:"date"`
	Status         string     `json:"status"` // 'present' or 'absent'
	MarkedByUserID *int64     `json:"marked_by_user_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
