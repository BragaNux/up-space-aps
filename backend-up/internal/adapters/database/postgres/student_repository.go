package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// StudentRepository e a implementacao em postgres do repositorio de alunos
type StudentRepository struct {
	db *DB
}

func NewStudentRepository(db *DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// lista de colunas reaproveitada em todo SELECT de aluno, pra nao repetir a lista toda vez
const studentColumns = `
	id, name, presence_status, check_in_at, guardian_user_id, photo_url, turma_id, group_name,
	teacher_user_id, teacher_name, birth_date, enrollment_code, blood_type, allergies, restrictions,
	medications, created_at, updated_at
`

// scanStudent le uma linha do banco pra dentro de um Student, tratando os varios campos opcionais (nulos)
func scanStudent(scan func(dest ...any) error) (*entities.Student, error) {
	student := &entities.Student{}
	var checkIn, birthDate sql.NullTime
	var guardianUserID, turmaID, teacherUserID sql.NullInt64
	var allergies pq.StringArray

	if err := scan(
		&student.ID, &student.Name, &student.PresenceStatus, &checkIn, &guardianUserID,
		&student.PhotoURL, &turmaID, &student.GroupName, &teacherUserID, &student.TeacherName, &birthDate,
		&student.EnrollmentCode, &student.BloodType, &allergies, &student.Restrictions,
		&student.Medications, &student.CreatedAt, &student.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if checkIn.Valid {
		student.CheckInAt = &checkIn.Time
	}
	if birthDate.Valid {
		student.BirthDate = &birthDate.Time
	}
	if guardianUserID.Valid {
		id := guardianUserID.Int64
		student.GuardianUserID = &id
	}
	if turmaID.Valid {
		id := turmaID.Int64
		student.TurmaID = &id
	}
	if teacherUserID.Valid {
		id := teacherUserID.Int64
		student.TeacherUserID = &id
	}
	student.Allergies = []string(allergies)
	if student.Allergies == nil {
		student.Allergies = []string{}
	}

	return student, nil
}

// GetActive busca o primeiro aluno cadastrado, usado no modo single-student/demo
func (r *StudentRepository) GetActive(ctx context.Context) (*entities.Student, error) {
	row := r.db.Conn().QueryRowContext(ctx, `SELECT `+studentColumns+` FROM students ORDER BY id ASC LIMIT 1`)
	return scanStudent(row.Scan)
}

// GetByID busca um aluno pelo id
func (r *StudentRepository) GetByID(ctx context.Context, id int64) (*entities.Student, error) {
	row := r.db.Conn().QueryRowContext(ctx, `SELECT `+studentColumns+` FROM students WHERE id = $1`, id)
	return scanStudent(row.Scan)
}

// ListByGuardian busca os filhos vinculados a um responsavel
func (r *StudentRepository) ListByGuardian(ctx context.Context, guardianUserID int64) ([]*entities.Student, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `SELECT `+studentColumns+` FROM students WHERE guardian_user_id = $1 ORDER BY id ASC`, guardianUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectStudents(rows)
}

// List busca todos os alunos
func (r *StudentRepository) List(ctx context.Context) ([]*entities.Student, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `SELECT `+studentColumns+` FROM students ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectStudents(rows)
}

// collectStudents percorre as linhas e devolve a lista de alunos escaneados
func collectStudents(rows *sql.Rows) ([]*entities.Student, error) {
	students := make([]*entities.Student, 0)
	for rows.Next() {
		student, err := scanStudent(rows.Scan)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, rows.Err()
}

// Create insere um aluno novo e gera o codigo de matricula (#ano-id) depois do insert, ja que precisa do id gerado
func (r *StudentRepository) Create(ctx context.Context, student *entities.Student) error {
	if student.Allergies == nil {
		student.Allergies = []string{}
	}
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO students (
			name, presence_status, guardian_user_id, photo_url, turma_id, group_name,
			teacher_user_id, teacher_name, birth_date, enrollment_code, blood_type,
			allergies, restrictions, medications
		)
		VALUES ($1, 'absent', $2, $3, $4, $5, $6, $7, $8, '', $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`, student.Name, student.GuardianUserID, student.PhotoURL, student.TurmaID, student.GroupName,
		student.TeacherUserID, student.TeacherName, student.BirthDate,
		student.BloodType, pq.Array(student.Allergies), student.Restrictions, student.Medications)

	if err := row.Scan(&student.ID, &student.CreatedAt, &student.UpdatedAt); err != nil {
		return err
	}

	// Enrollment codes are generated from the row's own id, so they are only
	// known after the insert: #<year>-<zero-padded id>.
	student.EnrollmentCode = fmt.Sprintf("#%d-%04d", student.CreatedAt.Year(), student.ID)
	_, err := r.db.Conn().ExecContext(ctx, `UPDATE students SET enrollment_code = $1 WHERE id = $2`, student.EnrollmentCode, student.ID)
	return err
}

// Update atualiza o cadastro do aluno (o codigo de matricula nunca muda depois de criado)
func (r *StudentRepository) Update(ctx context.Context, student *entities.Student) error {
	if student.Allergies == nil {
		student.Allergies = []string{}
	}
	// Note: enrollment_code is generated once at creation and never updated.
	_, err := r.db.Conn().ExecContext(ctx, `
		UPDATE students SET
			name = $1, guardian_user_id = $2, photo_url = $3, turma_id = $4, group_name = $5,
			teacher_user_id = $6, teacher_name = $7, birth_date = $8,
			blood_type = $9, allergies = $10, restrictions = $11, medications = $12, updated_at = now()
		WHERE id = $13
	`, student.Name, student.GuardianUserID, student.PhotoURL, student.TurmaID, student.GroupName,
		student.TeacherUserID, student.TeacherName, student.BirthDate,
		student.BloodType, pq.Array(student.Allergies), student.Restrictions, student.Medications, student.ID)
	return err
}

// Delete apaga o aluno pelo id
func (r *StudentRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn().ExecContext(ctx, `DELETE FROM students WHERE id = $1`, id)
	return err
}

// UpdatePresence marca o aluno como presente/ausente; ao marcar presente, grava o horario de entrada (informado ou agora mesmo)
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
