package entities

// Guardian e um responsavel autorizado a buscar/acompanhar um aluno
type Guardian struct {
	ID         int64  `json:"id"`
	StudentID  int64  `json:"student_id"`
	Name       string `json:"name"`
	Relation   string `json:"relation"`
	Phone      string `json:"phone"`
	AvatarURL  string `json:"avatar_url"`
	Authorized bool   `json:"authorized"`
}
