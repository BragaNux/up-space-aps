package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

// DeleteTimelineUseCase remove um evento da linha do tempo do aluno
type DeleteTimelineUseCase struct {
	repo ports.TimelinePort
}

func NewDeleteTimelineUseCase(repo ports.TimelinePort) *DeleteTimelineUseCase {
	return &DeleteTimelineUseCase{repo: repo}
}

// apaga o evento da timeline pelo id
func (u *DeleteTimelineUseCase) Execute(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
