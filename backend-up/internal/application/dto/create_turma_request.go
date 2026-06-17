package dto

// CreateTurmaRequest e o corpo esperado no POST /api/turmas
type CreateTurmaRequest struct {
	Name string `json:"name"`
}
