package entities

import "time"

// Student e o cadastro completo de um aluno (dados pessoais, saude, turma e responsaveis)
type Student struct {
	ID             int64       `json:"id"`
	Name           string      `json:"name"`
	PresenceStatus string      `json:"presence_status"`
	CheckInAt      *time.Time  `json:"check_in_at,omitempty"`
	GuardianUserID *int64      `json:"guardian_user_id,omitempty"`
	PhotoURL       string      `json:"photo_url"`
	TurmaID        *int64      `json:"turma_id,omitempty"`
	GroupName      string      `json:"group_name"`
	TeacherUserID  *int64      `json:"teacher_user_id,omitempty"`
	TeacherName    string      `json:"teacher_name"`
	BirthDate      *time.Time  `json:"birth_date,omitempty"`
	EnrollmentCode string      `json:"enrollment_code"`
	BloodType      string      `json:"blood_type"`
	Allergies      []string    `json:"allergies"`
	Restrictions   string      `json:"restrictions"`
	Medications    string      `json:"medications"`
	Guardians      []*Guardian `json:"guardians,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
