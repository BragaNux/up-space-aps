package repositories

import (
	"context"
	"time"

	"up-espaco/backend/internal/domain/entities"
)

// AttendanceRepository cuida do registro de presenca/falta dos alunos
type AttendanceRepository interface {
	Save(ctx context.Context, attendance *entities.Attendance) error                                              // salva (cria ou atualiza) a presenca de um aluno no dia
	GetByStudentAndDate(ctx context.Context, studentID int64, date time.Time) (*entities.Attendance, error)       // busca a presenca de um aluno numa data
	ListByTurmaAndDate(ctx context.Context, turmaID int64, date time.Time) ([]*entities.Attendance, error)        // lista a presenca de toda a turma num dia
	ListByStudent(ctx context.Context, studentID int64) ([]*entities.Attendance, error)                           // lista o historico de presenca de um aluno
}
