package usecases

import (
	"context"
	"errors"
	"time"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// SaveAttendanceUseCase registra a presenca/falta de um aluno
type SaveAttendanceUseCase struct {
	repo repositories.AttendanceRepository
}

func NewSaveAttendanceUseCase(repo repositories.AttendanceRepository) *SaveAttendanceUseCase {
	return &SaveAttendanceUseCase{repo: repo}
}

// valida aluno, data e status antes de salvar a presenca
func (u *SaveAttendanceUseCase) Execute(ctx context.Context, a *entities.Attendance) error {
	if a.StudentID <= 0 {
		return errors.New("id do estudante inválido")
	}
	if a.Date.IsZero() {
		return errors.New("data inválida")
	}
	if a.Status != "present" && a.Status != "absent" {
		return errors.New("status de presença inválido")
	}
	return u.repo.Save(ctx, a)
}

// GetTurmaAttendanceUseCase busca a presenca de uma turma inteira num dia
type GetTurmaAttendanceUseCase struct {
	repo repositories.AttendanceRepository
}

func NewGetTurmaAttendanceUseCase(repo repositories.AttendanceRepository) *GetTurmaAttendanceUseCase {
	return &GetTurmaAttendanceUseCase{repo: repo}
}

// valida turma e data antes de buscar a lista de presenca
func (u *GetTurmaAttendanceUseCase) Execute(ctx context.Context, turmaID int64, date time.Time) ([]*entities.Attendance, error) {
	if turmaID <= 0 {
		return nil, errors.New("id da turma inválido")
	}
	if date.IsZero() {
		return nil, errors.New("data inválida")
	}
	return u.repo.ListByTurmaAndDate(ctx, turmaID, date)
}

// ListAttendanceByStudentUseCase busca o historico de presenca de um aluno
type ListAttendanceByStudentUseCase struct {
	repo repositories.AttendanceRepository
}

func NewListAttendanceByStudentUseCase(repo repositories.AttendanceRepository) *ListAttendanceByStudentUseCase {
	return &ListAttendanceByStudentUseCase{repo: repo}
}

// valida o id do aluno antes de buscar o historico
func (u *ListAttendanceByStudentUseCase) Execute(ctx context.Context, studentID int64) ([]*entities.Attendance, error) {
	if studentID <= 0 {
		return nil, errors.New("id do estudante inválido")
	}
	return u.repo.ListByStudent(ctx, studentID)
}
