package dto

// CreateGuardianRequest e o corpo esperado no POST de cadastrar responsavel
type CreateGuardianRequest struct {
	Name       string `json:"name"`
	Relation   string `json:"relation"`
	Phone      string `json:"phone"`
	AvatarURL  string `json:"avatar_url"`
	Authorized bool   `json:"authorized"`
}
