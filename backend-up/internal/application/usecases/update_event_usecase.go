package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// UpdateEventUseCase atualiza um evento existente na agenda
type UpdateEventUseCase struct {
	repo ports.EventPort
}

func NewUpdateEventUseCase(repo ports.EventPort) *UpdateEventUseCase {
	return &UpdateEventUseCase{repo: repo}
}

// confere titulo, local e datas antes de salvar a atualizacao
func (u *UpdateEventUseCase) Execute(ctx context.Context, event *entities.Event) error {
	if event.Title == "" {
		return errors.New("título é obrigatório")
	}
	if event.Location == "" {
		return errors.New("local é obrigatório")
	}
	if event.EndsAt.Before(event.StartsAt) {
		return errors.New("a data de término deve ser depois da data de início")
	}
	return u.repo.Update(ctx, event)
}
