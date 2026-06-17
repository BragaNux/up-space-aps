package postgres

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// EventRepository e a implementacao em postgres do repositorio de eventos da agenda
type EventRepository struct {
	db *DB
}

func NewEventRepository(db *DB) *EventRepository {
	return &EventRepository{db: db}
}

// List busca todos os eventos em ordem cronologica
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

// GetByID busca um evento pelo id
func (r *EventRepository) GetByID(ctx context.Context, id int64) (*entities.Event, error) {
	row := r.db.Conn().QueryRowContext(ctx, `
		SELECT id, title, description, location, starts_at, ends_at, rsvp_count, created_at
		FROM events
		WHERE id = $1
	`, id)

	event := &entities.Event{}
	if err := row.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.StartsAt, &event.EndsAt, &event.RSVPCount, &event.CreatedAt); err != nil {
		return nil, err
	}
	return event, nil
}

// Create insere um evento novo (rsvp_count comeca em zero)
func (r *EventRepository) Create(ctx context.Context, event *entities.Event) error {
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO events (title, description, location, starts_at, ends_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, rsvp_count, created_at
	`, event.Title, event.Description, event.Location, event.StartsAt, event.EndsAt)
	return row.Scan(&event.ID, &event.RSVPCount, &event.CreatedAt)
}

// Update atualiza os dados de um evento existente
func (r *EventRepository) Update(ctx context.Context, event *entities.Event) error {
	_, err := r.db.Conn().ExecContext(ctx, `
		UPDATE events
		SET title = $1, description = $2, location = $3, starts_at = $4, ends_at = $5
		WHERE id = $6
	`, event.Title, event.Description, event.Location, event.StartsAt, event.EndsAt, event.ID)
	return err
}

// Delete apaga o evento pelo id
func (r *EventRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn().ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)
	return err
}

// ToggleRSVP inverte a confirmacao de presenca do usuario no evento (numa transacao, pra manter o contador certo)
func (r *EventRepository) ToggleRSVP(ctx context.Context, eventID int64, userID int64) (int64, bool, error) {
	tx, err := r.db.Conn().BeginTx(ctx, nil)
	if err != nil {
		return 0, false, err
	}
	defer tx.Rollback()

	var exists bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM event_rsvps WHERE event_id = $1 AND user_id = $2)
	`, eventID, userID).Scan(&exists); err != nil {
		return 0, false, err
	}

	var count int64
	var confirmed bool
	if exists {
		if _, err := tx.ExecContext(ctx, `DELETE FROM event_rsvps WHERE event_id = $1 AND user_id = $2`, eventID, userID); err != nil {
			return 0, false, err
		}
		if err := tx.QueryRowContext(ctx, `
			UPDATE events SET rsvp_count = rsvp_count - 1 WHERE id = $1 RETURNING rsvp_count
		`, eventID).Scan(&count); err != nil {
			return 0, false, err
		}
		confirmed = false
	} else {
		if _, err := tx.ExecContext(ctx, `INSERT INTO event_rsvps (event_id, user_id) VALUES ($1, $2)`, eventID, userID); err != nil {
			return 0, false, err
		}
		if err := tx.QueryRowContext(ctx, `
			UPDATE events SET rsvp_count = rsvp_count + 1 WHERE id = $1 RETURNING rsvp_count
		`, eventID).Scan(&count); err != nil {
			return 0, false, err
		}
		confirmed = true
	}

	if err := tx.Commit(); err != nil {
		return 0, false, err
	}

	return count, confirmed, nil
}

// RSVPedEventIDs devolve o conjunto de eventos em que o usuario ja confirmou presenca
func (r *EventRepository) RSVPedEventIDs(ctx context.Context, userID int64) (map[int64]bool, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `SELECT event_id FROM event_rsvps WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make(map[int64]bool)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids[id] = true
	}
	return ids, rows.Err()
}

var _ repositories.EventRepository = (*EventRepository)(nil)
