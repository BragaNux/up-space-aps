package handlers

import (
	"encoding/json"
	"net/http"

	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
	"up-espaco/backend/internal/domain/entities"
)

// TurmaHandler cuida das rotas de turmas/salas
type TurmaHandler struct {
	list   *usecases.ListTurmasUseCase
	create *usecases.CreateTurmaUseCase
	update *usecases.UpdateTurmaUseCase
	delete *usecases.DeleteTurmaUseCase
}

func NewTurmaHandler(list *usecases.ListTurmasUseCase, create *usecases.CreateTurmaUseCase, update *usecases.UpdateTurmaUseCase, delete *usecases.DeleteTurmaUseCase) *TurmaHandler {
	return &TurmaHandler{list: list, create: create, update: update, delete: delete}
}

// List devolve todas as turmas (GET /api/turmas)
func (h *TurmaHandler) List(w http.ResponseWriter, r *http.Request) {
	turmas, err := h.list.Execute(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, turmas)
}

// Create cria uma turma nova, so pra contas profissionais (POST /api/turmas)
func (h *TurmaHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireProfissional(r) {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem criar turmas")
		return
	}

	var req dto.CreateTurmaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	turma := &entities.Turma{Name: req.Name}
	if err := h.create.Execute(r.Context(), turma); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, turma)
}

// Update edita uma turma existente, so pra contas profissionais (PUT /api/turmas/{id})
func (h *TurmaHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireProfissional(r) {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem editar turmas")
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID da turma inválido")
		return
	}

	var req dto.CreateTurmaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	turma := &entities.Turma{ID: id, Name: req.Name}
	if err := h.update.Execute(r.Context(), turma); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, turma)
}

// Delete remove uma turma, so pra contas profissionais (DELETE /api/turmas/{id})
func (h *TurmaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireProfissional(r) {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem remover turmas")
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID da turma inválido")
		return
	}

	if err := h.delete.Execute(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// writeJSON serializa qualquer payload como JSON com o status code informado
func (h *TurmaHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError monta uma resposta de erro padrao {"error": "..."}
func (h *TurmaHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
