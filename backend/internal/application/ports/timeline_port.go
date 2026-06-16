package ports

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

type TimelinePort interface {
	ListForDay(ctx context.Context, date string) ([]*entities.TimelineEvent, error)
}
