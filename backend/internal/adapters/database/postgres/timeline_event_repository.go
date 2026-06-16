package postgres

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

type TimelineEventRepository struct {
	db *DB
}

func NewTimelineEventRepository(db *DB) *TimelineEventRepository {
	return &TimelineEventRepository{db: db}
}

func (r *TimelineEventRepository) ListForDay(ctx context.Context, date string) ([]*entities.TimelineEvent, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT id, title, description, occurred_at, created_at
		FROM timeline_events
		WHERE DATE(occurred_at) = $1
		ORDER BY occurred_at ASC
	`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*entities.TimelineEvent, 0)
	for rows.Next() {
		event := &entities.TimelineEvent{}
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.OccurredAt, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

var _ repositories.TimelineEventRepository = (*TimelineEventRepository)(nil)
