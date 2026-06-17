package dto

// ForgotPasswordRequest e o corpo esperado no POST /api/auth/forgot-password
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}
