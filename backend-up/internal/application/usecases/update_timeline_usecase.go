package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// UpdateTimelineUseCase atualiza um evento existente na timeline do aluno
type UpdateTimelineUseCase struct {
	repo ports.TimelinePort
}

func NewUpdateTimelineUseCase(repo ports.TimelinePort) *UpdateTimelineUseCase {
	return &UpdateTimelineUseCase{repo: repo}
}

// confere titulo e descricao antes de salvar a atualizacao
func (u *UpdateTimelineUseCase) Execute(ctx context.Context, event *entities.TimelineEvent) error {
	if event.Title == "" {
		return errors.New("título é obrigatório")
	}
	if event.Description == "" {
		return errors.New("descrição é obrigatória")
	}
	return u.repo.Update(ctx, event)
}
