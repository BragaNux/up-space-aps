package handlers

import (
	"encoding/json"
	"net/http"

	"up-espaco/backend/internal/application/usecases"
)

type TimelineHandler struct {
	listTimeline *usecases.ListTimelineUseCase
}

func NewTimelineHandler(listTimeline *usecases.ListTimelineUseCase) *TimelineHandler {
	return &TimelineHandler{listTimeline: listTimeline}
}

func (h *TimelineHandler) List(w http.ResponseWriter, r *http.Request) {
	events, err := h.listTimeline.Execute(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, events)
}

func (h *TimelineHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *TimelineHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
