package postgres

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

type EventRepository struct {
	db *DB
}

func NewEventRepository(db *DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) List(ctx context.Context) ([]*entities.Event, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT id, title, description, location, starts_at, ends_at, rsvp_count, created_at
		FROM events
		ORDER BY starts_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*entities.Event, 0)
	for rows.Next() {
		event := &entities.Event{}
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.StartsAt, &event.EndsAt, &event.RSVPCount, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (r *EventRepository) IncrementRSVP(ctx context.Context, id int64) (int64, error) {
	var count int64
	row := r.db.Conn().QueryRowContext(ctx, `
		UPDATE events
		SET rsvp_count = rsvp_count + 1
		WHERE id = $1
		RETURNING rsvp_count
	`, id)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *EventRepository) ListEvents(ctx context.Context) ([]*entities.Event, error) {
	return r.List(ctx)
}

func (r *EventRepository) RSVPEvent(ctx context.Context, id int64) (int64, error) {
	return r.IncrementRSVP(ctx, id)
}

var _ repositories.EventRepository = (*EventRepository)(nil)
