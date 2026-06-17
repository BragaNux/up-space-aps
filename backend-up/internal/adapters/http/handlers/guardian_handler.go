package handlers

import (
	"encoding/json"
	"net/http"

	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
	"up-espaco/backend/internal/domain/entities"
)

// GuardianHandler cuida das rotas dos responsaveis autorizados por aluno
type GuardianHandler struct {
	list   *usecases.ListGuardiansUseCase
	create *usecases.CreateGuardianUseCase
	update *usecases.UpdateGuardianUseCase
	delete *usecases.DeleteGuardianUseCase
}

func NewGuardianHandler(list *usecases.ListGuardiansUseCase, create *usecases.CreateGuardianUseCase, update *usecases.UpdateGuardianUseCase, delete *usecases.DeleteGuardianUseCase) *GuardianHandler {
	return &GuardianHandler{list: list, create: create, update: update, delete: delete}
}

// List devolve os responsaveis de um aluno (GET /api/students/{id}/guardians)
func (h *GuardianHandler) List(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do aluno inválido")
		return
	}

	guardians, err := h.list.Execute(r.Context(), studentID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, guardians)
}

// Create cadastra um responsavel novo pra um aluno (POST /api/students/{id}/guardians)
func (h *GuardianHandler) Create(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do aluno inválido")
		return
	}

	var req dto.CreateGuardianRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	guardian := &entities.Guardian{
		StudentID: studentID, Name: req.Name, Relation: req.Relation,
		Phone: req.Phone, AvatarURL: req.AvatarURL, Authorized: req.Authorized,
	}

	if err := h.create.Execute(r.Context(), guardian); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, guardian)
}

// Update edita os dados de um responsavel (PUT /api/guardians/{id})
func (h *GuardianHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do contato inválido")
		return
	}

	var req dto.CreateGuardianRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	guardian := &entities.Guardian{
		ID: id, Name: req.Name, Relation: req.Relation,
		Phone: req.Phone, AvatarURL: req.AvatarURL, Authorized: req.Authorized,
	}

	if err := h.update.Execute(r.Context(), guardian); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, guardian)
}

// Delete remove um responsavel (DELETE /api/guardians/{id})
func (h *GuardianHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do contato inválido")
		return
	}

	if err := h.delete.Execute(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeJSON serializa qualquer payload como JSON com o status code informado
func (h *GuardianHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError monta uma resposta de erro padrao {"error": "..."}
func (h *GuardianHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
