package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"up-espaco/backend/internal/application/usecases"

	"github.com/gorilla/mux"
)

type EventHandler struct {
	listEvents *usecases.ListEventsUseCase
	rsvpEvent  *usecases.RSVPEventUseCase
}

func NewEventHandler(listEvents *usecases.ListEventsUseCase, rsvpEvent *usecases.RSVPEventUseCase) *EventHandler {
	return &EventHandler{listEvents: listEvents, rsvpEvent: rsvpEvent}
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	events, err := h.listEvents.Execute(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, events)
}

func (h *EventHandler) RSVP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		h.writeError(w, http.StatusBadRequest, "invalid event id")
		return
	}

	rsvpCount, err := h.rsvpEvent.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]int64{"rsvp_count": rsvpCount})
}

func (h *EventHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *EventHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
