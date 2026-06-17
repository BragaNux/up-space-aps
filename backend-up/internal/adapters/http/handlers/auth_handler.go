package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"up-espaco/backend/internal/adapters/http/middleware"
	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
)

// AuthHandler cuida das rotas de autenticacao e perfil (register, login, /me, etc)
type AuthHandler struct {
	register       *usecases.RegisterUserUseCase
	login          *usecases.LoginUseCase
	forgotPassword *usecases.ForgotPasswordUseCase
	getMe          *usecases.GetMeUseCase
	updateProfile  *usecases.UpdateProfileUseCase
	listByRole     *usecases.ListUsersByRoleUseCase
}

func NewAuthHandler(
	register *usecases.RegisterUserUseCase,
	login *usecases.LoginUseCase,
	forgotPassword *usecases.ForgotPasswordUseCase,
	getMe *usecases.GetMeUseCase,
	updateProfile *usecases.UpdateProfileUseCase,
	listByRole *usecases.ListUsersByRoleUseCase,
) *AuthHandler {
	return &AuthHandler{
		register: register, login: login, forgotPassword: forgotPassword,
		getMe: getMe, updateProfile: updateProfile, listByRole: listByRole,
	}
}

// Register cria uma conta nova (POST /api/auth/register)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	user, err := h.register.Execute(r.Context(), req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		if errors.Is(err, usecases.ErrEmailAlreadyInUse) {
			h.writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, user)
}

// Login confere email/senha e devolve o token JWT (POST /api/auth/login)
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	token, user, err := h.login.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "E-mail ou senha inválidos")
		return
	}

	h.writeJSON(w, http.StatusOK, dto.LoginResponse{Token: token, User: user})
}

// ForgotPassword sempre responde "ok" pra nao revelar se o email existe (POST /api/auth/forgot-password)
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	_ = h.forgotPassword.Execute(r.Context(), req.Email)
	h.writeJSON(w, http.StatusOK, map[string]string{"Status": "se o e-mail existir, um link de recuperação foi enviado"})
}

// Me devolve os dados do usuario logado (GET /api/me)
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Autenticação necessária")
		return
	}

	user, err := h.getMe.Execute(r.Context(), userID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Usuário não encontrado")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// UpdateMe atualiza o perfil do usuario logado (PUT /api/me)
func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Autenticação necessária")
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	user, err := h.updateProfile.Execute(r.Context(), userID, req.Name, req.Phone, req.Address, req.AvatarURL)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// ListByRole lista usuarios por papel, so acessivel pra contas profissionais (GET /api/users)
func (h *AuthHandler) ListByRole(w http.ResponseWriter, r *http.Request) {
	role, ok := middleware.RoleFromContext(r.Context())
	if !ok || role != "profissional" {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem listar usuários")
		return
	}

	targetRole := r.URL.Query().Get("role")
	users, err := h.listByRole.Execute(r.Context(), targetRole)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, users)
}

// writeJSON serializa qualquer payload como JSON com o status code informado
func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError monta uma resposta de erro padrao {"error": "..."}
func (h *AuthHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
