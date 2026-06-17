package dto

// PresenceRequest e o corpo esperado no PATCH de marcar presenca do aluno
type PresenceRequest struct {
	Status string `json:"status"`
}
