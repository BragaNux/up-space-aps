package postgres

import (
	"context"
	"database/sql"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

type StudentRepository struct {
	db *DB
}

func NewStudentRepository(db *DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) GetActive(ctx context.Context) (*entities.Student, error) {
	row := r.db.Conn().QueryRowContext(ctx, `
		SELECT id, name, presence_status, check_in_at, created_at, updated_at
		FROM students
		ORDER BY id ASC
		LIMIT 1
	`)

	student := &entities.Student{}
	var checkIn sql.NullTime
	if err := row.Scan(&student.ID, &student.Name, &student.PresenceStatus, &checkIn, &student.CreatedAt, &student.UpdatedAt); err != nil {
		return nil, err
	}
	if checkIn.Valid {
		student.CheckInAt = &checkIn.Time
	}
	return student, nil
}

func (r *StudentRepository) UpdatePresence(ctx context.Context, id int64, status string, checkInAt *string) error {
	if status == "present" {
		if checkInAt != nil {
			_, err := r.db.Conn().ExecContext(ctx, `
				UPDATE students
				SET presence_status = $1, check_in_at = $2, updated_at = now()
				WHERE id = $3
			`, status, *checkInAt, id)
			return err
		}
		_, err := r.db.Conn().ExecContext(ctx, `
			UPDATE students
			SET presence_status = $1, check_in_at = now(), updated_at = now()
			WHERE id = $2
		`, status, id)
		return err
	}

	_, err := r.db.Conn().ExecContext(ctx, `
		UPDATE students
		SET presence_status = $1, updated_at = now()
		WHERE id = $2
	`, status, id)
	return err
}

var _ repositories.StudentRepository = (*StudentRepository)(nil)
