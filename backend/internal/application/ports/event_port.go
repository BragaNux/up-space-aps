package ports

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

type EventPort interface {
	ListEvents(ctx context.Context) ([]*entities.Event, error)
	RSVPEvent(ctx context.Context, id int64) (int64, error)
}
