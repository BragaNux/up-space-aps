package dto

// RegisterRequest e o corpo esperado no POST /api/auth/register
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
