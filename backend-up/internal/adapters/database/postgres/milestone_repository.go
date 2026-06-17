package postgres

import (
	"context"
	"database/sql"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// MilestoneRepository e a implementacao em postgres do repositorio de marcos de desenvolvimento
type MilestoneRepository struct {
	db *DB
}

func NewMilestoneRepository(db *DB) *MilestoneRepository {
	return &MilestoneRepository{db: db}
}

// scanMilestone le uma linha do banco pra dentro de um Milestone, tratando achieved_at nulo (marco ainda nao alcancado)
func scanMilestone(scan func(dest ...any) error) (*entities.Milestone, error) {
	m := &entities.Milestone{}
	var achievedAt sql.NullTime
	if err := scan(&m.ID, &m.StudentID, &m.Title, &m.Category, &m.Description, &achievedAt, &m.Done, &m.CreatedAt); err != nil {
		return nil, err
	}
	if achievedAt.Valid {
		m.AchievedAt = &achievedAt.Time
	}
	return m, nil
}

// List busca os marcos de um aluno, mais recentes primeiro
func (r *MilestoneRepository) List(ctx context.Context, studentID int64) ([]*entities.Milestone, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT id, student_id, title, category, description, achieved_at, done, created_at
		FROM milestones
		WHERE student_id = $1
		ORDER BY achieved_at DESC NULLS LAST, created_at DESC
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	milestones := make([]*entities.Milestone, 0)
	for rows.Next() {
		m, err := scanMilestone(rows.Scan)
		if err != nil {
			return nil, err
		}
		milestones = append(milestones, m)
	}
	return milestones, rows.Err()
}

// GetByID busca um marco pelo id
func (r *MilestoneRepository) GetByID(ctx context.Context, id int64) (*entities.Milestone, error) {
	row := r.db.Conn().QueryRowContext(ctx, `
		SELECT id, student_id, title, category, description, achieved_at, done, created_at
		FROM milestones WHERE id = $1
	`, id)
	return scanMilestone(row.Scan)
}

// Create insere um marco novo
func (r *MilestoneRepository) Create(ctx context.Context, m *entities.Milestone) error {
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO milestones (student_id, title, category, description, achieved_at, done)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`, m.StudentID, m.Title, m.Category, m.Description, m.AchievedAt, m.Done)
	return row.Scan(&m.ID, &m.CreatedAt)
}

// Update atualiza um marco existente
func (r *MilestoneRepository) Update(ctx context.Context, m *entities.Milestone) error {
	_, err := r.db.Conn().ExecContext(ctx, `
		UPDATE milestones
		SET title = $1, category = $2, description = $3, achieved_at = $4, done = $5
		WHERE id = $6
	`, m.Title, m.Category, m.Description, m.AchievedAt, m.Done, m.ID)
	return err
}

// Delete apaga o marco pelo id
func (r *MilestoneRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn().ExecContext(ctx, `DELETE FROM milestones WHERE id = $1`, id)
	return err
}

var _ repositories.MilestoneRepository = (*MilestoneRepository)(nil)
