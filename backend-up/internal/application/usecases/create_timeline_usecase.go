package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// CreateTimelineUseCase cria um evento novo na linha do tempo do aluno
type CreateTimelineUseCase struct {
	repo ports.TimelinePort
}

func NewCreateTimelineUseCase(repo ports.TimelinePort) *CreateTimelineUseCase {
	return &CreateTimelineUseCase{repo: repo}
}

// confere aluno, titulo e descricao antes de criar o evento
func (u *CreateTimelineUseCase) Execute(ctx context.Context, event *entities.TimelineEvent) error {
	if event.StudentID == 0 {
		return errors.New("id do aluno é obrigatório")
	}
	if event.Title == "" {
		return errors.New("título é obrigatório")
	}
	if event.Description == "" {
		return errors.New("descrição é obrigatória")
	}
	return u.repo.Create(ctx, event)
}
