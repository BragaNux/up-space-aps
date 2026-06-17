package handlers

import (
	"encoding/json"
	"net/http"

	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
	"up-espaco/backend/internal/domain/entities"
)

// MilestoneHandler cuida das rotas de marcos de desenvolvimento dos alunos
type MilestoneHandler struct {
	list   *usecases.ListMilestonesUseCase
	create *usecases.CreateMilestoneUseCase
	update *usecases.UpdateMilestoneUseCase
	delete *usecases.DeleteMilestoneUseCase
}

func NewMilestoneHandler(list *usecases.ListMilestonesUseCase, create *usecases.CreateMilestoneUseCase, update *usecases.UpdateMilestoneUseCase, delete *usecases.DeleteMilestoneUseCase) *MilestoneHandler {
	return &MilestoneHandler{list: list, create: create, update: update, delete: delete}
}

// List devolve os marcos de um aluno, ou lista vazia se nao informar o aluno (GET /api/milestones?student_id=)
func (h *MilestoneHandler) List(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseStudentIDQuery(r)
	if err != nil {
		h.writeJSON(w, http.StatusOK, []any{})
		return
	}

	milestones, err := h.list.Execute(r.Context(), studentID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, milestones)
}

// Create registra um marco novo pra um aluno (POST /api/milestones)
func (h *MilestoneHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateMilestoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	milestone := &entities.Milestone{
		StudentID: req.StudentID, Title: req.Title, Category: req.Category,
		Description: req.Description, AchievedAt: req.AchievedAt, Done: req.Done,
	}

	if err := h.create.Execute(r.Context(), milestone); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, milestone)
}

// Update edita um marco existente (PUT /api/milestones/{id})
func (h *MilestoneHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do marco inválido")
		return
	}

	var req dto.CreateMilestoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	milestone := &entities.Milestone{
		ID: id, StudentID: req.StudentID, Title: req.Title, Category: req.Category,
		Description: req.Description, AchievedAt: req.AchievedAt, Done: req.Done,
	}

	if err := h.update.Execute(r.Context(), milestone); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, milestone)
}

// Delete remove um marco (DELETE /api/milestones/{id})
func (h *MilestoneHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do marco inválido")
		return
	}

	if err := h.delete.Execute(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeJSON serializa qualquer payload como JSON com o status code informado
func (h *MilestoneHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError monta uma resposta de erro padrao {"error": "..."}
func (h *MilestoneHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
