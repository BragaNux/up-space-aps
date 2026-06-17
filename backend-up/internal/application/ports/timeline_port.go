package ports

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// TimelinePort e a porta que os usecases de timeline usam pra falar com o repositorio
type TimelinePort interface {
	ListForDay(ctx context.Context, studentID int64, date string) ([]*entities.TimelineEvent, error)
	GetByID(ctx context.Context, id int64) (*entities.TimelineEvent, error)
	Create(ctx context.Context, event *entities.TimelineEvent) error
	Update(ctx context.Context, event *entities.TimelineEvent) error
	Delete(ctx context.Context, id int64) error
}
