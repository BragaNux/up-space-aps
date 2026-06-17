package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

// DeleteEventUseCase remove um evento da agenda
type DeleteEventUseCase struct {
	repo ports.EventPort
}

func NewDeleteEventUseCase(repo ports.EventPort) *DeleteEventUseCase {
	return &DeleteEventUseCase{repo: repo}
}

// apaga o evento pelo id
func (u *DeleteEventUseCase) Execute(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
