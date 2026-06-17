package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"up-espaco/backend/internal/adapters/http/middleware"
	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

var errInvalidStudentID = errors.New("ID do aluno ausente ou inválido")

// StudentHandler cuida das rotas de cadastro e presenca dos alunos
type StudentHandler struct {
	getStudent     *usecases.GetStudentUseCase
	getStudentByID *usecases.GetStudentByIDUseCase
	listStudents   *usecases.ListStudentsUseCase
	listMyChildren *usecases.ListMyChildrenUseCase
	createStudent  *usecases.CreateStudentUseCase
	updateStudent  *usecases.UpdateStudentUseCase
	deleteStudent  *usecases.DeleteStudentUseCase
	updatePresence *usecases.UpdatePresenceUseCase
	listGuardians  *usecases.ListGuardiansUseCase
	userRepo       repositories.UserRepository
	turmaRepo      repositories.TurmaRepository
}

func NewStudentHandler(
	getStudent *usecases.GetStudentUseCase,
	getStudentByID *usecases.GetStudentByIDUseCase,
	listStudents *usecases.ListStudentsUseCase,
	listMyChildren *usecases.ListMyChildrenUseCase,
	createStudent *usecases.CreateStudentUseCase,
	updateStudent *usecases.UpdateStudentUseCase,
	deleteStudent *usecases.DeleteStudentUseCase,
	updatePresence *usecases.UpdatePresenceUseCase,
	listGuardians *usecases.ListGuardiansUseCase,
	userRepo repositories.UserRepository,
	turmaRepo repositories.TurmaRepository,
) *StudentHandler {
	return &StudentHandler{
		getStudent:     getStudent,
		getStudentByID: getStudentByID,
		listStudents:   listStudents,
		listMyChildren: listMyChildren,
		createStudent:  createStudent,
		updateStudent:  updateStudent,
		deleteStudent:  deleteStudent,
		updatePresence: updatePresence,
		listGuardians:  listGuardians,
		userRepo:       userRepo,
		turmaRepo:      turmaRepo,
	}
}

// withGuardians anexa a lista de responsaveis no objeto do aluno antes de devolver pro front
func (h *StudentHandler) withGuardians(r *http.Request, student *entities.Student) *entities.Student {
	if student == nil {
		return nil
	}
	guardians, err := h.listGuardians.Execute(r.Context(), student.ID)
	if err == nil {
		student.Guardians = guardians
	}
	return student
}

// GetStudent devolve o aluno "ativo" do modo single-student/demo (GET /api/student)
func (h *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	student, err := h.getStudent.Execute(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, h.withGuardians(r, student))
}

// List devolve todos os alunos cadastrados (GET /api/students)
func (h *StudentHandler) List(w http.ResponseWriter, r *http.Request) {
	students, err := h.listStudents.Execute(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, students)
}

// ListMyChildren devolve os filhos do responsavel logado (GET /api/me/children)
func (h *StudentHandler) ListMyChildren(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Autenticação necessária")
		return
	}

	students, err := h.listMyChildren.Execute(r.Context(), userID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, students)
}

// Get busca um aluno especifico, ja com os responsaveis anexados (GET /api/students/{id})
func (h *StudentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do aluno inválido")
		return
	}

	student, err := h.getStudentByID.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Aluno não encontrado")
		return
	}

	h.writeJSON(w, http.StatusOK, h.withGuardians(r, student))
}

// resolveTurmaAndTeacher busca a turma e o professor pelos ids e preenche os campos de exibicao do aluno
// (group_name/teacher_name), que ficam desnormalizados pra facilitar a leitura no front
func (h *StudentHandler) resolveTurmaAndTeacher(r *http.Request, student *entities.Student, turmaID, teacherUserID *int64) error {
	if turmaID != nil {
		turma, err := h.turmaRepo.GetByID(r.Context(), *turmaID)
		if err != nil {
			return errors.New("Turma não encontrada")
		}
		student.TurmaID = turmaID
		student.GroupName = turma.Name
	}

	if teacherUserID != nil {
		teacher, err := h.userRepo.GetByID(r.Context(), *teacherUserID)
		if err != nil {
			return errors.New("Professor(a) não encontrado(a)")
		}
		if teacher.Role != "profissional" {
			return errors.New("O professor(a) selecionado precisa ter o papel profissional")
		}
		student.TeacherUserID = teacherUserID
		student.TeacherName = teacher.Name
	}

	return nil
}

// Create cadastra um aluno novo, resolvendo responsavel (por id ou email), turma e professor (POST /api/students)
func (h *StudentHandler) Create(w http.ResponseWriter, r *http.Request) {
	role, ok := middleware.RoleFromContext(r.Context())
	if !ok || role != "profissional" {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem cadastrar crianças")
		return
	}

	var req dto.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	guardianUserID := req.GuardianUserID
	if guardianUserID == nil && req.GuardianEmail != "" {
		guardian, err := h.userRepo.GetByEmail(r.Context(), req.GuardianEmail)
		if err != nil {
			h.writeError(w, http.StatusBadRequest, "E-mail do responsável não encontrado")
			return
		}
		if guardian.Role != "responsavel" {
			h.writeError(w, http.StatusBadRequest, "A conta informada precisa ter o papel de responsável")
			return
		}
		guardianUserID = &guardian.ID
	}

	student := &entities.Student{
		Name: req.Name, GuardianUserID: guardianUserID, PhotoURL: req.PhotoURL,
		GroupName: req.GroupName, TeacherName: req.TeacherName, BirthDate: req.BirthDate,
		EnrollmentCode: req.EnrollmentCode, BloodType: req.BloodType, Allergies: req.Allergies,
		Restrictions: req.Restrictions, Medications: req.Medications,
	}

	if err := h.resolveTurmaAndTeacher(r, student, req.TurmaID, req.TeacherUserID); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.createStudent.Execute(r.Context(), student); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, student)
}

// Update edita o cadastro de um aluno (PUT /api/students/{id})
func (h *StudentHandler) Update(w http.ResponseWriter, r *http.Request) {
	role, ok := middleware.RoleFromContext(r.Context())
	if !ok || role != "profissional" {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem atualizar crianças")
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do aluno inválido")
		return
	}

	var req dto.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	student := &entities.Student{
		ID: id, Name: req.Name, GuardianUserID: req.GuardianUserID, PhotoURL: req.PhotoURL,
		GroupName: req.GroupName, TeacherName: req.TeacherName, BirthDate: req.BirthDate,
		EnrollmentCode: req.EnrollmentCode, BloodType: req.BloodType, Allergies: req.Allergies,
		Restrictions: req.Restrictions, Medications: req.Medications,
	}

	if err := h.resolveTurmaAndTeacher(r, student, req.TurmaID, req.TeacherUserID); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.updateStudent.Execute(r.Context(), student); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, student)
}

// Delete remove um aluno (DELETE /api/students/{id})
func (h *StudentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	role, ok := middleware.RoleFromContext(r.Context())
	if !ok || role != "profissional" {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem remover crianças")
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do aluno inválido")
		return
	}

	if err := h.deleteStudent.Execute(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePresence marca o aluno como presente/ausente (PATCH /api/students/{id}/presence)
func (h *StudentHandler) UpdatePresence(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do aluno inválido")
		return
	}

	var req dto.PresenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if err := h.updatePresence.Execute(r.Context(), id, req.Status); err != nil {
		if err.Error() == "status de presença inválido" {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"status": "atualizado"})
}

// parseIDParam le o {id} da URL (via mux) e converte pra int64
func parseIDParam(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		return 0, err
	}
	return id, nil
}

// parseStudentIDQuery le o query param obrigatorio "student_id", usado pra filtrar posts, timeline e marcos de um aluno
func parseStudentIDQuery(r *http.Request) (int64, error) {
	raw := r.URL.Query().Get("student_id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errInvalidStudentID
	}
	return id, nil
}

// writeJSON serializa qualquer payload como JSON com o status code informado
func (h *StudentHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError monta uma resposta de erro padrao {"error": "..."}
func (h *StudentHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
