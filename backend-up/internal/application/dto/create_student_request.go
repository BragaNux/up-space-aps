package dto

import "time"

// CreateStudentRequest e o corpo esperado no POST /api/students
type CreateStudentRequest struct {
	Name           string     `json:"name"`
	GuardianUserID *int64     `json:"guardian_user_id,omitempty"`
	GuardianEmail  string     `json:"guardian_email,omitempty"`
	PhotoURL       string     `json:"photo_url"`
	TurmaID        *int64     `json:"turma_id,omitempty"`
	GroupName      string     `json:"group_name"`
	TeacherUserID  *int64     `json:"teacher_user_id,omitempty"`
	TeacherName    string     `json:"teacher_name"`
	BirthDate      *time.Time `json:"birth_date,omitempty"`
	EnrollmentCode string     `json:"enrollment_code"`
	BloodType      string     `json:"blood_type"`
	Allergies      []string   `json:"allergies"`
	Restrictions   string     `json:"restrictions"`
	Medications    string     `json:"medications"`
}
