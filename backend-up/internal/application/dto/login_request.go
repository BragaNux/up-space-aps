package dto

// LoginRequest e o corpo esperado no POST /api/auth/login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse e o que devolvemos depois de um login valido (token + dados do usuario)
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
