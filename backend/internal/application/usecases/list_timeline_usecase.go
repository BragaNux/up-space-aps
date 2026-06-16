package usecases

import (
	"context"
	"time"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

type ListTimelineUseCase struct {
	repo ports.TimelinePort
}

func NewListTimelineUseCase(repo ports.TimelinePort) *ListTimelineUseCase {
	return &ListTimelineUseCase{repo: repo}
}

func (u *ListTimelineUseCase) Execute(ctx context.Context) ([]*entities.TimelineEvent, error) {
	return u.repo.ListForDay(ctx, time.Now().Format("2006-01-02"))
}
