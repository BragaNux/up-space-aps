package postgres

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// TimelineEventRepository e a implementacao em postgres do repositorio de timeline dos alunos
type TimelineEventRepository struct {
	db *DB
}

func NewTimelineEventRepository(db *DB) *TimelineEventRepository {
	return &TimelineEventRepository{db: db}
}

// ListForDay busca os eventos da timeline de um aluno numa data especifica
func (r *TimelineEventRepository) ListForDay(ctx context.Context, studentID int64, date string) ([]*entities.TimelineEvent, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT id, student_id, title, description, occurred_at, created_at
		FROM timeline_events
		WHERE student_id = $1 AND DATE(occurred_at) = $2
		ORDER BY occurred_at ASC
	`, studentID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*entities.TimelineEvent, 0)
	for rows.Next() {
		event := &entities.TimelineEvent{}
		if err := rows.Scan(&event.ID, &event.StudentID, &event.Title, &event.Description, &event.OccurredAt, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

// GetByID busca um evento da timeline pelo id
func (r *TimelineEventRepository) GetByID(ctx context.Context, id int64) (*entities.TimelineEvent, error) {
	row := r.db.Conn().QueryRowContext(ctx, `
		SELECT id, student_id, title, description, occurred_at, created_at
		FROM timeline_events
		WHERE id = $1
	`, id)

	event := &entities.TimelineEvent{}
	if err := row.Scan(&event.ID, &event.StudentID, &event.Title, &event.Description, &event.OccurredAt, &event.CreatedAt); err != nil {
		return nil, err
	}
	return event, nil
}

// Create insere um evento novo na timeline
func (r *TimelineEventRepository) Create(ctx context.Context, event *entities.TimelineEvent) error {
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO timeline_events (student_id, title, description, occurred_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, event.StudentID, event.Title, event.Description, event.OccurredAt)
	return row.Scan(&event.ID, &event.CreatedAt)
}

// Update atualiza um evento existente da timeline
func (r *TimelineEventRepository) Update(ctx context.Context, event *entities.TimelineEvent) error {
	_, err := r.db.Conn().ExecContext(ctx, `
		UPDATE timeline_events
		SET title = $1, description = $2, occurred_at = $3
		WHERE id = $4
	`, event.Title, event.Description, event.OccurredAt, event.ID)
	return err
}

// Delete apaga o evento da timeline pelo id
func (r *TimelineEventRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn().ExecContext(ctx, `DELETE FROM timeline_events WHERE id = $1`, id)
	return err
}

var _ repositories.TimelineEventRepository = (*TimelineEventRepository)(nil)
