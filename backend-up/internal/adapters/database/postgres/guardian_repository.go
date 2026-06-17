package postgres

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// GuardianRepository e a implementacao em postgres do repositorio de responsaveis
type GuardianRepository struct {
	db *DB
}

func NewGuardianRepository(db *DB) *GuardianRepository {
	return &GuardianRepository{db: db}
}

// scanGuardian le uma linha do banco pra dentro de um Guardian
func scanGuardian(scan func(dest ...any) error) (*entities.Guardian, error) {
	g := &entities.Guardian{}
	if err := scan(&g.ID, &g.StudentID, &g.Name, &g.Relation, &g.Phone, &g.AvatarURL, &g.Authorized); err != nil {
		return nil, err
	}
	return g, nil
}

// ListByStudent busca os responsaveis de um aluno
func (r *GuardianRepository) ListByStudent(ctx context.Context, studentID int64) ([]*entities.Guardian, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT id, student_id, name, relation, phone, avatar_url, authorized
		FROM student_guardians
		WHERE student_id = $1
		ORDER BY id ASC
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	guardians := make([]*entities.Guardian, 0)
	for rows.Next() {
		g, err := scanGuardian(rows.Scan)
		if err != nil {
			return nil, err
		}
		guardians = append(guardians, g)
	}
	return guardians, rows.Err()
}

// GetByID busca um responsavel pelo id
func (r *GuardianRepository) GetByID(ctx context.Context, id int64) (*entities.Guardian, error) {
	row := r.db.Conn().QueryRowContext(ctx, `
		SELECT id, student_id, name, relation, phone, avatar_url, authorized
		FROM student_guardians WHERE id = $1
	`, id)
	return scanGuardian(row.Scan)
}

// Create insere um responsavel novo
func (r *GuardianRepository) Create(ctx context.Context, g *entities.Guardian) error {
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO student_guardians (student_id, name, relation, phone, avatar_url, authorized)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, g.StudentID, g.Name, g.Relation, g.Phone, g.AvatarURL, g.Authorized)
	return row.Scan(&g.ID)
}

// Update atualiza os dados de um responsavel
func (r *GuardianRepository) Update(ctx context.Context, g *entities.Guardian) error {
	_, err := r.db.Conn().ExecContext(ctx, `
		UPDATE student_guardians
		SET name = $1, relation = $2, phone = $3, avatar_url = $4, authorized = $5
		WHERE id = $6
	`, g.Name, g.Relation, g.Phone, g.AvatarURL, g.Authorized, g.ID)
	return err
}

// Delete apaga o responsavel pelo id
func (r *GuardianRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn().ExecContext(ctx, `DELETE FROM student_guardians WHERE id = $1`, id)
	return err
}

var _ repositories.GuardianRepository = (*GuardianRepository)(nil)
