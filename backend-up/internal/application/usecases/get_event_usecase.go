package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// GetEventUseCase busca um evento especifico da agenda
type GetEventUseCase struct {
	repo ports.EventPort
}

func NewGetEventUseCase(repo ports.EventPort) *GetEventUseCase {
	return &GetEventUseCase{repo: repo}
}

// busca o evento pelo id
func (u *GetEventUseCase) Execute(ctx context.Context, id int64) (*entities.Event, error) {
	return u.repo.GetByID(ctx, id)
}
