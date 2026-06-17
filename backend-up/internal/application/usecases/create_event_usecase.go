package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// CreateEventUseCase cria um evento novo na agenda
type CreateEventUseCase struct {
	repo ports.EventPort
}

func NewCreateEventUseCase(repo ports.EventPort) *CreateEventUseCase {
	return &CreateEventUseCase{repo: repo}
}

// confere titulo, local e se a data de termino vem depois da de inicio
func (u *CreateEventUseCase) Execute(ctx context.Context, event *entities.Event) error {
	if event.Title == "" {
		return errors.New("título é obrigatório")
	}
	if event.Location == "" {
		return errors.New("local é obrigatório")
	}
	if event.EndsAt.Before(event.StartsAt) {
		return errors.New("a data de término deve ser depois da data de início")
	}
	return u.repo.Create(ctx, event)
}
