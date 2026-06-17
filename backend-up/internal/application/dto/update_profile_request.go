package dto

// UpdateProfileRequest e o corpo esperado no PUT /api/me
type UpdateProfileRequest struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	AvatarURL string `json:"avatar_url"`
}
