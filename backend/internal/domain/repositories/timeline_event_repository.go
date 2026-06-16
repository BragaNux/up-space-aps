package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

type TimelineEventRepository interface {
	ListForDay(ctx context.Context, date string) ([]*entities.TimelineEvent, error)
}
