package handlers

import (
	"encoding/json"
	"net/http"

	"up-espaco/backend/internal/adapters/http/middleware"
	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
	"up-espaco/backend/internal/domain/entities"
)

// EventHandler cuida das rotas da agenda de eventos
type EventHandler struct {
	listEvents  *usecases.ListEventsUseCase
	getEvent    *usecases.GetEventUseCase
	createEvent *usecases.CreateEventUseCase
	updateEvent *usecases.UpdateEventUseCase
	deleteEvent *usecases.DeleteEventUseCase
	rsvpEvent   *usecases.RSVPEventUseCase
}

func NewEventHandler(
	listEvents *usecases.ListEventsUseCase,
	getEvent *usecases.GetEventUseCase,
	createEvent *usecases.CreateEventUseCase,
	updateEvent *usecases.UpdateEventUseCase,
	deleteEvent *usecases.DeleteEventUseCase,
	rsvpEvent *usecases.RSVPEventUseCase,
) *EventHandler {
	return &EventHandler{
		listEvents: listEvents, getEvent: getEvent, createEvent: createEvent,
		updateEvent: updateEvent, deleteEvent: deleteEvent, rsvpEvent: rsvpEvent,
	}
}

// List devolve os eventos da agenda, marcando os RSVPs do usuario logado se houver (GET /api/events)
func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	var userID *int64
	if id, ok := middleware.UserIDFromContext(r.Context()); ok {
		userID = &id
	}

	events, err := h.listEvents.Execute(r.Context(), userID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, events)
}

// Get busca um evento especifico (GET /api/events/{id})
func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do evento inválido")
		return
	}

	event, err := h.getEvent.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Evento não encontrado")
		return
	}

	h.writeJSON(w, http.StatusOK, event)
}

// Create cria um evento novo, so pra contas profissionais (POST /api/events)
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireProfissional(r) {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem criar eventos")
		return
	}

	var req dto.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	event := &entities.Event{
		Title: req.Title, Description: req.Description, Location: req.Location,
		StartsAt: req.StartsAt, EndsAt: req.EndsAt,
	}

	if err := h.createEvent.Execute(r.Context(), event); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, event)
}

// Update edita um evento existente, so pra contas profissionais (PUT /api/events/{id})
func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireProfissional(r) {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem editar eventos")
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do evento inválido")
		return
	}

	var req dto.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	event := &entities.Event{
		ID: id, Title: req.Title, Description: req.Description, Location: req.Location,
		StartsAt: req.StartsAt, EndsAt: req.EndsAt,
	}

	if err := h.updateEvent.Execute(r.Context(), event); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, event)
}

// Delete remove um evento, so pra contas profissionais (DELETE /api/events/{id})
func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireProfissional(r) {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem remover eventos")
		return
	}

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

// RSVP confirma ou cancela a presenca do usuario logado no evento (POST /api/events/{id}/rsvp)
func (h *EventHandler) RSVP(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do evento inválido")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Autenticação necessária")
		return
	}

	count, confirmed, err := h.rsvpEvent.Execute(r.Context(), id, userID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]any{"rsvp_count": count, "confirmed": confirmed})
}

// writeJSON serializa qualquer payload como JSON com o status code informado
func (h *EventHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError monta uma resposta de erro padrao {"error": "..."}
func (h *EventHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
