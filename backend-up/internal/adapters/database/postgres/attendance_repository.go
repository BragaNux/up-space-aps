package postgres

import (
	"context"
	"time"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// AttendanceRepository e a implementacao em postgres do repositorio de presenca/falta
type AttendanceRepository struct {
	db *DB
}

func NewAttendanceRepository(db *DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

// Save grava a presenca do dia; se ja existir um registro pra aquele aluno+data, so atualiza status e quem marcou
func (r *AttendanceRepository) Save(ctx context.Context, a *entities.Attendance) error {
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO attendance (student_id, date, status, marked_by_user_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (student_id, date)
		DO UPDATE SET status = EXCLUDED.status, marked_by_user_id = EXCLUDED.marked_by_user_id
		RETURNING id, created_at
	`, a.StudentID, a.Date.Format("2006-01-02"), a.Status, a.MarkedByUserID)

	return row.Scan(&a.ID, &a.CreatedAt)
}

// GetByStudentAndDate busca o registro de presenca de um aluno numa data especifica
func (r *AttendanceRepository) GetByStudentAndDate(ctx context.Context, studentID int64, date time.Time) (*entities.Attendance, error) {
	row := r.db.Conn().QueryRowContext(ctx, `
		SELECT id, student_id, date, status, marked_by_user_id, created_at
		FROM attendance
		WHERE student_id = $1 AND date = $2
	`, studentID, date.Format("2006-01-02"))

	a := &entities.Attendance{}
	if err := row.Scan(&a.ID, &a.StudentID, &a.Date, &a.Status, &a.MarkedByUserID, &a.CreatedAt); err != nil {
		return nil, err
	}
	return a, nil
}

// ListByTurmaAndDate busca a presenca de todos os alunos de uma turma numa data
func (r *AttendanceRepository) ListByTurmaAndDate(ctx context.Context, turmaID int64, date time.Time) ([]*entities.Attendance, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT a.id, a.student_id, a.date, a.status, a.marked_by_user_id, a.created_at
		FROM attendance a
		JOIN students s ON s.id = a.student_id
		WHERE s.turma_id = $1 AND a.date = $2
	`, turmaID, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*entities.Attendance, 0)
	for rows.Next() {
		a := &entities.Attendance{}
		if err := rows.Scan(&a.ID, &a.StudentID, &a.Date, &a.Status, &a.MarkedByUserID, &a.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

// ListByStudent busca o historico de presenca de um aluno, mais recente primeiro
func (r *AttendanceRepository) ListByStudent(ctx context.Context, studentID int64) ([]*entities.Attendance, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT id, student_id, date, status, marked_by_user_id, created_at
		FROM attendance
		WHERE student_id = $1
		ORDER BY date DESC
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*entities.Attendance, 0)
	for rows.Next() {
		a := &entities.Attendance{}
		if err := rows.Scan(&a.ID, &a.StudentID, &a.Date, &a.Status, &a.MarkedByUserID, &a.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

var _ repositories.AttendanceRepository = (*AttendanceRepository)(nil)
