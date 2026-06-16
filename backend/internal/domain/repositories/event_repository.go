package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

type EventRepository interface {
	List(ctx context.Context) ([]*entities.Event, error)
	IncrementRSVP(ctx context.Context, id int64) (int64, error)
}
