package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"up-espaco/backend/internal/adapters/http/middleware"
	"up-espaco/backend/internal/application/usecases"
	"up-espaco/backend/internal/domain/entities"
)

// AttendanceHandler cuida das rotas de presenca/falta dos alunos
type AttendanceHandler struct {
	saveUseCase          *usecases.SaveAttendanceUseCase
	getUseCase           *usecases.GetTurmaAttendanceUseCase
	listByStudentUseCase *usecases.ListAttendanceByStudentUseCase
}

func NewAttendanceHandler(
	saveUseCase *usecases.SaveAttendanceUseCase,
	getUseCase *usecases.GetTurmaAttendanceUseCase,
	listByStudentUseCase *usecases.ListAttendanceByStudentUseCase,
) *AttendanceHandler {
	return &AttendanceHandler{
		saveUseCase:          saveUseCase,
		getUseCase:           getUseCase,
		listByStudentUseCase: listByStudentUseCase,
	}
}

// Save registra a presenca/falta de um aluno num dia, so pra profissionais (POST /api/students/{id}/attendance)
func (h *AttendanceHandler) Save(w http.ResponseWriter, r *http.Request) {
	role, ok := middleware.RoleFromContext(r.Context())
	if !ok || role != "profissional" {
		h.writeError(w, http.StatusForbidden, "Somente profissionais podem registrar presenças")
		return
	}

	studentID, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do estudante inválido")
		return
	}

	var req struct {
		Date   string `json:"date"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date[:10])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Data inválida (formato esperado YYYY-MM-DD)")
		return
	}

	userID, _ := middleware.UserIDFromContext(r.Context())
	var userIDPtr *int64
	if userID != 0 {
		userIDPtr = &userID
	}

	attendance := &entities.Attendance{
		StudentID:      studentID,
		Date:           parsedDate,
		Status:         req.Status,
		MarkedByUserID: userIDPtr,
	}

	if err := h.saveUseCase.Execute(r.Context(), attendance); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, attendance)
}

// ListByTurma devolve a presenca de uma turma inteira num dia, so pra profissionais (GET /api/turmas/{id}/attendance)
func (h *AttendanceHandler) ListByTurma(w http.ResponseWriter, r *http.Request) {
	role, ok := middleware.RoleFromContext(r.Context())
	if !ok || role != "profissional" {
		h.writeError(w, http.StatusForbidden, "Somente profissionais podem listar presenças")
		return
	}

	turmaID, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID da turma inválido")
		return
	}

	dateParam := r.URL.Query().Get("date")
	var targetDate time.Time
	if dateParam == "" {
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("2006-01-02", dateParam[:10])
		if err != nil {
			h.writeError(w, http.StatusBadRequest, "Data inválida (formato esperado YYYY-MM-DD)")
			return
		}
	}

	list, err := h.getUseCase.Execute(r.Context(), turmaID, targetDate)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, list)
}

// ListByStudent devolve o historico de presenca de um aluno (GET /api/students/{id}/attendance)
func (h *AttendanceHandler) ListByStudent(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do estudante inválido")
		return
	}

	list, err := h.listByStudentUseCase.Execute(r.Context(), studentID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, list)
}

// writeJSON serializa qualquer payload como JSON com o status code informado
func (h *AttendanceHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError monta uma resposta de erro padrao {"error": "..."}
func (h *AttendanceHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
