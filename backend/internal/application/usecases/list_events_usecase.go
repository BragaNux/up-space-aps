package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

type ListEventsUseCase struct {
	repo ports.EventPort
}

func NewListEventsUseCase(repo ports.EventPort) *ListEventsUseCase {
	return &ListEventsUseCase{repo: repo}
}

func (u *ListEventsUseCase) Execute(ctx context.Context) ([]*entities.Event, error) {
	return u.repo.ListEvents(ctx)
}
