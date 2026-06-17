package handlers

import (
	"encoding/json"
	"net/http"

	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
	"up-espaco/backend/internal/domain/entities"
)

// TimelineHandler cuida das rotas da linha do tempo de cada aluno
type TimelineHandler struct {
	listTimeline *usecases.ListTimelineUseCase
	getEvent     *usecases.GetTimelineEventUseCase
	createEvent  *usecases.CreateTimelineUseCase
	updateEvent  *usecases.UpdateTimelineUseCase
	deleteEvent  *usecases.DeleteTimelineUseCase
}

func NewTimelineHandler(
	listTimeline *usecases.ListTimelineUseCase,
	getEvent *usecases.GetTimelineEventUseCase,
	createEvent *usecases.CreateTimelineUseCase,
	updateEvent *usecases.UpdateTimelineUseCase,
	deleteEvent *usecases.DeleteTimelineUseCase,
) *TimelineHandler {
	return &TimelineHandler{
		listTimeline: listTimeline, getEvent: getEvent, createEvent: createEvent,
		updateEvent: updateEvent, deleteEvent: deleteEvent,
	}
}

// List devolve os eventos de hoje na timeline do aluno (GET /api/timeline?student_id=)
func (h *TimelineHandler) List(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseStudentIDQuery(r)
	if err != nil {
		h.writeJSON(w, http.StatusOK, []any{})
		return
	}

	events, err := h.listTimeline.Execute(r.Context(), studentID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, events)
}

// Get busca um evento especifico da timeline (GET /api/timeline/{id})
func (h *TimelineHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do evento inválido")
		return
	}

	event, err := h.getEvent.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Registro da agenda não encontrado")
		return
	}

	h.writeJSON(w, http.StatusOK, event)
}

// Create cria um evento novo na timeline de um aluno (POST /api/timeline?student_id=)
func (h *TimelineHandler) Create(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseStudentIDQuery(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Parâmetro student_id ausente ou inválido")
		return
	}

	var req dto.CreateTimelineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	event := &entities.TimelineEvent{StudentID: studentID, Title: req.Title, Description: req.Description, OccurredAt: req.OccurredAt}
	if err := h.createEvent.Execute(r.Context(), event); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, event)
}

// Update edita um evento existente da timeline (PUT /api/timeline/{id})
func (h *TimelineHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do evento inválido")
		return
	}

	var req dto.CreateTimelineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	event := &entities.TimelineEvent{ID: id, Title: req.Title, Description: req.Description, OccurredAt: req.OccurredAt}
	if err := h.updateEvent.Execute(r.Context(), event); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, event)
}

// Delete remove um evento da timeline (DELETE /api/timeline/{id})
func (h *TimelineHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do evento inválido")
		return
	}

	if err := h.deleteEvent.Execute(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeJSON serializa qualquer payload como JSON com o status code informado
func (h *TimelineHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError monta uma resposta de erro padrao {"error": "..."}
func (h *TimelineHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
