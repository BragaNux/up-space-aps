package handlers

import (
	"encoding/json"
	"net/http"

	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
)

type StudentHandler struct {
	getStudent     *usecases.GetStudentUseCase
	updatePresence *usecases.UpdatePresenceUseCase
}

func NewStudentHandler(getStudent *usecases.GetStudentUseCase, updatePresence *usecases.UpdatePresenceUseCase) *StudentHandler {
	return &StudentHandler{getStudent: getStudent, updatePresence: updatePresence}
}

func (h *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	student, err := h.getStudent.Execute(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, student)
}

func (h *StudentHandler) UpdatePresence(w http.ResponseWriter, r *http.Request) {
	var req dto.PresenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.updatePresence.Execute(r.Context(), req.Status); err != nil {
		if err.Error() == "invalid presence status" {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *StudentHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *StudentHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
