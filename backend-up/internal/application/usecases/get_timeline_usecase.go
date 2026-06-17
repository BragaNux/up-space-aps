package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// GetTimelineEventUseCase busca um evento especifico da timeline do aluno
type GetTimelineEventUseCase struct {
	repo ports.TimelinePort
}

func NewGetTimelineEventUseCase(repo ports.TimelinePort) *GetTimelineEventUseCase {
	return &GetTimelineEventUseCase{repo: repo}
}

// busca o evento da timeline pelo id
func (u *GetTimelineEventUseCase) Execute(ctx context.Context, id int64) (*entities.TimelineEvent, error) {
	return u.repo.GetByID(ctx, id)
}
