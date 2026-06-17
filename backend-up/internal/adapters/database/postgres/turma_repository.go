package postgres

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// TurmaRepository e a implementacao em postgres do repositorio de turmas
type TurmaRepository struct {
	db *DB
}

func NewTurmaRepository(db *DB) *TurmaRepository {
	return &TurmaRepository{db: db}
}

// List busca todas as turmas ordenadas por nome
func (r *TurmaRepository) List(ctx context.Context) ([]*entities.Turma, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT id, name, created_at, updated_at FROM turmas ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	turmas := make([]*entities.Turma, 0)
	for rows.Next() {
		t := &entities.Turma{}
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		turmas = append(turmas, t)
	}
	return turmas, rows.Err()
}

// GetByID busca uma turma pelo id
func (r *TurmaRepository) GetByID(ctx context.Context, id int64) (*entities.Turma, error) {
	row := r.db.Conn().QueryRowContext(ctx, `
		SELECT id, name, created_at, updated_at FROM turmas WHERE id = $1
	`, id)
	t := &entities.Turma{}
	if err := row.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return t, nil
}

// Create insere uma turma nova e preenche id/datas geradas pelo banco
func (r *TurmaRepository) Create(ctx context.Context, t *entities.Turma) error {
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO turmas (name) VALUES ($1) RETURNING id, created_at, updated_at
	`, t.Name)
	return row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

// Update atualiza o nome da turma
func (r *TurmaRepository) Update(ctx context.Context, t *entities.Turma) error {
	_, err := r.db.Conn().ExecContext(ctx, `
		UPDATE turmas SET name = $1, updated_at = now() WHERE id = $2
	`, t.Name, t.ID)
	return err
}

// Delete apaga a turma pelo id
func (r *TurmaRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn().ExecContext(ctx, `DELETE FROM turmas WHERE id = $1`, id)
	return err
}

var _ repositories.TurmaRepository = (*TurmaRepository)(nil)
